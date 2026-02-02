# frp

[![Build Status](https://circleci.com/gh/fatedier/frp.svg?style=shield)](https://circleci.com/gh/fatedier/frp)
[![GitHub release](https://img.shields.io/github/tag/fatedier/frp.svg?label=release)](https://github.com/fatedier/frp/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/fatedier/frp)](https://goreportcard.com/report/github.com/fatedier/frp)
[![GitHub Releases Stats](https://img.shields.io/github/downloads/fatedier/frp/total.svg?logo=github)](https://somsubhra.github.io/github-release-stats/?username=fatedier&repository=frp)

[README](README.md) | [中文文档](README)

frp 是一个专注于内网穿透的高性能的反向代理应用，支持 TCP、UDP、HTTP、HTTPS 等多种协议，且支持 P2P 通信。可以将内网服务以安全、便捷的方式通过具有公网 IP 节点的中转暴露到公网。

## Sponsors

frp 是一个完全开源的项目，我们的开发工作完全依靠赞助者们的支持。如果你愿意加入他们的行列，请考虑 [赞助 frp 的开发](https://github.com/sponsors/fatedier)。

<h3 align="center">Gold Sponsors</h3>
<!--gold sponsors start-->
<p align="center">
  <a href="https://www.recall.ai/?utm_source=github&utm_medium=sponsorship&utm_campaign=fatedier-frp" target="_blank">
    <b>Recall.ai - API for meeting recordings</b><br>
    <br>
    <sup>If you're looking for a meeting recording API, consider checking out Recall.ai, an API that records Zoom, Google Meet, Microsoft Teams, in-person meetings, and more.</sup>
  </a>
</p>
<p align="center">
  <a href="https://go.warp.dev/frp" target="_blank">
    <img width="360px" src="https://raw.githubusercontent.com/warpdotdev/brand-assets/refs/heads/main/Github/Sponsor/Warp-Github-LG-01.png">
    <br>
    <b>Warp, built for collaborating with AI Agents</b>
    <br>
        <sub>Available for macOS, Linux and Windows</sub>
  </a>
</p>
<p align="center">
  <a href="https://jb.gg/frp" target="_blank">
    <img width="420px" src="https://raw.githubusercontent.com/fatedier/frp/dev/doc/pic/sponsor_jetbrains.jpg">
        <br>
        <b>The complete IDE crafted for professional Go developers</b>
  </a>
</p>
<p align="center">
  <a href="https://github.com/daytonaio/daytona" target="_blank">
    <img width="420px" src="https://raw.githubusercontent.com/fatedier/frp/dev/doc/pic/sponsor_daytona.png">
        <br>
        <b>Secure and Elastic Infrastructure for Running Your AI-Generated Code</b>
  </a>
</p>
<p align="center">
  <a href="https://github.com/beclab/Olares" target="_blank">
    <img width="420px" src="https://raw.githubusercontent.com/fatedier/frp/dev/doc/pic/sponsor_olares.jpeg">
        <br>
        <b>The sovereign cloud that puts you in control</b>
        <br>
        <sub>An open source, self-hosted alternative to public clouds, built for data ownership and privacy</sub>
  </a>
</p>
<!--gold sponsors end-->

## 为什么使用 frp ？

通过在具有公网 IP 节点上部署 frp 服务端，可以轻松地将内网服务穿透到公网，同时提供诸多专业的功能特性，这包括：

* 客户端服务端通信支持 TCP、QUIC、KCP 以及 Websocket 等多种协议。
* 采用 TCP 连接流式复用，在单个连接间承载更多请求，节省连接建立时间，降低请求延迟。
* 代理组间的负载均衡。
* 端口复用，多个服务通过同一个服务端端口暴露。
* 支持 P2P 通信，流量不经过服务器中转，充分利用带宽资源。
* 多个原生支持的客户端插件（静态文件查看，HTTPS/HTTP 协议转换，HTTP、SOCK5 代理等），便于独立使用 frp 客户端完成某些工作。
* 高度扩展性的服务端插件系统，易于结合自身需求进行功能扩展。
* 服务端和客户端 UI 页面。

## 开发状态

frp 目前已被很多公司广泛用于测试、生产环境。

master 分支用于发布稳定版本，dev 分支用于开发，您可以尝试下载最新的 release 版本进行测试。

我们正在进行 v2 大版本的开发，将会尝试在各个方面进行重构和升级，且不会与 v1 版本进行兼容，预计会持续较长的一段时间。

现在的 v0 版本将会在合适的时间切换为 v1 版本并且保证兼容性，后续只做 bug 修复和优化，不再进行大的功能性更新。

### 关于 v2 的一些说明

v2 版本的复杂度和难度比我们预期的要高得多。我只能利用零散的时间进行开发，而且由于上下文经常被打断，效率极低。由于这种情况可能会持续一段时间，我们仍然会在当前版本上进行一些优化和迭代，直到我们有更多空闲时间 来推进大版本的重构，或者也有可能放弃一次性的重构，而是采用渐进的方式在当前版本上逐步做一些可能会导致不兼容的修改。

v2 的构想是基于我多年在云原生领域，特别是在 K8s 和 ServiceMesh 方面的工作经验和思考。它的核心是一个现代化的四层和七层代理，类似于 envoy。这个代理本身高度可扩展，不仅可以用于实现内网穿透的功能，还可以应用于 更多领域。在这个高度可扩展的内核基础上，我们将实现 frp v1 中的所有功能，并且能够以一种更加优雅的方式实现原先架构中无法实现或不易实现的功能。同时，我们将保持高效的开发和迭代能力。

除此之外，我希望 frp 本身也成为一个高度可扩展的系统和平台，就像我们可以基于 K8s 提供一系列扩展能力一样。在 K8s 上，我们可以根据企业需求进行定制化开发，例如使用 CRD、controller 模式、webhook、CSI 和 CNI 等。在 frp v1 中，我们引入了服务端插件的概念，实现了一些简单的扩展性。但是，它实际上依赖于简单的 HTTP 协议，并且需要用户自己启动独立的进程和管理。这种方式远远不够灵活和方便，而且现实世界的需求千差万别，我们不能期望一个由少数人维护的非营利性开源项目能够满足所有人的需求。

最后，我们意识到像配置管理、权限验证、证书管理和管理 API 等模块的当前设计并不够现代化。尽管我们可能在 v1 版本中进行一些优化，但确保兼容性是一个令人头疼的问题，需要投入大量精力来解决。

非常感谢您对 frp 的支持。
## 基于开源frp基础上增加IP4P
## 新特性：IP4P 与 TXT 记录支持 (STUN 穿透优化)

本版本增加了对 IP4P 和 DNS TXT 记录解析的支持，旨在解决 STUN 穿透端口频繁变动的问题。

### 1. 简介
STUN 穿透的端口通常不固定。通过 Cloudflare 页面规则虽然可以重定向域名，但很多应用程序不支持重定向或非标准端口。本功能借鉴了 NATMAP 的思路，利用 DNS 记录（AAAA 或 TXT）来动态更新和解析公网 IP 与端口。

### 2. 实现原理
*   **IP4P**: 将 IPv4 地址和端口编码进 IPv6 地址中（`2001::{port}:{ipv4-hi16}:{ipv4-lo16}`），通过 AAAA 记录发布。
*   **TXT 记录**: 将 `公网IP:端口` 进行 Base64 编码（如 `183.6.66.666:6666` -> `MTgzLjYuNjYuNjY2OjY2NjY=`），存入 TXT 记录中。

### 3. 配置步骤

#### 第一步：内网配置 (frps & frpc)
在内网同时部署 frps 和 frpc。
*   **frps【使用官方包就行，不限制】**: 正常配置，监听 7000 端口。
*   **frpc【必须使用ip4p版本】**: 配置 `stcp` 或 `sudp` 将本地服务（如 Jellyfin, Nginx）暴露。

**frps 正常部署即可，配置没有特殊要求**


**被访问的frpc.toml 示例**:
```toml
serverAddr = "192.168.0.1"
serverPort = 7000

auth.method = "token"
auth.token = "frps定义token"
# console or real logFile path like ./frpc.log
log.to = "/var/log/frpc.log"
# trace, debug, info, warn, error
log.level = "info"
log.maxDays = 3
# disable log colors when log.to is console, default is false
log.disablePrintColor = true

[[proxies]]
name = "frp_ssh_op"
type = "stcp"
# 只有与此处设置的 secretKey 一致的用户才能访问此服务
secretKey = "秘钥"
localIP = "192.168.0.1"
localPort = 22

[[proxies]]
name = "frp_alist"
type = "stcp"
secretKey = "秘钥"
localIP = "192.168.0.1"
localPort = 5244

```

#### 第二步：设置 DNS 记录自动更新
使用脚本监控 IP 变动并调用 Cloudflare API 更新 TXT 记录。

**脚本片段**:
```bash
# 替换为你的 Cloudflare 信息
domain="ip4p.yourdomain.com"
zone="your_zone_id"
txt_id="your_txt_record_id"
email="your_email"
key="your_api_key"

# 获取 IP 并 Base64 编码
addr_64=$(echo -n ${ipAddr} | base64)

# 更新 TXT 记录
curl -s -X PUT "https://api.cloudflare.com/client/v4/zones/$zone/dns_records/$txt_id" \
-H "X-Auth-Email: $email" -H "X-Auth-Key: $key" \
-H "Content-Type: application/json" \
--data '{"type":"TXT","name":"'$domain'","content":"'$addr_64'","ttl":60,"proxied":false}'
```

#### 第三步：异地配置 (Visitor frpc)
在访问端使用支持该功能的 frpc 版本。

**frpc.ini 示例**:
完整的配置示例请参见 `conf/frpc_ip4p_example.toml`。
```toml
serverAddr = "ip4p.yourdomain.com"

[ip4p]
mode = "lookup_text"


[auth]
method = "token"
token = "****"

[log]
to = "./frpc.log"
level = "info"
maxDays = 3
disablePrintColor = true

[[visitors]]
name = "visitor_frp_ssh_op"
type = "stcp"
# 要访问的 stcp 代理的名字
serverName = "frp_ssh_op"
secretKey = "秘钥"
# 绑定本地端口以访问 SSH 服务
bindAddr = "127.0.0.1"
bindPort = 20001

[[visitors]]
name = "visitor_frp_alist"
type = "stcp"
# 要访问的 stcp 代理的名字
serverName = "frp_alist"
secretKey = "秘钥"
# 绑定本地端口以访问 SSH 服务
bindAddr = "127.0.0.1"
bindPort = 15244
```

访问 `127.0.0.1:15244` 即可连接内网alist服务。

## 文档

完整文档已经迁移至 [https://gofrp.org](https://gofrp.org)。

## 为 frp 做贡献

frp 是一个免费且开源的项目，我们欢迎任何人为其开发和进步贡献力量。

* 在使用过程中出现任何问题，可以通过 [issues](https://github.com/fatedier/frp/issues) 来反馈。
* Bug 的修复可以直接提交 Pull Request 到 dev 分支。
* 如果是增加新的功能特性，请先创建一个 issue 并做简单描述以及大致的实现方法，提议被采纳后，就可以创建一个实现新特性的 Pull Request。
* 欢迎对说明文档做出改善，帮助更多的人使用 frp，特别是英文文档。
* 贡献代码请提交 PR 至 dev 分支，master 分支仅用于发布稳定可用版本。
* 如果你有任何其他方面的问题或合作，欢迎发送邮件至 fatedier@gmail.com 。

**提醒：和项目相关的问题请在 [issues](https://github.com/fatedier/frp/issues) 中反馈，这样方便其他有类似问题的人可以快速查找解决方法，并且也避免了我们重复回答一些问题。**

## 关联项目

* [gofrp/plugin](https://github.com/gofrp/plugin) - frp 插件仓库，收录了基于 frp 扩展机制实现的各种插件，满足各种场景下的定制化需求。
* [gofrp/tiny-frpc](https://github.com/gofrp/tiny-frpc) - 基于 ssh 协议实现的 frp 客户端的精简版本(最低约 3.5MB 左右)，支持常用的部分功能，适用于资源有限的设备。
