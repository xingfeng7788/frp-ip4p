# IP4P 配置指南

IP4P (IP for Proxy) 是一个允许 frpc 通过 DNS 记录动态解析服务器地址和端口的功能。这对于服务器地址或端口可能动态变化的场景非常有用。

## 概述

IP4P 支持两种模式:

1. **lookup_text** (默认): 使用标准 DNS TXT 记录查询和 IPv6 编码
2. **api**: 使用 DNS 服务商 API 获取 TXT 记录

## 配置说明

### 模式 1: lookup_text (默认)

这是默认模式，与现有配置保持向后兼容。

```toml
serverAddr = "frp.example.com"
serverPort = 7000

[ip4p]
mode = "lookup_text"
```

或者直接省略 `[ip4p]` 配置段:

```toml
serverAddr = "frp.example.com"
serverPort = 7000
```

### 模式 2: API 模式

API 模式允许您使用 DNS 服务商的 API 来获取 TXT 记录。这比标准 DNS 查询更可靠和快速。

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

**如何获取 Cloudflare 凭证:**
1. 登录 Cloudflare 控制台
2. 进入 "我的个人资料" > "API 令牌"
3. 创建一个具有 "Zone:DNS:Read" 权限的令牌
4. 从域名概览页面获取 Zone ID

#### 腾讯云 DNSPod

```toml
serverAddr = "frp.example.com"
serverPort = 7000

[ip4p]
mode = "api"
provider = "tencent"
apiKey = "your-tencent-secret-id"
apiSecret = "your-tencent-secret-key"
```

**如何获取腾讯云凭证:**
1. 登录腾讯云控制台
2. 进入 "访问管理" > "API 密钥"
3. 创建或使用现有的 SecretId 和 SecretKey

#### 阿里云 DNS

```toml
serverAddr = "frp.example.com"
serverPort = 7000

[ip4p]
mode = "api"
provider = "alibaba"
apiKey = "your-alibaba-access-key-id"
apiSecret = "your-alibaba-access-key-secret"
```

**如何获取阿里云凭证:**
1. 登录阿里云控制台
2. 进入 "AccessKey 管理"
3. 创建或使用现有的 AccessKey ID 和 AccessKey Secret

## DNS 记录格式

### TXT 记录格式

TXT 记录应包含 base64 编码的 `地址:端口` 格式字符串。

示例:
```
# 原始内容: 192.168.1.100:7000
# Base64 编码: MTkyLjE2OC4xLjEwMDo3MDAw

frp.example.com 的 TXT 记录: "MTkyLjE2OC4xLjEwMDo3MDAw"
```

### IPv6 IP4P 编码

IP4P 还支持 IPv6 编码，其中地址和端口被编码在前缀为 `2001:0000` 的 IPv6 地址中:

```
2001:0000:xxxx:xxxx:xxxx:PPPP:AAAA:AAAA
                          ^^^^  ^^^^^^^^^
                          端口   IPv4 地址
```

示例:
- IPv4: 192.168.1.100
- 端口: 7000 (0x1B58)
- IPv6: 2001:0000:0000:0000:0000:1B58:C0A8:0164

## 查询流程

IP4P 查询按以下顺序进行:

1. **TXT 记录查询**: 尝试获取 TXT 记录 (通过 API 或标准 DNS)
2. **IPv6 IP4P 编码**: 尝试解析带有 IP4P 编码的 IPv6 地址
3. **回退**: 使用配置文件中的原始地址和端口

如果任何步骤成功，将使用解析的地址和端口。否则，将回退到下一个方法。

## 错误处理

- 如果 API 模式初始化失败 (例如，凭证无效)，将自动回退到 `lookup_text` 模式
- 如果 TXT 记录查询失败，将尝试 IPv6 IP4P 编码
- 如果所有方法都失败，将使用原始配置的地址和端口
- 所有错误都会被记录以便调试

## 使用场景

1. **动态 IP 地址**: 服务器 IP 频繁变化 (例如，家庭网络的动态 IP)
2. **负载均衡**: 通过更新 DNS 记录在多个服务器之间分配客户端
3. **故障转移**: 自动将客户端重定向到备用服务器
4. **安全性**: 通过 DNS 解析隐藏实际服务器地址
5. **多区域**: 根据 DNS 解析将客户端路由到最近的服务器

## 安全注意事项

1. **API 密钥**: 安全存储 API 密钥，避免提交到版本控制系统
2. **权限**: 为 API 令牌使用最小必需权限
3. **DNS 安全**: 考虑使用 DNSSEC 防止 DNS 欺骗
4. **速率限制**: 注意 DNS 服务商 API 的速率限制

## 故障排除

启用调试日志以查看 IP4P 解析详情:

```toml
[log]
level = "debug"
```

常见问题:

1. **API 认证失败**: 检查 API 密钥和密钥是否正确
2. **未找到 TXT 记录**: 验证 DNS 记录是否正确配置
3. **解析错误的地址**: 检查 TXT 记录格式 (必须是 base64 编码)
4. **超时**: 检查到 DNS 服务商 API 的网络连接

## 配置示例

完整的配置示例请参见 `conf/frpc_ip4p_example.toml`。

## 实现细节

### 代码结构

```
pkg/config/v1/ip4p.go          # IP4P 配置定义
pkg/util/dns/provider.go       # DNS 服务商接口
pkg/util/dns/cloudflare.go     # Cloudflare 实现
pkg/util/dns/tencent.go        # 腾讯云实现
pkg/util/dns/alibaba.go        # 阿里云实现
client/ip4p.go                 # IP4P 查询逻辑
client/connector.go            # 集成到连接器
```

### 配置参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| mode | string | 否 | 查询模式: `lookup_text` 或 `api`，默认 `lookup_text` |
| provider | string | API模式必填 | DNS 服务商: `cloudflare`、`tencent` 或 `alibaba` |
| apiKey | string | API模式必填 | API 密钥或访问密钥 ID |
| apiSecret | string | 部分必填 | API 密钥 (腾讯云和阿里云需要) |
| zoneId | string | Cloudflare必填 | Cloudflare Zone ID |

### 查询优先级

1. API 模式 TXT 记录查询 (如果配置)
2. 标准 DNS TXT 记录查询
3. IPv6 IP4P 编码查询
4. 使用原始配置值

### 向后兼容性

- 如果未配置 `[ip4p]` 段，行为与之前版本完全相同
- 默认使用 `lookup_text` 模式，保持现有功能
- 所有现有配置文件无需修改即可继续工作
