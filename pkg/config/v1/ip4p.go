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

package v1

import "github.com/fatedier/frp/pkg/util/util"

type IP4PMode string

const (
	IP4PModeLookupText IP4PMode = "lookup_text"
	IP4PModeAPI        IP4PMode = "api"
)

type DNSProvider string

const (
	DNSProviderCloudflare DNSProvider = "cloudflare"
	DNSProviderTencent    DNSProvider = "tencent"
	DNSProviderAlibaba    DNSProvider = "alibaba"
)

type IP4PConfig struct {
	// Mode specifies the IP4P lookup mode.
	// Valid values are "lookup_text" (default) and "api".
	Mode IP4PMode `json:"mode,omitempty"`

	// Provider specifies the DNS provider when using API mode.
	// Valid values are "cloudflare", "tencent", "alibaba".
	Provider DNSProvider `json:"provider,omitempty"`

	// APIKey is the API key for the DNS provider.
	APIKey string `json:"apiKey,omitempty"`

	// APISecret is the API secret for the DNS provider (required for some providers).
	APISecret string `json:"apiSecret,omitempty"`

	// ZoneID is the zone ID for Cloudflare.
	ZoneID string `json:"zoneId,omitempty"`
}

func (c *IP4PConfig) Complete() {
	c.Mode = util.EmptyOr(c.Mode, IP4PModeLookupText)
}
