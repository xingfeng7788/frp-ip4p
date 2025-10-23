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
	"github.com/fatedier/golib/log"
	"net"
	"strconv"
)

func lookupIP4P(addr string, port int) (string, int) {
	log.Info("try txt record for addr: ", addr)
	addrs, err := net.DefaultResolver.LookupTXT(context.Background(), addr)
	log.Info("txt record result: ", addrs)
	if err == nil {
		for _, addr_64 := range addrs {
			decodeBytes, err := base64.StdEncoding.DecodeString(addr_64)
			if err != nil {
				continue
			}
			addr_s, port_s, _ := net.SplitHostPort(string(decodeBytes))
			log.Info("try to parse txt record: ", addr_s, ":", port_s)
			port_i, err := strconv.Atoi(port_s)
			if err != nil {
				continue
			}
			return addr_s, port_i
		}
	}
	log.Info("try txt record failed, try ip4p then")
	ips, err := net.DefaultResolver.LookupIP(context.Background(), "ip6", addr)
	if err == nil {
		for _, ip := range ips {
			if len(ip) == 16 {
				if ip[0] == 0x20 && ip[1] == 0x01 &&
					ip[2] == 0x00 && ip[3] == 0x00 {
					addr = net.IPv4(ip[12], ip[13], ip[14], ip[15]).String()
					port = int(ip[10])<<8 | int(ip[11])
					return addr, port
				}
			}
		}
	}
	log.Info("try ip4p record failed, try normal mode")
	log.Info("try to resolve addr: ", addr)
	log.Info("try to resolve port: ", port)
	return addr, port
}
