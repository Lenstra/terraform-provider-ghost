package client

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSite(t *testing.T) {
	ctx := context.Background()
	client := TestServer(t)

	site, err := client.Site().Read(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, site.Title)
}

func TestUsers(t *testing.T) {
	ctx := context.Background()
	client := TestServer(t)

	users, err := client.Users().List(ctx)
	require.NoError(t, err)
	require.Len(t, users, 1)
	require.Equal(t, "Test User", users[0].Name)
}

func TestTheme(t *testing.T) {
	ctx := context.Background()
	client := TestServer(t)

	f, err := os.Open("../tests/casper.zip")
	require.NoError(t, err)
	theme, err := client.Themes().Upload(ctx, "test-theme", f)
	require.NoError(t, err)
	require.Equal(t, "test-theme", theme.Name)
	require.Equal(t, false, theme.Active)

	theme, err = client.Themes().Activate(ctx, "casper")
	require.NoError(t, err)
	require.Equal(t, "casper", theme.Name)
	require.Equal(t, true, theme.Active)
}

func TestWebhooks(t *testing.T) {
	ctx := context.Background()
	client := TestServer(t)

	wh := &Webhook{
		Event:     "post.added",
		TargetUrl: "https://example.com/hook/",
	}
	wh, err := client.Webhooks().Create(ctx, wh)
	require.NoError(t, err)
	require.NotEmpty(t, wh.Id)
	require.Equal(t, "post.added", wh.Event)
	require.Equal(t, "https://example.com/hook/", wh.TargetUrl)

	wh.TargetUrl = "https://example.com/updated-hook/"
	wh, err = client.Webhooks().Update(ctx, wh)
	require.NoError(t, err)
	require.Equal(t, "https://example.com/updated-hook/", wh.TargetUrl)

	err = client.Webhooks().Delete(ctx, wh.Id)
	require.NoError(t, err)
}
