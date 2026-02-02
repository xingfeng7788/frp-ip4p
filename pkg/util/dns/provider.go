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
	"fmt"

	v1 "github.com/fatedier/frp/pkg/config/v1"
)

// Provider defines the interface for DNS providers
type Provider interface {
	// GetTXTRecords retrieves TXT records for the given domain
	GetTXTRecords(ctx context.Context, domain string) ([]string, error)
}

// NewProvider creates a new DNS provider based on the configuration
func NewProvider(cfg *v1.IP4PConfig) (Provider, error) {
	switch cfg.Provider {
	case v1.DNSProviderCloudflare:
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("cloudflare API key is required")
		}
		if cfg.ZoneID == "" {
			return nil, fmt.Errorf("cloudflare zone ID is required")
		}
		return NewCloudflareProvider(cfg.APIKey, cfg.ZoneID), nil
	case v1.DNSProviderTencent:
		if cfg.APIKey == "" || cfg.APISecret == "" {
			return nil, fmt.Errorf("tencent API key and secret are required")
		}
		return NewTencentProvider(cfg.APIKey, cfg.APISecret), nil
	case v1.DNSProviderAlibaba:
		if cfg.APIKey == "" || cfg.APISecret == "" {
			return nil, fmt.Errorf("alibaba API key and secret are required")
		}
		return NewAlibabaProvider(cfg.APIKey, cfg.APISecret), nil
	default:
		return nil, fmt.Errorf("unsupported DNS provider: %s", cfg.Provider)
	}
}
