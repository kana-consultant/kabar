package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"seo-backend/internal/models"
	"strings"
	"time"
)

type HTTPClient struct {
}

func (c *HTTPClient) Send(
	method string,
	url string,
	body map[string]interface{},
	headers map[string]string,
	timeout int,
	retry int,
	apiKey string,
) ([]byte, error) {

	jsonBody, _ := json.Marshal(body)

	var lastErr error

	for i := 0; i < retry; i++ {

		req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		if apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+apiKey)
		}

		client := &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(2 * time.Second)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return body, nil
		}

		lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		time.Sleep(2 * time.Second)
	}

	return nil, lastErr
}

func BuildPayload(mapping map[string]interface{}, draft models.Draft) map[string]interface{} {

	result := map[string]interface{}{}

	for k, v := range mapping {

		if str, ok := v.(string); ok {

			str = strings.ReplaceAll(str, "{title}", draft.Title)
			str = strings.ReplaceAll(str, "{topic}", draft.Topic)
			str = strings.ReplaceAll(str, "{content}", draft.Article)

			excerpt := draft.Article
			if len(excerpt) > 200 {
				excerpt = excerpt[:200]
			}

			str = strings.ReplaceAll(str, "{excerpt}", excerpt)

			if draft.ImageURL != nil {
				str = strings.ReplaceAll(str, "{image_url}", *draft.ImageURL)
			}

			result[k] = str
		} else {
			result[k] = v
		}
	}

	return result
}
