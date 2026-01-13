<div align="center">

![OpenAuth Logo](./web/public/logo.svg)

# OpenAuth

**An open-source Identity and Access Management (IAM) platform providing enterprise-grade Single Sign-On (SSO), Multi-Factor Authentication (MFA), user management, and application integration capabilities.**

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://golang.org/)
[![React](https://img.shields.io/badge/React-18+-61DAFB?logo=react)](https://reactjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.3+-3178C6?logo=typescript)](https://www.typescriptlang.org/)

[English](./README.md) | [中文](./README_zh.md)

</div>



## ✨ Features

### Core Capabilities
- ✅ **User Authentication**: Login, registration, password management
- ✅ **Multi-Factor Authentication (MFA)**: Support for TOTP, SMS, Email
- ✅ **User Management**: Complete user CRUD operations
- ✅ **Application Management**: Application creation, configuration, and management
- ✅ **Role-Based Access Control (RBAC)**: Role and permission management
- ✅ **Security Policies**: Password policies, MFA policies, whitelist management
- ✅ **Audit Logging**: Complete operation audit records
- ✅ **Session Management**: Session viewing, deletion, and statistics
- ✅ **Organization Management**: Organization and user group management
- ✅ **Notification Services**: Email and SMS notifications
- ✅ **API Documentation**: Swagger documentation support

### SSO Protocol Support
- ✅ **OAuth 2.0 / OIDC**: Standard OAuth 2.0 and OpenID Connect protocols
- ✅ **SAML 2.0**: Enterprise-grade SAML Single Sign-On (SSO, Metadata)
- ✅ **LDAP**: LDAP directory service integration

### Technical Highlights
- 🎨 **Modern UI**: Based on React + Ant Design with minimalist black and white theme
- 🚀 **High Performance**: Go backend, PostgreSQL + Redis
- 🔒 **Secure & Reliable**: JWT authentication, bcrypt password encryption, MFA support
- 📦 **Easy Deployment**: Docker support, simple configuration
- 🌐 **Internationalization**: Multi-language support (i18n)

## 🏗️ Tech Stack

### Backend
- **Language**: Go 1.23+
- **Framework**: Gin
- **Database**: PostgreSQL 14+
- **Cache**: Redis 7+
- **ORM**: GORM
- **Authentication**: JWT

### Frontend
- **Framework**: React 18 + TypeScript
- **Build Tool**: Vite
- **UI Component Library**: Ant Design 5
- **State Management**: Zustand
- **HTTP Client**: Axios
- **Data Fetching**: TanStack Query

## 📋 Prerequisites

- Go 1.23 or higher
- Node.js 18+ and npm/yarn
- PostgreSQL 14+
- Redis 7+
- Make (optional, for running scripts)

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/hanyouqing/OpenAuth.git
cd OpenAuth
```

### 2. Configure Database

Create a PostgreSQL database:

```bash
createdb openauth
```

### 3. Configure Application

Copy the configuration file and modify it:

```bash
cp configs/config.yaml configs/config.local.yaml
```

Edit `configs/config.local.yaml` and set the database and Redis connection information.

### 4. Run Backend

```bash
# Install dependencies
go mod download

# Run database migrations
go run cmd/server/main.go migrate

# Start the server
go run cmd/server/main.go
```

The backend service will start at `http://localhost:8080`.

### 5. Run Frontend

```bash
cd web

# Install dependencies
npm install

# Start development server
npm run dev
```

The frontend application will start at `http://localhost:3000`.

## 🔐 Default Account

The system automatically creates a default administrator account:

- **Username**: `admin`
- **Password**: `admin123`

**⚠️ Important**: Please change the password immediately after first login!

## 📖 API Documentation

### Swagger Documentation

After starting the service, visit Swagger UI to view the complete API documentation:
```
http://localhost:8080/swagger/index.html
```

Generate Swagger documentation:
```bash
make swagger
```

### Main API Endpoints

#### Authentication
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh Token
- `POST /api/v1/auth/register` - Register
- `POST /api/v1/auth/forgot-password` - Forgot Password
- `POST /api/v1/auth/reset-password` - Reset Password

#### User Management
- `GET /api/v1/users` - User list
- `GET /api/v1/users/me` - Current user information
- `POST /api/v1/users` - Create user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

#### Application Management
- `GET /api/v1/applications` - Application list
- `POST /api/v1/applications` - Create application
- `PUT /api/v1/applications/:id` - Update application
- `DELETE /api/v1/applications/:id` - Delete application

#### MFA Management
- `GET /api/v1/mfa/devices` - MFA device list
- `POST /api/v1/mfa/devices/totp` - Create TOTP device
- `POST /api/v1/mfa/devices/totp/verify` - Verify TOTP
- `POST /api/v1/mfa/devices/sms` - Send SMS verification code
- `POST /api/v1/mfa/devices/email` - Send Email verification code

#### Role & Permissions
- `GET /api/v1/roles` - Role list
- `POST /api/v1/roles` - Create role
- `POST /api/v1/roles/:id/permissions` - Assign permissions
- `GET /api/v1/permissions` - Permission list

#### Session Management
- `GET /api/v1/sessions` - Session list
- `DELETE /api/v1/sessions/:id` - Delete session
- `GET /api/v1/sessions/active/count` - Active session count

#### Organization
- `GET /api/v1/organizations` - Organization list
- `POST /api/v1/organizations` - Create organization
- `POST /api/v1/organizations/:id/users` - Add user to organization
- `GET /api/v1/groups` - User group list
- `POST /api/v1/groups` - Create user group

#### SSO Protocols
- `GET /oauth2/authorize` - OAuth 2.0 authorization endpoint
- `POST /oauth2/token` - OAuth 2.0 token endpoint
- `GET /oauth2/userinfo` - OIDC UserInfo endpoint
- `GET /saml/sso` - SAML SSO endpoint
- `GET /saml/metadata` - SAML Metadata endpoint

For more API documentation, please refer to [API Documentation](./docs/API.md).

## 🐳 Docker Deployment

### Using Docker Compose

```bash
docker-compose up -d
```

This will start:
- PostgreSQL database
- Redis cache
- OpenAuth backend service
- OpenAuth frontend service

### Environment Variables

You can configure the application through environment variables:

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

## 📁 Project Structure

```
openauth/
├── cmd/
│   └── server/              # Application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── database/            # Database connection and migrations
│   ├── models/              # Data models
│   ├── handlers/            # HTTP handlers
│   ├── services/            # Business logic
│   ├── middleware/          # Middleware
│   ├── auth/                # Authentication related
│   └── sso/                 # SSO protocol implementation
├── web/                     # Frontend code
├── migrations/              # Database migration files
├── docs/                    # Documentation
│   ├── PRD.md              # Product Requirements Document
│   └── ARCHITECTURE.md      # Architecture documentation
├── configs/                 # Configuration files
└── README.md
```

## 🔧 Development

### Backend Development

```bash
# Install dependencies
make deps

# Run service
make run

# Run tests
make test

# Format code
make fmt

# Build
make build

# Generate Swagger documentation
make swagger
```

### Frontend Development

```bash
cd web

# Install dependencies
npm install

# Run development server
npm run dev

# Run tests
npm test

# Format code
npm run format

# Build for production
npm run build
```

## 🧪 Testing

```bash
# Backend tests
go test ./...

# Frontend tests
cd web && npm test

# Test with coverage
make test-coverage
```

For more testing information, please refer to [Testing Documentation](./docs/TESTING.md).

## 📝 License

This project is licensed under the **Apache License 2.0**. See the [LICENSE](./LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the project
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

Please read our [Contributing Guidelines](./CONTRIBUTING.md) (if available) for more details.

## 📄 Documentation

- [Product Requirements Document (PRD)](./docs/PRD.md)
- [System Architecture Documentation](./docs/ARCHITECTURE.md)
- [API Documentation](./docs/API.md)
- [Implementation Status](./docs/IMPLEMENTATION_STATUS.md)
- [Testing Documentation](./docs/TESTING.md)
- [Swagger Documentation](http://localhost:8080/swagger/index.html) (available after starting the service)

## 🗺️ Roadmap

### Phase 1: Core Features (In Progress)
- [x] User authentication and management
- [x] Application management
- [x] Basic MFA (TOTP)
- [x] OAuth 2.0/OIDC support
- [x] SAML 2.0 support
- [x] LDAP support

### Phase 2: Management Features
- [x] Role and permission management
- [x] Security policy configuration
- [x] Audit log viewing
- [ ] System monitoring

### Phase 3: Advanced Features
- [x] Organization management
- [x] Conditional access policies
- [x] Webhook integration
- [x] API key management

## 🐛 Issue Reporting

If you find any issues, please submit them on [GitHub Issues](https://github.com/hanyouqing/OpenAuth/issues).

## 🙏 Acknowledgments

This project references the following excellent projects:
- [Okta](https://www.okta.com/) - Identity management platform
- [Authing](https://www.authing.co/) - Identity cloud service
- [Keycloak](https://www.keycloak.org/) - Open source identity management

## 📧 Contact

For questions or suggestions, please contact us through:

- GitHub Issues: [Submit an Issue](https://github.com/hanyouqing/OpenAuth/issues)
- Email: [To be added]

---

<div align="center">

**OpenAuth** - Open-source identity authentication platform, making identity management simpler.

Made with ❤️ by the OpenAuth community

[![GitHub stars](https://img.shields.io/github/stars/hanyouqing/OpenAuth?style=social)](https://github.com/hanyouqing/OpenAuth/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/hanyouqing/OpenAuth?style=social)](https://github.com/hanyouqing/OpenAuth/network/members)

</div>
