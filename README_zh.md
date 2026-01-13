<div align="center">

![OpenAuth Logo](./web/public/logo.svg)

# OpenAuth

**一个开源的身份认证与访问管理（IAM）平台，提供企业级的单点登录（SSO）、多因素认证（MFA）、用户管理和应用集成能力。**

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://golang.org/)
[![React](https://img.shields.io/badge/React-18+-61DAFB?logo=react)](https://reactjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.3+-3178C6?logo=typescript)](https://www.typescriptlang.org/)

[English](./README.md) | [中文](./README_zh.md)

</div>

## ✨ 特性

### 核心功能
- ✅ **用户认证**：登录、注册、密码管理
- ✅ **多因素认证（MFA）**：支持 TOTP、SMS、Email
- ✅ **用户管理**：完整的用户 CRUD 操作
- ✅ **应用管理**：应用创建、配置和管理
- ✅ **角色权限**：基于角色的访问控制（RBAC）
- ✅ **安全策略**：密码策略、MFA 策略、白名单管理
- ✅ **审计日志**：完整的操作审计记录
- ✅ **会话管理**：会话查看、删除、统计
- ✅ **组织架构**：组织管理和用户组管理
- ✅ **通知服务**：邮件和短信通知
- ✅ **API 文档**：Swagger 文档支持

### SSO 协议支持
- ✅ **OAuth 2.0 / OIDC**：标准 OAuth 2.0 和 OpenID Connect 协议
- ✅ **SAML 2.0**：企业级 SAML 单点登录（SSO、Metadata）
- ✅ **LDAP**：LDAP 目录服务集成

### 技术特点
- 🎨 **现代化 UI**：基于 React + Ant Design，极简黑白主题
- 🚀 **高性能**：Go 后端，PostgreSQL + Redis
- 🔒 **安全可靠**：JWT 认证、bcrypt 密码加密、MFA 支持
- 📦 **易于部署**：Docker 支持，配置简单
- 🌐 **国际化**：支持多语言（i18n）

## 🏗️ 技术栈

### 后端
- **语言**：Go 1.21+
- **框架**：Gin
- **数据库**：PostgreSQL 14+
- **缓存**：Redis 7+
- **ORM**：GORM
- **认证**：JWT

### 前端
- **框架**：React 18 + TypeScript
- **构建工具**：Vite
- **UI 组件库**：Ant Design 5
- **状态管理**：Zustand
- **HTTP 客户端**：Axios
- **数据获取**：TanStack Query

## 📋 前置要求

- Go 1.21 或更高版本
- Node.js 18+ 和 npm/yarn
- PostgreSQL 14+
- Redis 7+
- Make（可选，用于运行脚本）

## 🚀 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/hanyouqing/OpenAuth.git
cd OpenAuth
```

### 2. 配置数据库

创建 PostgreSQL 数据库：

```bash
createdb openauth
```

### 3. 配置应用

复制配置文件并修改：

```bash
cp configs/config.yaml configs/config.local.yaml
```

编辑 `configs/config.local.yaml`，设置数据库和 Redis 连接信息。

### 4. 运行后端

```bash
# 安装依赖
go mod download

# 运行数据库迁移
go run cmd/server/main.go migrate

# 启动服务器
go run cmd/server/main.go
```

后端服务将在 `http://localhost:8080` 启动。

### 5. 运行前端

```bash
cd web

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端应用将在 `http://localhost:3000` 启动。

## 🔐 默认账户

系统会自动创建默认管理员账户：

- **用户名**：`admin`
- **密码**：`admin123`

**⚠️ 重要**：首次登录后请立即修改密码！

## 📖 API 文档

### Swagger 文档

启动服务后，访问 Swagger UI 查看完整的 API 文档：
```
http://localhost:8080/swagger/index.html
```

生成 Swagger 文档：
```bash
make swagger
```

### 主要 API 端点

#### 认证接口
- `POST /api/v1/auth/login` - 登录
- `POST /api/v1/auth/refresh` - 刷新 Token
- `POST /api/v1/auth/register` - 注册
- `POST /api/v1/auth/forgot-password` - 忘记密码
- `POST /api/v1/auth/reset-password` - 重置密码

#### 用户管理
- `GET /api/v1/users` - 用户列表
- `GET /api/v1/users/me` - 当前用户信息
- `POST /api/v1/users` - 创建用户
- `PUT /api/v1/users/:id` - 更新用户
- `DELETE /api/v1/users/:id` - 删除用户

#### 应用管理
- `GET /api/v1/applications` - 应用列表
- `POST /api/v1/applications` - 创建应用
- `PUT /api/v1/applications/:id` - 更新应用
- `DELETE /api/v1/applications/:id` - 删除应用

#### MFA 管理
- `GET /api/v1/mfa/devices` - MFA 设备列表
- `POST /api/v1/mfa/devices/totp` - 创建 TOTP 设备
- `POST /api/v1/mfa/devices/totp/verify` - 验证 TOTP
- `POST /api/v1/mfa/devices/sms` - 发送 SMS 验证码
- `POST /api/v1/mfa/devices/email` - 发送 Email 验证码

#### 角色权限
- `GET /api/v1/roles` - 角色列表
- `POST /api/v1/roles` - 创建角色
- `POST /api/v1/roles/:id/permissions` - 分配权限
- `GET /api/v1/permissions` - 权限列表

#### 会话管理
- `GET /api/v1/sessions` - 会话列表
- `DELETE /api/v1/sessions/:id` - 删除会话
- `GET /api/v1/sessions/active/count` - 活跃会话数

#### 组织架构
- `GET /api/v1/organizations` - 组织列表
- `POST /api/v1/organizations` - 创建组织
- `POST /api/v1/organizations/:id/users` - 添加用户到组织
- `GET /api/v1/groups` - 用户组列表
- `POST /api/v1/groups` - 创建用户组

#### SSO 协议
- `GET /oauth2/authorize` - OAuth 2.0 授权端点
- `POST /oauth2/token` - OAuth 2.0 Token 端点
- `GET /oauth2/userinfo` - OIDC UserInfo 端点
- `GET /saml/sso` - SAML SSO 端点
- `GET /saml/metadata` - SAML Metadata 端点

更多 API 文档请参考 [API 文档](./docs/API.md)。

## 🐳 Docker 部署

### 使用 Docker Compose

```bash
docker-compose up -d
```

这将启动：
- PostgreSQL 数据库
- Redis 缓存
- OpenAuth 后端服务
- OpenAuth 前端服务

### 环境变量

可以通过环境变量配置应用：

```bash
export OPENAUTH_DATABASE_HOST=localhost
export OPENAUTH_DATABASE_PORT=5432
export OPENAUTH_DATABASE_USER=postgres
export OPENAUTH_DATABASE_PASSWORD=postgres
export OPENAUTH_DATABASE_DBNAME=openauth
export OPENAUTH_REDIS_HOST=localhost
export OPENAUTH_REDIS_PORT=6379
export OPENAUTH_JWT_SECRET=your-secret-key
```

## 📁 项目结构

```
openauth/
├── cmd/
│   └── server/              # 应用入口
├── internal/
│   ├── config/              # 配置管理
│   ├── database/            # 数据库连接和迁移
│   ├── models/              # 数据模型
│   ├── handlers/            # HTTP 处理器
│   ├── services/            # 业务逻辑
│   ├── middleware/          # 中间件
│   ├── auth/                # 认证相关
│   └── sso/                 # SSO 协议实现
├── web/                     # 前端代码
├── migrations/              # 数据库迁移文件
├── docs/                    # 文档
│   ├── PRD.md              # 产品需求文档
│   └── ARCHITECTURE.md      # 架构文档
├── configs/                 # 配置文件
└── README.md
```

## 🔧 开发

### 后端开发

```bash
# 安装依赖
make deps

# 运行服务
make run

# 运行测试
make test

# 代码格式化
make fmt

# 构建
make build

# 生成 Swagger 文档
make swagger
```

### 前端开发

```bash
cd web

# 安装依赖
npm install

# 运行开发服务器
npm run dev

# 运行测试
npm test

# 代码格式化
npm run format

# 构建生产版本
npm run build
```

## 🧪 测试

```bash
# 后端测试
go test ./...

# 前端测试
cd web && npm test
```

## 📝 许可证

本项目采用 MIT 许可证。详见 [LICENSE](./LICENSE) 文件。

## 🤝 贡献

欢迎贡献代码！请遵循以下步骤：

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 文档

- [产品需求文档 (PRD)](./docs/PRD.md)
- [系统架构文档](./docs/ARCHITECTURE.md)
- [API 文档](./docs/API.md)
- [实现状态文档](./docs/IMPLEMENTATION_STATUS.md)
- [Swagger 文档](http://localhost:8080/swagger/index.html)（启动服务后访问）

## 🗺️ 路线图

### Phase 1: 基础功能（进行中）
- [x] 用户认证和管理
- [x] 应用管理
- [x] 基础 MFA（TOTP）
- [ ] OAuth 2.0/OIDC 支持
- [ ] SAML 2.0 支持
- [ ] LDAP 支持

### Phase 2: 管理功能
- [ ] 角色权限管理
- [ ] 安全策略配置
- [ ] 审计日志查看
- [ ] 系统监控

### Phase 3: 高级功能
- [ ] 组织架构管理
- [ ] 条件访问策略
- [ ] Webhook 集成
- [ ] API 密钥管理

## 🐛 问题反馈

如果发现问题，请在 [GitHub Issues](https://github.com/hanyouqing/OpenAuth/issues) 提交。

## 🙏 致谢

本项目参考了以下优秀项目：
- [Okta](https://www.okta.com/) - 身份管理平台
- [Authing](https://www.authing.co/) - 身份云服务
- [Keycloak](https://www.keycloak.org/) - 开源身份管理

## 📧 联系方式

如有问题或建议，请通过以下方式联系：

- GitHub Issues: [提交问题](https://github.com/hanyouqing/OpenAuth/issues)
- Email: [待添加]

---

**OpenAuth** - 开源身份认证平台，让身份管理更简单。
