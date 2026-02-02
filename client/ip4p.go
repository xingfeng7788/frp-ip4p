// Copyright 2022 hev, r@hev.cc
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
	"encoding/base64"
	"net"
	"strconv"

	"github.com/fatedier/golib/log"

	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/util/dns"
)

// IP4PLookup handles IP4P address resolution with configurable modes
type IP4PLookup struct {
	config   *v1.IP4PConfig
	provider dns.Provider
}

// NewIP4PLookup creates a new IP4P lookup instance
func NewIP4PLookup(config *v1.IP4PConfig) (*IP4PLookup, error) {
	lookup := &IP4PLookup{
		config: config,
	}

	// Initialize DNS provider if API mode is enabled
	if config.Mode == v1.IP4PModeAPI {
		provider, err := dns.NewProvider(config)
		if err != nil {
			log.Warn("failed to initialize DNS provider, falling back to lookup_text mode: ", err)
			// Fallback to lookup_text mode
			config.Mode = v1.IP4PModeLookupText
		} else {
			lookup.provider = provider
		}
	}

	return lookup, nil
}

// Lookup performs IP4P address resolution
func (l *IP4PLookup) Lookup(ctx context.Context, addr string, port int) (string, int) {
	switch l.config.Mode {
	case v1.IP4PModeAPI:
		return l.lookupViaAPI(ctx, addr, port)
	default:
		return l.lookupViaText(ctx, addr, port)
	}
}

// lookupViaAPI uses DNS provider API to get TXT records
func (l *IP4PLookup) lookupViaAPI(ctx context.Context, addr string, port int) (string, int) {
	log.Info("IP4P: using API mode with provider: ", l.config.Provider)

	if l.provider == nil {
		log.Warn("IP4P: provider not initialized, falling back to text mode")
		return l.lookupViaText(ctx, addr, port)
	}

	records, err := l.provider.GetTXTRecords(ctx, addr)
	if err != nil {
		log.Warn("IP4P: failed to get TXT records via API: ", err)
		return l.fallbackLookup(ctx, addr, port)
	}

	log.Info("IP4P: API returned TXT records: ", records)
	resolvedAddr, resolvedPort := l.parseTXTRecords(records, addr, port)
	if resolvedAddr != addr || resolvedPort != port {
		return resolvedAddr, resolvedPort
	}

	return l.fallbackLookup(ctx, addr, port)
}

// lookupViaText uses standard DNS TXT record lookup
func (l *IP4PLookup) lookupViaText(ctx context.Context, addr string, port int) (string, int) {
	log.Info("IP4P: trying TXT record lookup for addr: ", addr)

	records, err := net.DefaultResolver.LookupTXT(ctx, addr)
	if err == nil {
		log.Info("IP4P: TXT record result: ", records)
		resolvedAddr, resolvedPort := l.parseTXTRecords(records, addr, port)
		if resolvedAddr != addr || resolvedPort != port {
			return resolvedAddr, resolvedPort
		}
	} else {
		log.Info("IP4P: TXT record lookup failed: ", err)
	}

	return l.fallbackLookup(ctx, addr, port)
}

// parseTXTRecords parses TXT records to extract address and port
func (l *IP4PLookup) parseTXTRecords(records []string, defaultAddr string, defaultPort int) (string, int) {
	for _, record := range records {
		decodeBytes, err := base64.StdEncoding.DecodeString(record)
		if err != nil {
			continue
		}

		addrStr, portStr, err := net.SplitHostPort(string(decodeBytes))
		if err != nil {
			continue
		}

		log.Info("IP4P: parsed TXT record: ", addrStr, ":", portStr)

		portInt, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}

		return addrStr, portInt
	}

	return defaultAddr, defaultPort
}

// fallbackLookup tries IP4P IPv6 encoding and then returns original values
func (l *IP4PLookup) fallbackLookup(ctx context.Context, addr string, port int) (string, int) {
	log.Info("IP4P: trying IPv6 IP4P encoding")

	ips, err := net.DefaultResolver.LookupIP(ctx, "ip6", addr)
	if err == nil {
		for _, ip := range ips {
			if len(ip) == 16 {
				// Check for IP4P encoding: 2001:0000:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx
				if ip[0] == 0x20 && ip[1] == 0x01 &&
					ip[2] == 0x00 && ip[3] == 0x00 {
					resolvedAddr := net.IPv4(ip[12], ip[13], ip[14], ip[15]).String()
					resolvedPort := int(ip[10])<<8 | int(ip[11])
					log.Info("IP4P: found IPv6 encoding, resolved to: ", resolvedAddr, ":", resolvedPort)
					return resolvedAddr, resolvedPort
				}
			}
		}
	} else {
		log.Info("IP4P: IPv6 lookup failed: ", err)
	}

	log.Info("IP4P: all methods failed, using original values: ", addr, ":", port)
	return addr, port
}

// lookupIP4P is a convenience function that maintains backward compatibility
func lookupIP4P(addr string, port int) (string, int) {
	// Use default lookup_text mode for backward compatibility
	config := &v1.IP4PConfig{
		Mode: v1.IP4PModeLookupText,
	}
	lookup, _ := NewIP4PLookup(config)
	return lookup.Lookup(context.Background(), addr, port)
}
