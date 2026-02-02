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
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

type TencentProvider struct {
	secretID  string
	secretKey string
	client    *http.Client
}

type tencentResponse struct {
	Response struct {
		RecordList []struct {
			Value string `json:"Value"`
			Type  string `json:"Type"`
		} `json:"RecordList"`
		Error struct {
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"Error"`
	} `json:"Response"`
}

func NewTencentProvider(secretID, secretKey string) *TencentProvider {
	return &TencentProvider{
		secretID:  secretID,
		secretKey: secretKey,
		client:    &http.Client{},
	}
}

func (p *TencentProvider) GetTXTRecords(ctx context.Context, domain string) ([]string, error) {
	// Extract domain and subdomain
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid domain format: %s", domain)
	}

	mainDomain := strings.Join(parts[len(parts)-2:], ".")
	subdomain := strings.Join(parts[:len(parts)-2], ".")
	if subdomain == "" {
		subdomain = "@"
	}

	timestamp := time.Now().Unix()
	params := map[string]string{
		"Action":     "DescribeRecordList",
		"Version":    "2021-03-23",
		"Region":     "",
		"Domain":     mainDomain,
		"Subdomain":  subdomain,
		"RecordType": "TXT",
	}

	signature := p.generateSignature(params, timestamp)

	url := "https://dnspod.tencentcloudapi.com/?" + p.buildQueryString(params)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", signature)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-TC-Action", "DescribeRecordList")
	req.Header.Set("X-TC-Version", "2021-03-23")
	req.Header.Set("X-TC-Timestamp", fmt.Sprintf("%d", timestamp))

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to query tencent API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var tcResp tencentResponse
	if err := json.Unmarshal(body, &tcResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if tcResp.Response.Error.Code != "" {
		return nil, fmt.Errorf("tencent API error: %s - %s",
			tcResp.Response.Error.Code, tcResp.Response.Error.Message)
	}

	var records []string
	for _, record := range tcResp.Response.RecordList {
		if record.Type == "TXT" {
			records = append(records, record.Value)
		}
	}

	return records, nil
}

func (p *TencentProvider) generateSignature(params map[string]string, timestamp int64) string {
	// Simplified signature generation for Tencent Cloud API v3
	canonicalQueryString := p.buildQueryString(params)
	stringToSign := fmt.Sprintf("GET\ndnspod.tencentcloudapi.com\n/\n%s", canonicalQueryString)

	h := hmac.New(sha256.New, []byte(p.secretKey))
	h.Write([]byte(stringToSign))
	signature := hex.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("TC3-HMAC-SHA256 Credential=%s, SignedHeaders=content-type;host, Signature=%s",
		p.secretID, signature)
}

func (p *TencentProvider) buildQueryString(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	return strings.Join(parts, "&")
}
