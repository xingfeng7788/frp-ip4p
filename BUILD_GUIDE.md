# FRP 构建和打包指南

## 快速开始

### 1. 本地构建（当前平台）

构建 frpc 和 frps 到 `bin/` 目录：

```bash
make
```

或者分别构建：

```bash
# 只构建客户端
make frpc

# 只构建服务端
make frps
```

构建完成后，二进制文件位于：
- `bin/frpc` - 客户端
- `bin/frps` - 服务端

### 2. 交叉编译（多平台）

编译所有支持的平台：

```bash
make -f Makefile.cross-compiles
```

编译结果在 `release/` 目录下，包含以下平台：
- darwin (macOS): amd64, arm64
- linux: amd64, arm, arm64, mips64, mips64le, mips, mipsle, riscv64, loong64
- windows: amd64, arm64
- freebsd: amd64
- openbsd: amd64
- android: arm64

### 3. 完整打包（推荐）

生成所有平台的发布包（.tar.gz 和 .zip）：

```bash
./package.sh
```

打包结果在 `release/packages/` 目录下，每个包包含：
- frpc / frpc.exe
- frps / frps.exe
- frpc.toml (配置文件)
- frps.toml (配置文件)
- LICENSE

包命名格式：`frp_<version>_<os>_<arch>.tar.gz` 或 `.zip`

## 详细命令说明

### 基础命令

```bash
# 格式化代码
make fmt

# 代码检查
make vet

# 运行测试
make test

# 运行 E2E 测试
make e2e

# 清理构建产物
make clean
```

### 自定义构建

#### 指定 LDFLAGS

```bash
# 添加版本信息等
LDFLAGS="-s -w -X main.version=1.0.0" make frpc
```

#### 单平台交叉编译

```bash
# Linux AMD64
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -tags frpc -o frpc_linux_amd64 ./cmd/frpc

# Windows AMD64
env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w" -tags frpc -o frpc_windows_amd64.exe ./cmd/frpc

# macOS ARM64 (Apple Silicon)
env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags "-s -w" -tags frpc -o frpc_darwin_arm64 ./cmd/frpc

# Linux ARM64
env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags "-s -w" -tags frpc -o frpc_linux_arm64 ./cmd/frpc
```

## 常用打包场景

### 场景 1: 快速测试本地构建

```bash
# 构建当前平台
make

# 测试
./bin/frpc --version
./bin/frps --version
```

### 场景 2: 构建特定平台

```bash
# 只构建 Linux AMD64
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -tags frpc -o bin/frpc_linux_amd64 ./cmd/frpc
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -tags frps -o bin/frps_linux_amd64 ./cmd/frps
```

### 场景 3: 发布版本（完整打包）

```bash
# 1. 确保代码已提交
git status

# 2. 运行测试
make test

# 3. 完整打包
./package.sh

# 4. 检查打包结果
ls -lh release/packages/
```

### 场景 4: 只打包常用平台

如果只需要打包常用平台，可以修改 `package.sh` 中的变量：

```bash
# 编辑 package.sh，修改这些行：
os_all='linux windows darwin'
arch_all='amd64 arm64'
```

然后运行：

```bash
./package.sh
```

## 构建选项说明

### CGO_ENABLED=0
禁用 CGO，生成纯静态二进制文件，无外部依赖

### -trimpath
移除文件系统路径，使构建可重现

### -ldflags "-s -w"
- `-s`: 去除符号表
- `-w`: 去除 DWARF 调试信息
- 减小二进制文件大小

### -tags frpc / -tags frps
构建标签，用于条件编译

## 验证构建

### 检查版本

```bash
./bin/frpc --version
./bin/frps --version
```

### 检查二进制文件

```bash
# 查看文件大小
ls -lh bin/

# 查看文件类型
file bin/frpc
file bin/frps

# 检查依赖（Linux）
ldd bin/frpc  # 应该显示 "not a dynamic executable"
```

## 常见问题

### Q: 构建失败，提示找不到 go 命令
A: 确保已安装 Go 1.24.0 或更高版本：
```bash
go version
```

### Q: 交叉编译失败
A: 确保设置了正确的环境变量：
```bash
export GO111MODULE=on
export CGO_ENABLED=0
```

### Q: 打包脚本权限不足
A: 添加执行权限：
```bash
chmod +x package.sh
```

### Q: 如何减小二进制文件大小
A: 使用 UPX 压缩（可选）：
```bash
# 安装 UPX
# macOS: brew install upx
# Linux: apt-get install upx-ucl

# 压缩二进制文件
upx --best --lzma bin/frpc
upx --best --lzma bin/frps
```

## 推荐工作流

### 开发阶段

```bash
# 1. 修改代码
# 2. 格式化
make fmt

# 3. 本地构建测试
make

# 4. 运行测试
make test

# 5. 测试功能
./bin/frpc -c conf/frpc.toml
```

### 发布阶段

```bash
# 1. 确保所有测试通过
make alltest

# 2. 完整打包
./package.sh

# 3. 验证打包结果
ls -lh release/packages/

# 4. 测试主要平台的二进制文件
# （可以在 Docker 容器中测试不同平台）
```

## Docker 构建（可选）

如果需要在 Docker 中构建：

```bash
# 创建构建镜像
docker run --rm -v "$PWD":/workspace -w /workspace golang:1.24 make

# 交叉编译
docker run --rm -v "$PWD":/workspace -w /workspace golang:1.24 make -f Makefile.cross-compiles

# 完整打包
docker run --rm -v "$PWD":/workspace -w /workspace golang:1.24 ./package.sh
```

## 快速命令参考

```bash
# 本地构建
make                                    # 构建当前平台
make frpc                               # 只构建客户端
make frps                               # 只构建服务端

# 交叉编译
make -f Makefile.cross-compiles         # 所有平台

# 打包
./package.sh                            # 完整打包

# 测试
make test                               # 单元测试
make e2e                                # E2E 测试
make alltest                            # 所有测试

# 清理
make clean                              # 清理构建产物
rm -rf release/                         # 清理发布文件
```

## 输出目录结构

```
frp/
├── bin/                    # 本地构建输出
│   ├── frpc
│   └── frps
├── release/                # 交叉编译输出
│   ├── frpc_linux_amd64
│   ├── frps_linux_amd64
│   ├── frpc_windows_amd64.exe
│   ├── frps_windows_amd64.exe
│   └── packages/           # 打包输出
│       ├── frp_x.x.x_linux_amd64.tar.gz
│       ├── frp_x.x.x_windows_amd64.zip
│       └── ...
```
