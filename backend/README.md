# EasyCpp Backend

EasyCpp 后端服务，基于 Golang + Gin，使用 PostgreSQL 做主数据持久化、Redis 做会话存储，JWT 做访问令牌鉴权

当前版本已提供 C++ 在线编译执行接口（同步执行）。

## 完整目录结构

```text
backend/
├─ cmd/
│  └─ api/
│     └─ main.go
├─ internal/
│  ├─ app/
│  │  ├─ container.go
│  │  └─ router.go
│  ├─ config/
│  │  └─ config.go
│  ├─ dto/
│  │  ├─ auth_dto.go
│  │  ├─ cpp_ide_dto.go
│  │  ├─ comment_dto.go
│  │  ├─ post_dto.go
│  │  └─ user_dto.go
│  ├─ handler/
│  │  ├─ auth_handler.go
│  │  ├─ comment_handler.go
│  │  ├─ cpp_ide_handler.go
│  │  ├─ post_handler.go
│  │  └─ user_handler.go
│  ├─ middleware/
│  │  ├─ auth.go
│  │  └─ cors.go
│  ├─ model/
│  │  ├─ comment.go
│  │  ├─ post.go
│  │  └─ user.go
│  ├─ platform/
│  │  ├─ postgres/
│  │  │  └─ postgres.go
│  │  └─ rediscache/
│  │     └─ redis.go
│  ├─ repository/
│  │  ├─ comment_repository.go
│  │  ├─ errors.go
│  │  ├─ pg_errors.go
│  │  ├─ post_repository.go
│  │  ├─ schema.go
│  │  ├─ session_repository.go
│  │  └─ user_repository.go
│  ├─ security/
│  │  └─ jwt.go
│  └─ service/
│     ├─ auth_service.go
│     ├─ comment_service.go
│     ├─ cpp_ide_service.go
│     ├─ post_service.go
│     └─ user_service.go
├─ pkg/
│  └─ response/
│     └─ response.go
├─ go.mod
├─ go.sum
└─ README.md
```

## 模块说明

### cmd

- `cmd/api/main.go`
- 应用启动入口
- 负责初始化应用容器、启动 HTTP 服务、处理进程退出时的资源释放

### internal/app

- `router.go`: 定义路由分组、公开接口和鉴权接口
- `container.go`: 依赖注入容器，装配 Config、PostgreSQL、Redis、Repository、Service、Handler
- 这是应用层“编排中枢”，不放业务细节

### internal/config

- 统一读取环境变量并提供默认值
- 管理服务监听地址、数据库连接、Redis 连接、JWT 参数、密码哈希强度等

### internal/dto

- 请求/响应的数据传输对象（Data Transfer Object）
- 负责 HTTP 入参校验标签（`binding`）
- 避免将数据库模型直接暴露到 API 层
- 新增 `cpp_ide_dto.go`：定义 C++ 编译执行请求与响应结构

### internal/handler

- 控制器层（HTTP Handler）
- 负责参数解析、调用 Service、返回统一响应
- 不直接写数据库访问逻辑
- 新增 `cpp_ide_handler.go`：提供 `POST /api/v1/ide/cpp/run` 接口

### internal/middleware

- `cors.go`: 跨域处理
- `auth.go`: Bearer Token 解析、JWT 校验、Redis 会话校验、用户上下文注入

### internal/model

- 核心领域模型（用户、贴文、评论）
- 在 Service 和 Repository 之间传递的核心对象

### internal/platform

- 基础设施客户端封装
- `postgres`: PostgreSQL 连接池创建与健康检查
- `rediscache`: Redis 客户端创建与健康检查

### internal/repository

- 数据访问层（等价于很多项目里的 DAO 职责）
- `*_repository.go`: 具体 SQL/Redis 读写实现
- `schema.go`: 启动时自动建表
- `errors.go` / `pg_errors.go`: 统一数据层错误与数据库错误映射

### internal/security

- JWT 的生成、解析、签名校验相关逻辑
- 封装 `Claims`，统一令牌字段规范

### internal/service

- 业务逻辑层
- 组合多个 Repository 实现业务规则
- 例如注册登录流程、会话管理、作者权限校验、评论写入前贴文存在性校验
- 新增 `cpp_ide_service.go`：执行 C++ 代码的写盘、编译、运行与超时控制

### pkg/response

- 统一 API 响应结构封装（成功/失败）
- 保持前后端对响应格式的稳定约定

## API 概览

### 公开接口

