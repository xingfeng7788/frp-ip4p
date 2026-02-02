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

package client

import (
	"context"
	"testing"

	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/stretchr/testify/assert"
)

func TestNewIP4PLookup(t *testing.T) {
	tests := []struct {
		name    string
		config  *v1.IP4PConfig
		wantErr bool
	}{
		{
			name: "default lookup_text mode",
			config: &v1.IP4PConfig{
				Mode: v1.IP4PModeLookupText,
			},
			wantErr: false,
		},
		{
			name: "api mode without provider",
			config: &v1.IP4PConfig{
				Mode: v1.IP4PModeAPI,
			},
			wantErr: false, // Should fallback to lookup_text
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lookup, err := NewIP4PLookup(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, lookup)
			}
		})
	}
}

func TestIP4PLookup_Fallback(t *testing.T) {
	config := &v1.IP4PConfig{
		Mode: v1.IP4PModeLookupText,
	}
	lookup, err := NewIP4PLookup(config)
	assert.NoError(t, err)

	// Test with non-existent domain, should return original values
	addr, port := lookup.Lookup(context.Background(), "non-existent-domain-12345.example", 7000)
	assert.Equal(t, "non-existent-domain-12345.example", addr)
	assert.Equal(t, 7000, port)
}

func TestIP4PLookup_ParseTXTRecords(t *testing.T) {
	config := &v1.IP4PConfig{
		Mode: v1.IP4PModeLookupText,
	}
	lookup, err := NewIP4PLookup(config)
	assert.NoError(t, err)

	tests := []struct {
		name         string
		records      []string
		defaultAddr  string
		defaultPort  int
		expectedAddr string
		expectedPort int
	}{
		{
			name:         "valid base64 encoded record",
			records:      []string{"MTkyLjE2OC4xLjEwMDo3MDAw"}, // 192.168.1.100:7000
			defaultAddr:  "example.com",
			defaultPort:  8000,
			expectedAddr: "192.168.1.100",
			expectedPort: 7000,
		},
		{
			name:         "invalid record",
			records:      []string{"invalid-record"},
			defaultAddr:  "example.com",
			defaultPort:  8000,
			expectedAddr: "example.com",
			expectedPort: 8000,
		},
		{
			name:         "empty records",
			records:      []string{},
			defaultAddr:  "example.com",
			defaultPort:  8000,
			expectedAddr: "example.com",
			expectedPort: 8000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, port := lookup.parseTXTRecords(tt.records, tt.defaultAddr, tt.defaultPort)
			assert.Equal(t, tt.expectedAddr, addr)
			assert.Equal(t, tt.expectedPort, port)
		})
	}
}

func TestLookupIP4P_BackwardCompatibility(t *testing.T) {
	// Test the backward compatible function
	addr, port := lookupIP4P("non-existent-domain-12345.example", 7000)
	assert.Equal(t, "non-existent-domain-12345.example", addr)
	assert.Equal(t, 7000, port)
}
