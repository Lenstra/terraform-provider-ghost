package client

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/ory/dockertest"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

var testLogTraceEnvVar = "TEST_LOG_TRACE"

func TestServer(t *testing.T) *Client {
	t.Helper()

	ctx := context.Background()
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)

	err = pool.Client.Ping()
	require.NoError(t, err)

	ghost, err := pool.Run("ghost", "latest", []string{"NODE_ENV=development"})
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := pool.Purge(ghost); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	})

	port := ghost.GetPort("2368/tcp")
	address := fmt.Sprintf("http://localhost:%s", port)

	config := &Config{
		Address: address,
	}
	if os.Getenv(testLogTraceEnvVar) != "" {
		config.Logger = zerolog.New(os.Stdout)
	}
	client, err := NewClient(config)
	require.NoError(t, err)

	err = pool.Retry(func() error {
		req := &request{
			method:          "GET",
			path:            "",
			unauthenticated: true,
		}
		_, err := expect(200)(client.do(ctx, req))
		return err
	})
	require.NoError(t, err)

	req := &request{
		method: "POST",
		path:   "ghost/api/admin/authentication/setup/",
		body: map[string]any{
			"setup": []any{
				map[string]any{
					"blogTitle": "Test",
					"email":     "user@test.com",
					"name":      "Test User",
					"password":  "thooCha8ph",
				},
			},
		},
		unauthenticated: true,
	}
	_, err = expect(201)(client.do(ctx, req))
	require.NoError(t, err)

	req = &request{
		method: "POST",
		path:   "ghost/api/admin/session",
		body: map[string]any{
			"password": "thooCha8ph",
			"username": "user@test.com",
		},
		unauthenticated: true,
	}
	resp, err := expect(201)(client.do(ctx, req))
	require.NoError(t, err)

	cookie, err := getSessionCookie(resp)
	require.NoError(t, err)

	req = &request{
		method: "POST",
		path:   "ghost/api/admin/integrations/?include=api_keys,webhooks",
		body: map[string][]any{
			"integrations": {
				map[string]any{
					"name": "Test",
				},
			},
		},
		unauthenticated: true,
		cookies:         []*http.Cookie{cookie},
	}
	resp, err = expect(201)(client.do(ctx, req))
	require.NoError(t, err)

	type Integration struct {
		APIKeys []struct {
			Type   string `json:"type"`
			Secret string `json:"secret"`
		} `json:"api_keys"`
	}

	var integrations []Integration
	err = decode(resp.Body, "integrations", &integrations)
	require.NoError(t, err)
	require.Len(t, integrations, 1)
	require.Len(t, integrations[0].APIKeys, 2)
	require.Equal(t, "admin", integrations[0].APIKeys[1].Type)

	config.AdminAPIKey = integrations[0].APIKeys[1].Secret

	t.Setenv("GHOST_ADDRESS", config.Address)
	t.Setenv("GHOST_ADMIN_API_KEY", config.AdminAPIKey)

	return client
}

func getSessionCookie(resp *http.Response) (*http.Cookie, error) {
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "ghost-admin-api-session" {
			return cookie, nil
		}
	}
	return nil, fmt.Errorf("failed to find session cookie")
}