- `GET /api/v1/health`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/users`
- `GET /api/v1/users/:id`
- `GET /api/v1/posts`
- `GET /api/v1/posts/:id`
- `GET /api/v1/posts/:postId/comments`
- `POST /api/v1/ide/cpp/run`

### 需要鉴权（Bearer Token）

- `POST /api/v1/auth/logout`
- `GET /api/v1/auth/sessions`
- `DELETE /api/v1/auth/sessions/:sessionId`
- `POST /api/v1/posts`
- `DELETE /api/v1/posts/:id`
- `POST /api/v1/posts/:postId/comments`
- `DELETE /api/v1/comments/:commentId`

## 在线 IDE（C++）完整说明

### 功能边界

- 当前仅支持 C++（`g++`）编译与运行
- 接口为同步执行（单次请求返回完整结果）
- 为公开接口，不要求登录态
- 支持传入标准输入（`stdin`）
- 支持超时控制，超时后终止执行

### 接口定义

- 方法：`POST /api/v1/ide/cpp/run`
- Content-Type：`application/json`

请求体：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `code` | string | 是 | C++ 源码，长度 1~50000 |
| `stdin` | string | 否 | 标准输入，最大 10000 |
| `timeLimitMs` | int | 否 | 运行超时毫秒，`<=0` 使用默认值 |

返回体（`data`）：

| 字段 | 类型 | 说明 |
|---|---|---|
| `language` | string | 固定为 `cpp` |
| `compile` | object | 编译阶段结果 |
| `run` | object/null | 运行阶段结果；编译失败时为 `null` |
| `durationMs` | int64 | 整体耗时（毫秒） |

`compile` 字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| `succeeded` | bool | 编译是否成功 |
| `exitCode` | int | 编译进程退出码 |
| `stderr` | string | 编译错误输出 |

`run` 字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| `succeeded` | bool | 运行是否成功（退出码 0 且未超时） |
| `exitCode` | int | 运行进程退出码 |
| `stdout` | string | 标准输出 |
| `stderr` | string | 标准错误输出 |
| `timedOut` | bool | 是否触发超时 |

### 执行流程

1. 在临时目录写入 `main.cpp`
2. 调用 `g++ -std=c++17 -O2 -pipe` 编译
3. 编译成功后运行可执行文件（Windows 下为 `main.exe`）
4. 按 `timeLimitMs` 创建超时上下文
5. 采集 `stdout/stderr/exitCode` 并返回
6. 请求结束后清理临时目录

### 默认限制

| 限制项 | 值 |
|---|---|
| 默认超时 | `2000ms` |
| 最大超时 | `10000ms` |
| 输出截断上限 | `64KB` |

说明：超出输出上限后会在尾部附加 `[output truncated]`。

### 典型响应示例

编译与运行成功：

```json
{
	"code": 200,
	"message": "ok",
	"data": {
		"language": "cpp",
		"compile": {
			"succeeded": true,
			"exitCode": 0,
			"stderr": ""
		},
		"run": {
			"succeeded": true,
			"exitCode": 0,
			"stdout": "hello\n",
			"stderr": "",
			"timedOut": false
		},
		"durationMs": 38
	}
}
```

编译失败：

```json
{
	"code": 200,
	"message": "ok",
	"data": {
		"language": "cpp",
		"compile": {
			"succeeded": false,
			"exitCode": 1,
			"stderr": "main.cpp: In function ..."
		},
		"run": null,
		"durationMs": 12
	}
}
```

运行超时：

```json
{
	"code": 200,
	"message": "ok",
	"data": {
		"language": "cpp",
		"compile": {
			"succeeded": true,
			"exitCode": 0,
			"stderr": ""
		},
		"run": {
			"succeeded": false,
			"exitCode": 1,
			"stdout": "",
			"stderr": "Execution timed out.",
			"timedOut": true
		},
		"durationMs": 2005
	}
}
```

### 运行前置条件

- 机器需安装 `g++`，且可在 PATH 中执行
- 若编译器不存在，`compile.stderr` 会返回提示：`compiler not found: g++ is not installed or not in PATH`

### 安全说明（当前版本）

- 当前实现为本机进程执行模型，用于开发与原型验证
- 生产环境建议接入容器沙箱（隔离、禁网、资源限制、系统调用限制）

## 环境配置

推荐在 `backend` 目录使用 `.env`（或系统环境变量）管理配置。

### 必要配置项

| 变量名 | 默认值 | 说明 |
|---|---|---|
| `HTTP_ADDR` | `:8080` | HTTP 服务监听地址 |
| `POSTGRES_DSN` | `postgres://postgres:postgres@127.0.0.1:5432/easycpp?sslmode=disable` | PostgreSQL 连接串 |
| `REDIS_ADDR` | `127.0.0.1:6379` | Redis 地址 |
| `REDIS_PASSWORD` | 空 | Redis 密码 |
| `REDIS_DB` | `0` | Redis DB 编号 |
| `JWT_SECRET` | `change-me-in-production` | JWT 签名密钥（生产必须替换） |
| `JWT_ACCESS_TTL` | `2h` | Access Token 有效期 |
| `SESSION_TTL` | `24h` | Redis 会话有效期 |
| `BCRYPT_COST` | `12` | 密码哈希强度 |

