package outbound

import (
	"bytes"
	"context"
	"count-management-service/domain"
	"encoding/json"
	"fmt"
	"net/http"
)

type valueServiceClient struct {
	baseURL string
	client  *http.Client
}

func NewValueServiceClient(baseURL string) domain.ValueServiceClient {
	return &valueServiceClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *valueServiceClient) InitializeValue(ctx context.Context, itemId string, initialValue int) error {
	url := fmt.Sprintf("%s/api/v1/internal/counts", c.baseURL)
	body := map[string]interface{}{
		"itemId":       itemId,
		"initialValue": initialValue,
	}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to initialize value: status %d", resp.StatusCode)
	}

	return nil
}

func (c *valueServiceClient) DeleteValue(ctx context.Context, itemId string) error {
	url := fmt.Sprintf("%s/api/v1/internal/counts/%s", c.baseURL, itemId)
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode != 404 {
		return fmt.Errorf("failed to delete value: status %d", resp.StatusCode)
	}

	return nil
}
