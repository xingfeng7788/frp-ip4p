# IP4P 功能优化说明

## 概述

本次优化对 frp 的 IP4P (IP for Proxy) 功能进行了全面升级，增加了通过 DNS 服务商 API 获取 TXT 记录的能力，支持 Cloudflare、腾讯云和阿里云三大主流 DNS 服务商。

## 主要改进

### 1. 配置化的 IP4P 模式

新增 `ip4p` 配置段，支持两种模式：

- **lookup_text** (默认): 保持原有的 DNS TXT 记录查询和 IPv6 编码方式
- **api**: 通过 DNS 服务商 API 获取 TXT 记录

### 2. 多 DNS 服务商支持

实现了三大主流 DNS 服务商的 API 集成：

- **Cloudflare**: 使用 API Token 和 Zone ID
- **腾讯云 DNSPod**: 使用 SecretId 和 SecretKey
- **阿里云 DNS**: 使用 AccessKey ID 和 AccessKey Secret

### 3. 优化的查询流程

新的查询流程更加健壮和灵活：

1. 首先尝试配置的模式（API 或 lookup_text）
2. 如果失败，尝试 IPv6 IP4P 编码
3. 最后回退到原始配置的地址和端口

### 4. 向后兼容

- 完全兼容现有配置，无需修改即可继续使用
- 保留原有的 `lookupIP4P` 函数接口
- 默认行为与之前版本完全一致

## 文件结构

### 新增文件

```
pkg/config/v1/ip4p.go              # IP4P 配置定义
pkg/util/dns/provider.go           # DNS 服务商接口
pkg/util/dns/cloudflare.go         # Cloudflare API 实现
pkg/util/dns/tencent.go            # 腾讯云 DNSPod API 实现
pkg/util/dns/alibaba.go            # 阿里云 DNS API 实现
client/ip4p_test.go                # 单元测试
conf/frpc_ip4p_example.toml        # IP4P 配置示例
doc/ip4p.md                        # 英文文档
doc/ip4p_zh.md                     # 中文文档
```

### 修改文件

```
pkg/config/v1/client.go            # 添加 IP4P 配置字段
client/ip4p.go                     # 重构为面向对象设计
client/connector.go                # 集成 IP4P 查询
conf/frpc_full_example.toml        # 添加 IP4P 配置示例
```

## 配置示例

### 基本配置（默认模式）

```toml
serverAddr = "frp.example.com"
serverPort = 7000

[ip4p]
mode = "lookup_text"
```

### Cloudflare API 模式

```toml
serverAddr = "frp.example.com"
serverPort = 7000

[ip4p]
mode = "api"
provider = "cloudflare"
apiKey = "your-cloudflare-api-token"
zoneId = "your-cloudflare-zone-id"
```

### 腾讯云 API 模式

```toml
serverAddr = "frp.example.com"
serverPort = 7000

[ip4p]
mode = "api"
provider = "tencent"
apiKey = "your-tencent-secret-id"
apiSecret = "your-tencent-secret-key"
```

### 阿里云 API 模式

```toml
serverAddr = "frp.example.com"
serverPort = 7000

[ip4p]
mode = "api"
provider = "alibaba"
apiKey = "your-alibaba-access-key-id"
apiSecret = "your-alibaba-access-key-secret"
```

## 技术实现

### 1. 配置层 (pkg/config/v1/ip4p.go)

定义了 IP4P 的配置结构：

```go
type IP4PConfig struct {
    Mode      IP4PMode    // lookup_text 或 api
    Provider  DNSProvider // cloudflare, tencent, alibaba
    APIKey    string      // API 密钥
    APISecret string      // API 密钥（部分服务商需要）
    ZoneID    string      // Zone ID（Cloudflare 需要）
}
```

### 2. DNS 服务商层 (pkg/util/dns/)

实现了统一的 DNS 服务商接口：

```go
type Provider interface {
    GetTXTRecords(ctx context.Context, domain string) ([]string, error)
}
```

每个服务商都实现了该接口：

- **CloudflareProvider**: 使用 Cloudflare API v4
- **TencentProvider**: 使用腾讯云 DNSPod API v3
- **AlibabaProvider**: 使用阿里云 DNS API

### 3. 查询层 (client/ip4p.go)

重构为面向对象设计：

```go
type IP4PLookup struct {
    config   *v1.IP4PConfig
    provider dns.Provider
}

func (l *IP4PLookup) Lookup(ctx context.Context, addr string, port int) (string, int)
```

查询流程：

1. 根据配置模式选择查询方法
2. API 模式：调用 DNS 服务商 API
3. lookup_text 模式：使用标准 DNS 查询
4. 解析 TXT 记录（base64 解码）
5. 尝试 IPv6 IP4P 编码
6. 回退到原始值

