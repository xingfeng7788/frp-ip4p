# IP4P Configuration Guide

IP4P (IP for Proxy) is a feature that allows frpc to dynamically resolve server addresses and ports through DNS records. This is useful for scenarios where the server address or port may change dynamically.

## Overview

IP4P supports two modes:

1. **lookup_text** (default): Uses standard DNS TXT record lookup and IPv6 encoding
2. **api**: Uses DNS provider APIs to retrieve TXT records

## Configuration

### Mode 1: lookup_text (Default)

This is the default mode and maintains backward compatibility with existing configurations.

```toml
serverAddr = "frp.example.com"
serverPort = 7000

[ip4p]
mode = "lookup_text"
```

Or simply omit the `[ip4p]` section entirely:

```toml
serverAddr = "frp.example.com"
serverPort = 7000
```

### Mode 2: API Mode

API mode allows you to use DNS provider APIs to retrieve TXT records. This can be more reliable and faster than standard DNS lookups.

#### Cloudflare

```toml
serverAddr = "frp.example.com"
serverPort = 7000

[ip4p]
mode = "api"
provider = "cloudflare"
apiKey = "your-cloudflare-api-token"
zoneId = "your-cloudflare-zone-id"
```

**How to get Cloudflare credentials:**
1. Log in to your Cloudflare dashboard
2. Go to "My Profile" > "API Tokens"
3. Create a token with "Zone:DNS:Read" permissions
4. 通过创建或则编辑一个dns解析打开开发者工具 找到类似 https://dash.cloudflare.com/api/v4/zones/***/dns_records 调用的接口中找到zone_id


#### Tencent Cloud DNSPod

```toml
serverAddr = "frp.example.com"
serverPort = 7000

[ip4p]
mode = "api"
provider = "tencent"
apiKey = "your-tencent-secret-id"
apiSecret = "your-tencent-secret-key"
```

**How to get Tencent Cloud credentials:**
1. Log in to Tencent Cloud console
2. Go to "Access Management" > "API Keys"
3. Create or use existing SecretId and SecretKey

#### Alibaba Cloud DNS

```toml
serverAddr = "frp.example.com"
serverPort = 7000

[ip4p]
mode = "api"
provider = "alibaba"
apiKey = "your-alibaba-access-key-id"
apiSecret = "your-alibaba-access-key-secret"
```

**How to get Alibaba Cloud credentials:**
1. Log in to Alibaba Cloud console
2. Go to "AccessKey Management"
3. Create or use existing AccessKey ID and AccessKey Secret

## DNS Record Format

### TXT Record Format

The TXT record should contain a base64-encoded string of the format `address:port`.

Example:
```
# Original: 192.168.1.100:7000
# Base64 encoded: MTkyLjE2OC4xLjEwMDo3MDAw

TXT record for frp.example.com: "MTkyLjE2OC4xLjEwMDo3MDAw"
```

### IPv6 IP4P Encoding

IP4P also supports IPv6 encoding where the address and port are encoded in an IPv6 address with the prefix `2001:0000`:

```
2001:0000:xxxx:xxxx:xxxx:PPPP:AAAA:AAAA
                          ^^^^  ^^^^^^^^^
                          Port  IPv4 Address
```

Example:
- IPv4: 192.168.1.100
- Port: 7000 (0x1B58)
- IPv6: 2001:0000:0000:0000:0000:1B58:C0A8:0164

## Lookup Process

The IP4P lookup follows this order:

1. **TXT Record Lookup**: Try to get TXT records (via API or standard DNS)
2. **IPv6 IP4P Encoding**: Try to resolve IPv6 address with IP4P encoding
3. **Fallback**: Use the original address and port from configuration

If any step succeeds, the resolved address and port are used. Otherwise, it falls back to the next method.

## Error Handling

- If API mode fails to initialize (e.g., invalid credentials), it automatically falls back to `lookup_text` mode
- If TXT record lookup fails, it tries IPv6 IP4P encoding
- If all methods fail, the original configured address and port are used
- All errors are logged for debugging purposes

## Use Cases

1. **Dynamic IP Addresses**: Server IP changes frequently (e.g., home network with dynamic IP)
2. **Load Balancing**: Distribute clients across multiple servers by updating DNS records
3. **Failover**: Automatically redirect clients to backup servers
4. **Security**: Hide actual server address behind DNS resolution
5. **Multi-region**: Route clients to nearest server based on DNS resolution

## Security Considerations

1. **API Keys**: Store API keys securely, avoid committing them to version control
2. **Permissions**: Use minimal required permissions for API tokens
3. **DNS Security**: Consider using DNSSEC to prevent DNS spoofing
4. **Rate Limiting**: Be aware of DNS provider API rate limits

## Troubleshooting

Enable debug logging to see IP4P resolution details:

```toml
[log]
level = "debug"
```

Common issues:

1. **API authentication fails**: Check API key and secret are correct
2. **No TXT records found**: Verify DNS records are properly configured
3. **Wrong address resolved**: Check TXT record format (must be base64 encoded)
4. **Timeout**: Check network connectivity to DNS provider API

## Examples

See `conf/frpc_ip4p_example.toml` for complete configuration examples.
