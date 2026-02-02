// Copyright 2024 The frp Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dns

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type CloudflareProvider struct {
	apiKey string
	zoneID string
	client *http.Client
}

type cloudflareResponse struct {
	Success bool `json:"success"`
	Result  []struct {
		Content string `json:"content"`
		Type    string `json:"type"`
		Name    string `json:"name"`
	} `json:"result"`
	Errors []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

func NewCloudflareProvider(apiKey, zoneID string) *CloudflareProvider {
	return &CloudflareProvider{
		apiKey: apiKey,
		zoneID: zoneID,
		client: &http.Client{},
	}
}

func (p *CloudflareProvider) GetTXTRecords(ctx context.Context, domain string) ([]string, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?type=TXT&name=%s",
		p.zoneID, domain)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to query cloudflare API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cloudflare API returned status %d: %s", resp.StatusCode, string(body))
	}

	var cfResp cloudflareResponse
	if err := json.Unmarshal(body, &cfResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !cfResp.Success {
		if len(cfResp.Errors) > 0 {
			return nil, fmt.Errorf("cloudflare API error: %s", cfResp.Errors[0].Message)
		}
		return nil, fmt.Errorf("cloudflare API request failed")
	}

	var records []string
	for _, record := range cfResp.Result {
		if record.Type == "TXT" && strings.EqualFold(record.Name, domain) {
			records = append(records, record.Content)
		}
	}

	return records, nil
}
