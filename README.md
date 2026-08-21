# 科研开放数据集授权审核与封存服务

本项目为科研机构提供可追溯的数据集开放前审核流程。管理员创建审核批次并登记来源凭证，伦理审查员逐项作出决定，缺证项目可以退回补正；全部来源通过后生成不可编辑的发布快照，关闭批次时签发带连续哈希链的审计证书。服务只使用本地 JSON 账本，不依赖交易系统或外部平台。

## 构建、运行与测试

标准回归测试：

```text
go test ./...
```

运行有界自检（自动创建一条完整业务链并退出）：

```text
go run . -addr=127.0.0.1:19081 -selfcheck
```

启动 HTTP 服务：

```text
go run . -addr=127.0.0.1:19081
```

也可以设置 `PORT=19082`，服务会绑定 `127.0.0.1:19082`。账本默认写入 `ledger.json`，可通过 `-ledger` 指定路径。

## API

服务提供 `POST /v1/batches`、`POST /v1/batches/{batch_id}/sources`、`POST /v1/batches/{batch_id}/sources/batch`、`POST /v1/batches/{batch_id}/reviews/{source_id}`、`POST /v1/batches/{batch_id}/resubmit`、`POST /v1/batches/{batch_id}/freeze`、`POST /v1/batches/{batch_id}/close`、`GET /v1/batches/{batch_id}` 与 `GET /v1/batches/{batch_id}/audit`。

批量来源请求使用 `{"sources":[...]}`，会执行批次内重复校验并原子写入；`GET /v1/batches/{batch_id}/policy`（同样支持 `/precheck`）返回发布范围预检报告，`GET /v1/batches/{batch_id}/status` 返回审核进度和审计状态，`GET /v1/batches/{batch_id}/export` 导出冻结快照，`GET /v1/batches/{batch_id}/certificate/verify` 在线验证封存证书。写请求支持 `Idempotency-Key`，并可用 `If-Match` 或 JSON 中的 `expected_version` 做乐观并发校验。
