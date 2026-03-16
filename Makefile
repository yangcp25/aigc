# ==============================================================================
# 定义核心变量
# ==============================================================================
APP_NAME := aigc
BIN_DIR := bin
# 核心修正：编译和运行必须指向整个目录，而不是单个 main.go 文件
CMD_DIR := ./cmd/server

# 获取当前的 Git Commit Hash 作为版本号 (如果不在 git 仓库则回退为 dev)
VERSION := $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date "+%F %T")

# Go 编译注入参数 (方便在程序运行期间获取版本和构建时间)
LDFLAGS := -X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)'

# 动态获取 GOPATH 的第一个路径（完美防范 GVM 的多路径/找不到环境变量问题）
GOPATH_FIRST := $(shell go env GOPATH | awk -F':' '{print $$1}')

# ==============================================================================
# 定义 Make 指令
# ==============================================================================
.PHONY: all build run wire clean fmt test help

# 默认执行的指令
all: fmt wire build

## 生成 wire 依赖注入代码
wire:
	@echo "=> 正在使用 go generate 生成 Wire 代码..."
	# 核心黑科技 1：临时把 GOPATH/bin 塞进 PATH 里，确保能找到 wire 工具
	# 核心黑科技 2：必须带上 -tags wireinject，穿透隐身衣去读图纸
	@export PATH="$(GOPATH_FIRST)/bin:$$PATH" && go generate -tags wireinject ./...
	@echo "=> Wire 生成完成 ✅"

## 编译二进制文件
build: wire
	@echo "=> 正在编译 $(APP_NAME) ($(VERSION))..."
	@mkdir -p $(BIN_DIR)
	# 核心修正 3：编译 $(CMD_DIR) 整个包，把 main.go 和 wire_gen.go 一起打包
	@go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)
	@echo "=> 编译成功，输出至 $(BIN_DIR)/$(APP_NAME) ✅"

## 运行开发环境
run: wire
	@echo "=> 正在启动服务..."
	# 同样地，运行整个目录包
	@go run $(CMD_DIR)

## 格式化代码并整理依赖
fmt:
	@echo "=> 格式化代码..."
	@go fmt ./...
	@go mod tidy

## 运行单元测试
test:
	@echo "=> 运行测试..."
	@go test -v -count=1 ./...

## 清理构建产物
clean:
	@echo "=> 清理临时文件..."
	@rm -rf $(BIN_DIR)
	@rm -f $(CMD_DIR)/wire_gen.go

## 打印帮助信息
help:
	@echo "可用的 Make 命令:"
	@echo "  make all   - 格式化、生成注入代码并编译二进制产物 (默认)"
	@echo "  make wire  - 仅重新生成依赖注入代码"
	@echo "  make run   - 生成注入代码并直接启动服务 (日常开发用这个！)"
	@echo "  make build - 生成注入代码并编译输出到 bin 目录"
	@echo "  make fmt   - 格式化代码并清理 go.mod"
	@echo "  make clean - 清除编译产物和生成的注入文件"