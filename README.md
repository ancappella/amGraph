# Kratos Project Template

框架使用 [Kratos](https://github.com/go-kratos/kratos)。

## 服务启动

1. 准备 PostgreSQL（与 `configs/config.yaml` 中 `data.database.source` 一致），或按需改配置。
2. 生成代码并编译：

```bash
./scripts/codegen.sh
go build -o bin/amGraph ./cmd/amGraph
./bin/amGraph -conf ./configs
```

`./scripts/codegen.sh` 等价于 `make api && make config && go generate ./...`，但使用 **Homebrew 的 `gmake`** 并自动把 `protoc-gen-*` 装到项目内 `.tools/bin`，避免 macOS 自带 **`/usr/bin/make` 因 Command Line Tools / `xcrun` 损坏而失败**。

若仍想用 `make` 手写命令，在已 `brew install make` 的前提下：

```bash
export PATH="$(pwd)/.tools/bin:$(go env GOPATH)/bin:$PATH"   # 需已安装 protoc 插件，见 scripts/codegen.sh
gmake api && gmake config && go generate ./...
```

**彻底修复本机工具链**（可选，装好后可继续用 Apple `make`）：

```bash
sudo rm -rf /Library/Developer/CommandLineTools
xcode-select --install
```

无 `gmake` 时可用纯 bash（需本机已有 `protoc` 与各插件在 `PATH`）：`./scripts/gen-proto.sh`
