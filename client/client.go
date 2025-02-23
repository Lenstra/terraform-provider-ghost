package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	utils "github.com/Lenstra/go-utils/http"
	"github.com/go-viper/mapstructure/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

const (
	addressEnvVar                  = "GHOST_ADDRESS"
	adminApiKeyEnvVar              = "GHOST_ADMIN_API_KEY"
	certificateAuthorityPathEnvVar = "GHOST_CERTIFICATE_AUTHORITY_PATH"
	clientCertificatePathEnvVar    = "GHOST_CLIENT_CERTIFICATE_PATH"
	clientKeyPathEnvVar            = "GHOST_CLIENT_KEY_PATH"
)

type Config struct {
	Address     string
	AdminAPIKey string
	HttpClient  *http.Client
	TLSConfig   *TLSConfig
	Logger      zerolog.Logger
}

type TLSConfig struct {
	CertificateAuthorityPath string
	ClientCertificatePath    string
	ClientKeyPath            string
}

func defaultTLSConfig() *TLSConfig {
	return &TLSConfig{
		CertificateAuthorityPath: os.Getenv(certificateAuthorityPathEnvVar),
		ClientCertificatePath:    os.Getenv(clientCertificatePathEnvVar),
		ClientKeyPath:            os.Getenv(clientKeyPathEnvVar),
	}
}

func (t *TLSConfig) Transport() (*http.Transport, error) {
	if t == nil {
		return nil, nil
	}

	var certificates []tls.Certificate
	if t.ClientCertificatePath != "" || t.ClientKeyPath != "" {
		if t.ClientCertificatePath == "" {
			return nil, fmt.Errorf("ClientCertificatePath must be set when ClientKeyPath is")
		}
		if t.ClientKeyPath == "" {
			return nil, fmt.Errorf("ClientKeyPath must be set when ClientCertificatePath is")
		}
		cert, err := tls.LoadX509KeyPair(t.ClientCertificatePath, t.ClientKeyPath)
		if err != nil {
			return nil, err
		}
		certificates = append(certificates, cert)
	}

	var caCertPool *x509.CertPool
	if t.CertificateAuthorityPath != "" {
		caCertPool = x509.NewCertPool()
		caCert, err := os.ReadFile(t.CertificateAuthorityPath)
		if err != nil {
			return nil, err
		}

		caCertPool.AppendCertsFromPEM(caCert)
	}

	tlsConfig := &tls.Config{
		Certificates: certificates,
		RootCAs:      caCertPool,
	}
	return &http.Transport{TLSClientConfig: tlsConfig}, nil
}

func DefaultConfig() (*Config, error) {
	tlsConfig := defaultTLSConfig()
	transport, err := tlsConfig.Transport()
	if err != nil {
		return nil, err
	}

	return &Config{
		Address:     os.Getenv(addressEnvVar),
		AdminAPIKey: os.Getenv(adminApiKeyEnvVar),
		HttpClient: &http.Client{
			Timeout:   60 * time.Second,
			Transport: transport,
		},
	}, nil
}

type Client struct {
	conf *Config
}

func NewClient(config *Config) (*Client, error) {
	defConfig, err := DefaultConfig()
	if err != nil {
		return nil, err
	}

	if config.Address == "" {
		config.Address = defConfig.Address
	}
	if config.AdminAPIKey == "" {
		config.AdminAPIKey = defConfig.AdminAPIKey
	}
	if config.HttpClient == nil {
		config.HttpClient = defConfig.HttpClient
	}

	return &Client{conf: config}, nil
}

type request struct {
	method          string
	path            string
	body            any
	raw_body        io.Reader
	unauthenticated bool
	cookies         []*http.Cookie
	headers         map[string]string
}

func (c *Client) toHttpRequest(ctx context.Context, r *request) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", c.conf.Address, r.path)
	req, err := http.NewRequestWithContext(ctx, r.method, url, r.raw_body)
	if err != nil {
		return nil, err
	}

	if r.body != nil {
		body, err := json.Marshal(r.body)
		if err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(bytes.NewBuffer(body))

		req.Header.Set("Content-Type", "application/json")
	}

	if !r.unauthenticated {
		token, err := c.getToken()
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", fmt.Sprintf("Ghost %s", token))
	}

	for _, cookie := range r.cookies {
		req.AddCookie(cookie)
	}

	for key, value := range r.headers {
		req.Header.Add(key, value)
	}

	return req, nil
}

func (c *Client) getToken() (string, error) {
	parts := strings.Split(c.conf.AdminAPIKey, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("Invalid API key format")
	}
	id, secretHex := parts[0], parts[1]

	secret, err := hex.DecodeString(secretHex)
	if err != nil {
		return "", fmt.Errorf("Error decoding secret: %w", err)
	}

	iat := time.Now().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": iat,
		"exp": iat + 5*60,
		"aud": "/admin/",
	})

	token.Header["alg"] = "HS256"
	token.Header["typ"] = "JWT"
	token.Header["kid"] = id

	return token.SignedString(secret)
}

func (c *Client) do(ctx context.Context, r *request) (*http.Response, error) {
	req, err := c.toHttpRequest(ctx, r)
	if err != nil {
		return nil, err
	}

	if c.conf.Logger.Trace().Enabled() {
		var body []byte
		if req.Body != nil {
			body, err = io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
		}
		c.conf.Logger.Trace().Str("method", r.method).Str("path", r.path).Str("body", string(body)).Msg("request")
		if req.Body != nil {
			req.Body.Close()
			req.Body = io.NopCloser(bytes.NewBuffer(body))
		}
	}

	err = utils.LogRequest(c.conf.Logger.Trace(), req, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.conf.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	err = utils.LogResponse(c.conf.Logger.Trace(), resp, nil)
	if err != nil {
		return nil, err
	}

	if c.conf.Logger.Trace().Enabled() {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		c.conf.Logger.Trace().Str("method", r.method).Str("path", r.path).Int("status_code", resp.StatusCode).Str("body", string(body)).Msg("response")
		resp.Body.Close()
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	return resp, nil
}

func decode(r io.ReadCloser, name string, out any) error {
	defer r.Close()
	var data map[string]any
	dec := json.NewDecoder(r)
	if err := dec.Decode(&data); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  out,
	})
	if err != nil {
		return err
	}
	if err := decoder.Decode(data[name]); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	return nil
}

func expect(codes ...int) func(*http.Response, error) (*http.Response, error) {
	return func(res *http.Response, err error) (*http.Response, error) {
		if err != nil {
			return nil, err
		}
		if !slices.Contains(codes, res.StatusCode) {
			return res, fmt.Errorf("got status code %s, expected status code in %v", res.Status, codes)
		}
		return res, nil
	}
}
