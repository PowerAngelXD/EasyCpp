# EasyCpp

EasyCpp 是一个面向编程学习场景的全栈项目，包含：

- `backend`：Golang + Gin 的 REST API，提供用户、认证、帖子、评论与会话能力。
- `frontend`：Vue 3 + Vite 的学习社区原型界面，包含 Feed、帖子详情、Profile 与 Playground 页面。

当前仓库适合用于：

- 学习前后端分层与工程组织方式
- 作为小型社区产品/在线练习平台的起步模板
- 后续扩展代码执行沙箱、多语言编译、在线评测等能力

## 目录结构

```text
EasyCpp/
├─ backend/                  # Go 后端服务
│  ├─ cmd/api/main.go        # 启动入口
│  ├─ internal/              # 分层业务代码
│  ├─ pkg/response/          # 统一响应封装
│  └─ README.md
├─ frontend/                 # Vue 前端工程
│  ├─ src/components/
│  ├─ src/views/
│  ├─ src/router/
│  └─ README.md
├─ go.work
└─ README.MD
```

## 技术栈

### Backend

- Go `1.24`
- Gin（HTTP 路由）
- PostgreSQL（主数据）
- Redis（会话存储）
- JWT + Bearer Token（认证）

### Frontend

- Vue `3`
- Vue Router `4`
- Vite `8`
- C++ Playground（前端调用后端编译执行 API）

## 核心能力

### 已实现（后端）

- 用户注册与登录
- JWT 访问令牌签发与校验
- Redis 会话管理（登出、会话列表、撤销指定会话）
- 帖子创建、列表、详情、删除（作者权限）
- 评论创建、列表、删除（作者权限）
- 启动时自动建表（`users` / `posts` / `comments`）
- C++ 在线编译执行接口（同步执行）

### 已实现（前端）

- 应用壳层与侧边导航
- 首页 Feed 列表（当前使用本地 mock 数据）
- 帖子详情与评论线程 UI（原型）
- Profile 页面与贡献热力图（原型）
- Playground 页面（仅支持 C++，调用后端编译执行）
- 主题切换（本地存储 + 系统主题回退）

### 当前状态说明

- 前端 Playground 已对接后端 C++ 执行 API，其余业务页面仍以 UI/原型为主。
- 后端接口可独立运行并通过 curl/Postman 调试。

## 快速开始

### 1. 环境要求

- Go `>= 1.24`
- Node.js `>= 20`
- PostgreSQL `>= 14`（建议）
- Redis `>= 6`（建议）

### 2. 启动后端

```bash
cd backend
go mod tidy
go run ./cmd/api
```

默认监听：`http://127.0.0.1:8080`

### 3. 启动前端

```bash
cd frontend
npm install
npm run dev
```

默认访问：`http://127.0.0.1:5173`

如果前端与后端地址不同，可在 `frontend/.env` 配置：

```bash
VITE_API_BASE_URL=http://127.0.0.1:8080
```

### 4. 常用检查命令

### Backend

```bash
cd backend
go fmt ./...
go vet ./...
go test ./...
go build ./...
```

### Frontend

```bash
cd frontend
npm run build
npm run preview
npm run format:check
```

## 配置说明（后端）

后端通过环境变量读取配置：

| 变量名 | 默认值 | 说明 |
|---|---|---|
| `APP_ENV` | `development` | 运行环境 |
| `HTTP_ADDR` | `:8080` | HTTP 监听地址 |
| `POSTGRES_DSN` | `postgres://postgres:postgres@127.0.0.1:5432/easycpp?sslmode=disable` | PostgreSQL 连接串 |
| `REDIS_ADDR` | `127.0.0.1:6379` | Redis 地址 |
| `REDIS_PASSWORD` | 空 | Redis 密码 |
| `REDIS_DB` | `0` | Redis DB 编号 |
| `JWT_SECRET` | `change-me-in-production` | JWT 签名密钥 |
| `JWT_ACCESS_TTL` | `2h` | Access Token 过期时长 |
| `SESSION_TTL` | `24h` | 会话过期时长 |
| `BCRYPT_COST` | `12` | 密码哈希强度（最低 10） |

PowerShell 示例：

```powershell
$env:HTTP_ADDR=":8080"
$env:POSTGRES_DSN="postgres://postgres:postgres@127.0.0.1:5432/easycpp?sslmode=disable"
$env:REDIS_ADDR="127.0.0.1:6379"
$env:JWT_SECRET="replace-with-strong-secret"
```

## API 概览

基础前缀：`/api/v1`

### 公开接口

- `GET /health`
- `POST /auth/register`
- `POST /auth/login`
- `GET /users`
- `GET /users/:id`
- `GET /posts`
- `GET /posts/:id`
- `GET /posts/:postId/comments`
- `POST /ide/cpp/run`

### 需要鉴权（Bearer Token）

- `POST /auth/logout`
- `GET /auth/sessions`
- `DELETE /auth/sessions/:sessionId`
- `POST /posts`
- `DELETE /posts/:id`
- `POST /posts/:postId/comments`
- `DELETE /comments/:commentId`

## 开发路线图（Roadmap）

- Feed/Post/Profile 对接后端真实 API（替换 mock）
- 增加统一鉴权状态管理（前端）
- 增加后端单元测试与集成测试
- 引入容器化部署（Docker / Compose）
- C++ 执行链路增加容器沙箱与更严格资源隔离

## 贡献指南

欢迎通过 Issue / PR 参与贡献。建议遵循：

1. 保持 `handler -> service -> repository` 分层边界清晰。
2. 新接口使用 `dto` + `pkg/response` 保持响应与校验一致。
3. 提交前至少执行：
	- 后端：`go test ./...`、`go vet ./...`
	- 前端：`npm run build`、`npm run format:check`

## License

本项目采用 **GNU General Public License v3.0（GPL-3.0）**。

- 许可证全文见 [LICENSE](LICENSE)
- 你可以自由使用、修改和再分发本项目，但再分发的衍生作品必须继续以 GPL-3.0 开源
- 贡献到本仓库的代码默认视为同意以 GPL-3.0 进行许可
