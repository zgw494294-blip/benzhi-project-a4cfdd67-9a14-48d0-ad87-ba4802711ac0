# BENZHI_README

## 项目说明
- 项目：benzhi-project-a4cfdd67-9a14-48d0-ad87-ba4802711ac0
- 项目用途：科研开放数据集授权审核与封存服务已完成。服务采用 model、workflow、store、audit、httpapi 分层，实现批次、来源、审核、补证、冻结、关闭及审计查询全流程，使用本地 JSON 账本和连续哈希链保障可追溯性，支持幂等请求、乐观并发控制、回环地址配置与有界自检。
- Go 工具链：`golang:1.22`
- 前端工具链：无

## 标准构建、运行和测试命令
进入容器后执行：
cd '/app' && GOTOOLCHAIN=local go build ./...
cd /app && GOTOOLCHAIN=local go run . -addr=127.0.0.1:19081 -selfcheck
cd '/app' && GOTOOLCHAIN=local go test ./...

## Docker 构建和进入容器
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh benzhi-project-a4cfdd67-9a14-48d0-ad87-ba4802711ac0-amd64 linux/amd64
./build_benzhi_docker.sh benzhi-project-a4cfdd67-9a14-48d0-ad87-ba4802711ac0-arm64 linux/arm64
docker run -it benzhi-project-a4cfdd67-9a14-48d0-ad87-ba4802711ac0-amd64:latest

## 题目验证命令
1. 预期退出码 0：`go test ./...`
2. 预期退出码 0：`go run . -addr=127.0.0.1:19081 -selfcheck`