### 4. 集成层 (client/connector.go)

在连接器中集成 IP4P 查询：

```go
type defaultConnectorImpl struct {
    ctx        context.Context
    cfg        *v1.ClientCommonConfig
    ip4pLookup *IP4PLookup  // 新增
    // ...
}
```

在建立连接时使用 IP4P 查询：

```go
serverAddr, port := c.cfg.ServerAddr, c.cfg.ServerPort
if c.ip4pLookup != nil {
    serverAddr, port = c.ip4pLookup.Lookup(c.ctx, c.cfg.ServerAddr, c.cfg.ServerPort)
}
```

## DNS 记录格式

### TXT 记录

TXT 记录应包含 base64 编码的 `地址:端口` 字符串：

```
原始: 192.168.1.100:7000
Base64: MTkyLjE2OC4xLjEwMDo3MDAw
```

### IPv6 IP4P 编码

IPv6 地址格式：`2001:0000:xxxx:xxxx:xxxx:PPPP:AAAA:AAAA`

- `PPPP`: 端口号（16位）
- `AAAA:AAAA`: IPv4 地址（32位）

示例：
```
IPv4: 192.168.1.100
Port: 7000 (0x1B58)
IPv6: 2001:0000:0000:0000:0000:1B58:C0A8:0164
```

## 错误处理

### 自动降级

- API 模式初始化失败 → 降级到 lookup_text 模式
- API 查询失败 → 尝试标准 DNS 查询
- TXT 记录查询失败 → 尝试 IPv6 编码
- 所有方法失败 → 使用原始配置值

### 日志记录

所有关键步骤都有详细的日志记录：

```
IP4P: using API mode with provider: cloudflare
IP4P: API returned TXT records: [...]
IP4P: parsed TXT record: 192.168.1.100:7000
IP4P: found IPv6 encoding, resolved to: 192.168.1.100:7000
IP4P: all methods failed, using original values: example.com:7000
```

## 测试

### 单元测试 (client/ip4p_test.go)

包含以下测试用例：

1. `TestNewIP4PLookup`: 测试 IP4PLookup 创建
2. `TestIP4PLookup_Fallback`: 测试回退机制
3. `TestIP4PLookup_ParseTXTRecords`: 测试 TXT 记录解析
4. `TestLookupIP4P_BackwardCompatibility`: 测试向后兼容性

运行测试：

```bash
go test ./client -v -run TestIP4P
```

## 使用场景

1. **动态 IP**: 家庭宽带等动态 IP 环境
2. **负载均衡**: 通过 DNS 记录分配流量
3. **故障转移**: 自动切换到备用服务器
4. **多区域部署**: 根据地理位置路由
5. **安全隐藏**: 隐藏真实服务器地址

## 安全建议

1. **API 密钥管理**
   - 不要将 API 密钥提交到版本控制
   - 使用环境变量或密钥管理服务
   - 定期轮换 API 密钥

2. **权限最小化**
   - Cloudflare: 只授予 Zone:DNS:Read 权限
   - 腾讯云: 只授予 DNS 读取权限
   - 阿里云: 只授予 DNS 读取权限

3. **网络安全**
   - 使用 HTTPS 连接 API
   - 考虑使用 DNSSEC
   - 监控 API 调用频率

## 性能考虑

1. **缓存**: 考虑在应用层添加 DNS 结果缓存
2. **超时**: API 调用有合理的超时设置
3. **重试**: 失败时自动降级，不阻塞连接
4. **并发**: 支持并发查询，不影响性能

## 未来改进

1. **更多 DNS 服务商**: 可以添加更多服务商支持
2. **缓存机制**: 添加 DNS 查询结果缓存
3. **健康检查**: 定期验证解析结果的可用性
4. **监控指标**: 添加 Prometheus 指标
5. **配置热更新**: 支持动态更新 IP4P 配置

## 文档

- 英文文档: `doc/ip4p.md`
- 中文文档: `doc/ip4p_zh.md`
- 配置示例: `conf/frpc_ip4p_example.toml`
- 完整示例: `conf/frpc_full_example.toml`

## 兼容性

- **Go 版本**: 需要 Go 1.24.0+
- **frp 版本**: 基于当前主分支
- **配置格式**: TOML
- **向后兼容**: 完全兼容现有配置

## 总结

本次优化为 frp 的 IP4P 功能带来了以下改进：

✅ 支持通过 DNS 服务商 API 获取记录  
✅ 支持 Cloudflare、腾讯云、阿里云  
✅ 配置化的模式选择  
✅ 健壮的错误处理和降级机制  
✅ 完整的文档和示例  
✅ 单元测试覆盖  
✅ 向后兼容  

这些改进使得 IP4P 功能更加灵活、可靠和易用，适合各种动态 IP 场景。
