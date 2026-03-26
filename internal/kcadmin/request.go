package kcadmin

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func (c *Client) doJSON(ctx context.Context, method, path string, body any, target any) error {
	return c.doJSONWithRetry(ctx, method, path, body, target, true)
}

func (c *Client) doJSONWithRetry(ctx context.Context, method, path string, body any, target any, allowRetry bool) error {
	token, err := c.tokens.Token(ctx)
	if err != nil {
		return err
	}

	var payload io.Reader
	if body != nil {
		buffer, err := json.Marshal(body)
		if err != nil {
			return err
		}
		payload = bytes.NewBuffer(buffer)
	}

	request, err := http.NewRequestWithContext(ctx, method, strings.TrimRight(c.cfg.AdminAPIBaseURL, "/")+path, payload)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Accept", "application/json")
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := c.http.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized && allowRetry {
		c.tokens.Invalidate()
		return c.doJSONWithRetry(ctx, method, path, body, target, false)
	}

	if response.StatusCode >= http.StatusBadRequest {
		bodyBytes, _ := io.ReadAll(response.Body)
		return &APIError{
			StatusCode: response.StatusCode,
			Path:       path,
			Body:       string(bodyBytes),
		}
	}

	if target == nil || response.StatusCode == http.StatusNoContent {
		return nil
	}
	return decodeJSON(response.Body, target)
}

func decodeJSON(reader io.Reader, target any) error {
	return json.NewDecoder(reader).Decode(target)
}
