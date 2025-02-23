package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
)

type SiteClient struct {
	c *Client
}

func (c *Client) Site() *SiteClient {
	return &SiteClient{c}
}

func (s *SiteClient) Read(ctx context.Context) (*Site, error) {
	req := &request{
		method:          "GET",
		path:            "ghost/api/admin/site/",
		unauthenticated: true,
	}
	resp, err := expect(200)(s.c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	var site Site
	if err := decode(resp.Body, "site", &site); err != nil {
		return nil, err
	}
	return &site, nil
}

type UserClient struct {
	c *Client
}

func (c *Client) Users() *UserClient {
	return &UserClient{c}
}

func (u *UserClient) List(ctx context.Context) ([]User, error) {
	req := &request{
		method: "GET",
		path:   "ghost/api/admin/users/",
	}
	resp, err := expect(200)(u.c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	var users []User
	if err := decode(resp.Body, "users", &users); err != nil {
		return nil, err
	}
	return users, nil
}

type ThemeClient struct {
	c *Client
}

func (c *Client) Themes() *ThemeClient {
	return &ThemeClient{c}
}

func (t *ThemeClient) Upload(ctx context.Context, name string, f io.Reader) (*Theme, error) {
	form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)
	fw, err := writer.CreateFormFile("file", name+".zip")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(fw, f)
	if err != nil {
		return nil, err
	}

	writer.Close()

	req := &request{
		method:   "POST",
		path:     "ghost/api/admin/themes/upload",
		raw_body: form,
		headers: map[string]string{
			"Content-Type": writer.FormDataContentType(),
		},
	}
	resp, err := expect(200)(t.c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	var theme []Theme
	if err := decode(resp.Body, "themes", &theme); err != nil {
		return nil, err
	}
	return &theme[0], err
}

func (t *ThemeClient) Activate(ctx context.Context, name string) (*Theme, error) {
	req := &request{
		method: "PUT",
		path:   fmt.Sprintf("ghost/api/admin/themes/%s/activate/", name),
	}
	resp, err := expect(200)(t.c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	var theme []Theme
	if err := decode(resp.Body, "themes", &theme); err != nil {
		return nil, err
	}
	return &theme[0], err
}

type WebhookClient struct {
	c *Client
}

func (c *Client) Webhooks() *WebhookClient {
	return &WebhookClient{c}
}

func (w *WebhookClient) Create(ctx context.Context, webhook *Webhook) (*Webhook, error) {
	req := &request{
		method: "POST",
		path:   "ghost/api/admin/webhooks/",
		body: map[string][]Webhook{
			"webhooks": {*webhook},
		},
	}
	resp, err := expect(201)(w.c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	var data []Webhook
	if err := decode(resp.Body, "webhooks", &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	return &data[0], nil
}

func (w *WebhookClient) Update(ctx context.Context, webhook *Webhook) (*Webhook, error) {
	req := &request{
		method: "PUT",
		path:   fmt.Sprintf("ghost/api/admin/webhooks/%s/", webhook.Id),
		body: map[string][]Webhook{
			"webhooks": {*webhook},
		},
	}
	resp, err := expect(200)(w.c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	var data []Webhook
	if err := decode(resp.Body, "webhooks", &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	return &data[0], nil
}

func (w *WebhookClient) Delete(ctx context.Context, id string) error {
	req := &request{
		method: "DELETE",
		path:   fmt.Sprintf("ghost/api/admin/webhooks/%s/", id),
	}
	_, err := expect(204)(w.c.do(ctx, req))
	if err != nil {
		return err
	}
	return nil
}