### Windows PowerShell 示例

```powershell
$env:HTTP_ADDR=":8080"
$env:POSTGRES_DSN="postgres://postgres:postgres@127.0.0.1:5432/easycpp?sslmode=disable"
$env:REDIS_ADDR="127.0.0.1:6379"
$env:JWT_SECRET="replace-with-strong-secret"
```

## 开发与调试教程

### 1. 准备依赖

1. 安装 Go 1.22+
2. 启动 PostgreSQL（创建数据库 `easycpp`）
3. 启动 Redis（默认 `6379`）
4. 安装 C++ 编译器 `g++` 并确保在 PATH 中可执行

### 2. 安装依赖并启动后端

```bash
cd backend
go mod tidy
go run ./cmd/api
```

启动后默认监听 `http://127.0.0.1:8080`。

### 3. Build 教程（Windows / Linux）

下面是常用构建方式，建议先在 `backend` 目录执行。

#### Windows（PowerShell）

1. 本机构建（Windows 可执行文件）

```powershell
go build -o .\bin\easycpp-api.exe .\cmd\api
```

2. 构建并附带版本信息（可选）

```powershell
$version = "v0.1.0"
$commit = "dev"
go build -ldflags "-X main.version=$version -X main.commit=$commit" -o .\bin\easycpp-api.exe .\cmd\api
```

3. 在 Windows 上交叉编译 Linux 版本（可选）

```powershell
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o .\bin\easycpp-api-linux-amd64 .\cmd\api
Remove-Item Env:GOOS
Remove-Item Env:GOARCH
```

#### Linux（Bash）

1. 本机构建（Linux 可执行文件）

```bash
go build -o ./bin/easycpp-api ./cmd/api
```

2. 构建并附带版本信息（可选）

```bash
VERSION="v0.1.0"
COMMIT="dev"
go build -ldflags "-X main.version=$VERSION -X main.commit=$COMMIT" -o ./bin/easycpp-api ./cmd/api
```

3. 在 Linux 上交叉编译 Windows 版本（可选）

```bash
GOOS=windows GOARCH=amd64 go build -o ./bin/easycpp-api.exe ./cmd/api
```

#### 通用构建检查

```bash
go fmt ./...
go vet ./...
go build ./...
go test ./...
```

### 4. 快速验证接口

1. 健康检查

```bash
curl http://127.0.0.1:8080/api/v1/health
```

2. 注册

```bash
curl -X POST http://127.0.0.1:8080/api/v1/auth/register \
	-H "Content-Type: application/json" \
	-d '{"username":"tom","email":"tom@example.com","password":"12345678","bio":"hello"}'
```

3. 登录（记录返回的 `accessToken`）

```bash
curl -X POST http://127.0.0.1:8080/api/v1/auth/login \
	-H "Content-Type: application/json" \
	-d '{"email":"tom@example.com","password":"12345678"}'
```

4. 使用 Token 发帖

```bash
curl -X POST http://127.0.0.1:8080/api/v1/posts \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer <accessToken>" \
	-d '{"title":"First Post","summary":"intro","content":"This is my first post content.","language":"cpp","difficulty":"beginner","tags":["c++","basic"]}'
```

5. C++ 在线执行（公开接口）

```bash
curl -X POST http://127.0.0.1:8080/api/v1/ide/cpp/run \
	-H "Content-Type: application/json" \
	-d '{"code":"#include <iostream>\nint main(){std::cout<<\"hello\";return 0;}","stdin":"","timeLimitMs":2000}'
```

### 5. 常见调试命令

```bash
go fmt ./...
go vet ./...
go build ./...
go test ./...
```

### 6. 常见问题排查

1. 启动失败且提示数据库连接错误
- 检查 `POSTGRES_DSN` 用户名、密码、库名、端口

2. 登录成功但访问鉴权接口仍然 401
- 检查 `Authorization` 是否为 `Bearer <token>` 格式
- 检查 Redis 是否正常，确保会话未过期

3. 第一次启动缺表
- 应用已内置自动建表逻辑，若权限不足请确认数据库账户具有建表权限

## 开发建议

1. 新增业务优先遵循 `handler -> service -> repository` 分层
2. 统一通过 `pkg/response` 返回接口结果
3. 在 `dto` 中维护 API 入参校验规则
4. 对关键业务流程补充 `service` 层单测
