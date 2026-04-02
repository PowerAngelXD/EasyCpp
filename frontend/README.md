# EasyCpp Frontend
- EasyCpp 前端基于 Vue 3 + Vite，当前重点是学习社区界面原型与 C++ 在线 IDE 页面。
## 技术栈

- Vue 3
- Vue Router 4
- Vite 8

## 页面能力

- Home（Feed 原型，当前使用 mock 数据）
- Post（帖子详情 + 评论线程 UI 原型）
- Profile（个人页 + 热力图原型）
- Playground（仅支持 C++，调用后端编译执行 API）

## 项目结构

```text
frontend/
├─ src/
│  ├─ components/
│  │  ├─ AppShell.vue
│  │  ├─ CodeEditor.vue
│  │  ├─ SidebarNav.vue
│  │  ├─ TerminalPanel.vue
│  │  └─ ThemeToggle.vue
│  ├─ router/
│  │  └─ index.js
│  ├─ utils/
│  │  ├─ cppRunner.js
│  │  └─ theme.js
│  └─ views/
│     ├─ HomeView.vue
│     ├─ PlaygroundView.vue
│     ├─ PostView.vue
│     └─ ProfileView.vue
├─ index.html
├─ package.json
└─ vite.config.js
```

## 开发启动

```bash
npm install
npm run dev
```

默认地址：`http://127.0.0.1:5173`

## 构建与检查

```bash
npm run build
npm run preview
npm run format
npm run format:check
```

## C++ IDE 对接说明

Playground 页面会调用后端接口：

- `POST /api/v1/ide/cpp/run`

默认后端地址为 `http://127.0.0.1:8080`。你也可以通过环境变量覆盖。

在前端目录创建 `.env`：

```bash
VITE_API_BASE_URL=http://127.0.0.1:8080
```

注意：后端需要安装 `g++`，否则编译会失败。

# Vue 3 + Vite

This template should help get you started developing with Vue 3 in Vite. The template uses Vue 3 `<script setup>` SFCs, check out the [script setup docs](https://v3.vuejs.org/api/sfc-script-setup.html#sfc-script-setup) to learn more.

Learn more about IDE Support for Vue in the [Vue Docs Scaling up Guide](https://vuejs.org/guide/scaling-up/tooling.html#ide-support).
