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

func (c *valueServiceClient) GetValue(ctx context.Context, itemId string) (int, error) {
	url := fmt.Sprintf("%s/api/v1/internal/counts/%s", c.baseURL, itemId)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return 0, domain.ErrItemNotFound
	}
	if resp.StatusCode >= 400 {
		return 0, fmt.Errorf("failed to get value: status %d", resp.StatusCode)
	}

	var result struct {
		CurrentValue int `json:"currentValue"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.CurrentValue, nil
}

func (c *valueServiceClient) GetValues(ctx context.Context, itemIds []string) (map[string]int, error) {
	if len(itemIds) == 0 {
		return make(map[string]int), nil
	}

	url := fmt.Sprintf("%s/api/v1/internal/counts", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for _, id := range itemIds {
		q.Add("itemIds", id)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("failed to get values: status %d", resp.StatusCode)
	}

	var result struct {
		Counts []struct {
			ItemID       string `json:"itemId"`
			CurrentValue int    `json:"currentValue"`
		} `json:"counts"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	values := make(map[string]int)
	for _, count := range result.Counts {
		values[count.ItemID] = count.CurrentValue
	}
	return values, nil
}

func (c *valueServiceClient) GetHistory(ctx context.Context, itemId string) ([]domain.HistoryEntry, error) {
	url := fmt.Sprintf("%s/api/v1/counts/%s/history", c.baseURL, itemId)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, domain.ErrItemNotFound
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("failed to get history: status %d", resp.StatusCode)
	}

	var history []domain.HistoryEntry
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
		return nil, err
	}

	return history, nil
}
