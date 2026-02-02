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
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type AlibabaProvider struct {
	accessKeyID     string
	accessKeySecret string
	client          *http.Client
}

type alibabaResponse struct {
	DomainRecords struct {
		Record []struct {
			Value string `json:"Value"`
			Type  string `json:"Type"`
		} `json:"Record"`
	} `json:"DomainRecords"`
	Code      string `json:"Code"`
	Message   string `json:"Message"`
	RequestID string `json:"RequestId"`
}

func NewAlibabaProvider(accessKeyID, accessKeySecret string) *AlibabaProvider {
	return &AlibabaProvider{
		accessKeyID:     accessKeyID,
		accessKeySecret: accessKeySecret,
		client:          &http.Client{},
	}
}

func (p *AlibabaProvider) GetTXTRecords(ctx context.Context, domain string) ([]string, error) {
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

	params := map[string]string{
		"Action":           "DescribeDomainRecords",
		"Version":          "2015-01-09",
		"AccessKeyId":      p.accessKeyID,
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"SignatureNonce":   uuid.New().String(),
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"Format":           "JSON",
		"DomainName":       mainDomain,
		"RRKeyWord":        subdomain,
		"Type":             "TXT",
	}

	signature := p.generateSignature(params)
	params["Signature"] = signature

	queryString := p.buildQueryString(params)
	apiURL := "https://alidns.aliyuncs.com/?" + queryString

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to query alibaba API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var aliResp alibabaResponse
	if err := json.Unmarshal(body, &aliResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if aliResp.Code != "" {
		return nil, fmt.Errorf("alibaba API error: %s - %s", aliResp.Code, aliResp.Message)
	}

	var records []string
	for _, record := range aliResp.DomainRecords.Record {
		if record.Type == "TXT" {
			records = append(records, record.Value)
		}
	}

	return records, nil
}

func (p *AlibabaProvider) generateSignature(params map[string]string) string {
	// Sort parameters
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build canonical query string
	var parts []string
	for _, k := range keys {
		parts = append(parts, url.QueryEscape(k)+"="+url.QueryEscape(params[k]))
	}
	canonicalQueryString := strings.Join(parts, "&")

	// Build string to sign
	stringToSign := "GET&" + url.QueryEscape("/") + "&" + url.QueryEscape(canonicalQueryString)

	// Calculate signature
	h := hmac.New(sha1.New, []byte(p.accessKeySecret+"&"))
	h.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature
}

func (p *AlibabaProvider) buildQueryString(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		parts = append(parts, url.QueryEscape(k)+"="+url.QueryEscape(params[k]))
	}
	return strings.Join(parts, "&")
}
