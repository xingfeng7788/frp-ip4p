# FRP 配置文件格式说明

## 问题说明

如果你遇到 `json: unknown field "auth"` 错误，这是因为配置文件格式不正确。

## TOML 配置格式

FRP 使用 TOML 格式的配置文件。TOML 支持两种表示嵌套结构的方式：

### ❌ 错误格式（混合使用）

```toml
serverAddr = "example.com"

[ip4p]
mode = "api"
provider = "cloudflare"

auth.method = "token"      # ❌ 错误：不能在段外使用点号
auth.token = "xxx"         # ❌ 错误：不能在段外使用点号

log.to = "./frpc.log"      # ❌ 错误：不能在段外使用点号
log.level = "info"         # ❌ 错误：不能在段外使用点号
```

### ✅ 正确格式 1：使用段（推荐）

```toml
serverAddr = "example.com"

[ip4p]
mode = "api"
provider = "cloudflare"
apiKey = "your-api-key"
zoneId = "your-zone-id"

[auth]
method = "token"
token = "xxx"

[log]
to = "./frpc.log"
level = "info"
maxDays = 3
disablePrintColor = true

[[visitors]]
name = "visitor_name"
type = "stcp"
serverName = "server_name"
secretKey = "xxx"
bindAddr = "127.0.0.1"
bindPort = 20001
```

### ✅ 正确格式 2：完全使用点号（不推荐，但有效）

```toml
serverAddr = "example.com"

ip4p.mode = "api"
ip4p.provider = "cloudflare"
ip4p.apiKey = "your-api-key"
ip4p.zoneId = "your-zone-id"

auth.method = "token"
auth.token = "xxx"

log.to = "./frpc.log"
log.level = "info"
log.maxDays = 3
log.disablePrintColor = true

[[visitors]]
name = "visitor_name"
type = "stcp"
serverName = "server_name"
secretKey = "xxx"
bindAddr = "127.0.0.1"
bindPort = 20001
```

## 配置段说明

### 基本配置

```toml
# 服务器地址和端口
serverAddr = "ip4p.tianai.edu.kg"
serverPort = 7000  # 可选，默认 7000
```

### IP4P 配置（新增功能）

```toml
[ip4p]
mode = "api"                    # 模式：lookup_text 或 api
provider = "cloudflare"         # DNS 服务商：cloudflare, tencent, alibaba
apiKey = "your-api-key"         # API 密钥
apiSecret = "your-api-secret"   # API 密钥（腾讯云和阿里云需要）
zoneId = "your-zone-id"         # Zone ID（Cloudflare 需要）
```

### 认证配置

```toml
[auth]
method = "token"
token = "your-token"
```

### 日志配置

```toml
[log]
to = "./frpc.log"           # 日志文件路径，或 "console" 输出到控制台
level = "info"              # 日志级别：trace, debug, info, warn, error
maxDays = 3                 # 日志保留天数
disablePrintColor = true    # 禁用彩色输出
```

### 传输配置

```toml
[transport]
protocol = "tcp"            # 协议：tcp, kcp, quic, websocket, wss
poolCount = 5               # 连接池大小
tcpMux = true              # 是否启用 TCP 多路复用

[transport.tls]
enable = true              # 是否启用 TLS
```

### 访问者配置

```toml
[[visitors]]
name = "visitor_ssh"
type = "stcp"
serverName = "frp_ssh"
secretKey = "your-secret-key"
bindAddr = "127.0.0.1"
bindPort = 20001
```

## 完整配置示例

### 示例 1：使用 IP4P API 模式（Cloudflare）

```toml
serverAddr = "ip4p.tianai.edu.kg"

[ip4p]
mode = "api"
provider = "cloudflare"
apiKey = "your-cloudflare-api-token"
zoneId = "your-cloudflare-zone-id"

[auth]
method = "token"
token = "your-auth-token"

[log]
to = "./frpc.log"
level = "info"
maxDays = 3
disablePrintColor = true

[[visitors]]
name = "visitor_ssh"
type = "stcp"
serverName = "frp_ssh"
secretKey = "your-secret-key"
bindAddr = "127.0.0.1"
bindPort = 20001
```

### 示例 2：使用 IP4P 默认模式

```toml
serverAddr = "ip4p.tianai.edu.kg"

# IP4P 使用默认的 lookup_text 模式，可以省略配置
# [ip4p]
# mode = "lookup_text"

[auth]
method = "token"
token = "your-auth-token"

[log]
to = "./frpc.log"
level = "info"

[[visitors]]
name = "visitor_ssh"
type = "stcp"
serverName = "frp_ssh"
secretKey = "your-secret-key"
bindAddr = "127.0.0.1"
bindPort = 20001
```

### 示例 3：使用腾讯云 DNS API

```toml
serverAddr = "ip4p.tianai.edu.kg"

[ip4p]
mode = "api"
provider = "tencent"
apiKey = "your-tencent-secret-id"
apiSecret = "your-tencent-secret-key"

[auth]
method = "token"
token = "your-auth-token"

[log]
to = "./frpc.log"
level = "info"

[[visitors]]
name = "visitor_ssh"
type = "stcp"
serverName = "frp_ssh"
secretKey = "your-secret-key"
bindAddr = "127.0.0.1"
bindPort = 20001
```

## 配置文件位置

默认配置文件位置：
- Linux/macOS: `./frpc.toml` 或 `/etc/frp/frpc.toml`
- Windows: `.\frpc.toml`

使用自定义配置文件：
```bash
./frpc -c /path/to/your/frpc.toml
```

## 验证配置文件

使用 `verify` 命令验证配置文件格式：

```bash
./frpc verify -c frpc.toml
```

## 常见错误

### 1. `json: unknown field "auth"`

**原因**：混合使用了段格式和点号格式

**解决**：统一使用段格式（推荐）

```toml
# ❌ 错误
[ip4p]
mode = "api"

auth.method = "token"  # 混合格式

# ✅ 正确
[ip4p]
mode = "api"

[auth]
method = "token"
```

### 2. `unknown field "xxx"`

**原因**：字段名拼写错误或不存在

**解决**：检查字段名是否正确，参考官方文档

### 3. TOML 解析错误

**原因**：TOML 语法错误

**解决**：
- 检查引号是否配对
- 检查等号两边是否有空格
- 检查段名是否正确（`[section]` 或 `[[array]]`）

## 配置迁移

如果你有旧的配置文件使用点号格式，需要转换为段格式：

### 转换前（点号格式）

```toml
auth.method = "token"
auth.token = "xxx"
log.to = "./frpc.log"
log.level = "info"
```

### 转换后（段格式）

```toml
[auth]
method = "token"
token = "xxx"

[log]
to = "./frpc.log"
level = "info"
```

## 工具推荐

- **TOML 在线验证器**: https://www.toml-lint.com/
- **TOML 格式化工具**: https://toolkit.site/format.html?type=toml

## 参考文档

- FRP 官方文档: https://github.com/fatedier/frp
- TOML 规范: https://toml.io/
- IP4P 配置文档: `doc/ip4p_zh.md`
