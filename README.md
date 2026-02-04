# Go Monolith Frame

一个基于 Gin + GORM 的 Go 单体应用框架，提供开箱即用的企业级功能组件。

## 特性

- **配置管理**: 支持 Viper + YAML 配置文件，支持 .env 环境变量，多环境切换（dev/prod）
- **日志系统**: 基于 zap 的结构化日志，支持日志轮转、彩色控制台输出
- **数据库**: GORM ORM，支持自动迁移、连接池管理、慢查询检测、自定义 Zap Logger
- **缓存**: 抽象缓存层，支持 Memory/Redis 切换，Cache-Aside 模式
- **认证**: JWT 令牌生成与解析，bcrypt 密码加密，Auth 中间件保护路由
- **中间件**: 请求日志、Panic 恢复、CORS 跨域处理、JWT 认证
- **错误处理**: 统一的错误响应格式，业务错误码映射（模块化错误码）
- **依赖注入**: 显式依赖传递，提升可测试性
- **优雅关闭**: 信号处理，平滑服务停机

## 技术栈

| 组件 | 技术选型 |
|------|---------|
| Web 框架 | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io/) |
| 数据库驱动 | [gorm.io/driver/mysql](https://github.com/go-gorm/mysql) |
| 日志 | [zap](https://github.com/uber-go/zap) |
| 日志轮转 | [lumberjack](https://github.com/natefinch/lumberjack) |
| 缓存 | [redis/go-redis](https://github.com/redis/go-redis) |
| 配置 | [Viper](https://github.com/spf13/viper) |
| 环境变量 | [godotenv](https://github.com/joho/godotenv) |
| JWT | [golang-jwt/jwt](https://github.com/golang-jwt/jwt) |
| 密码加密 | [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto/bcrypt) |

## 快速开始

### 环境要求

- Go 1.23+
- MySQL 5.7+ / 8.0+
- Redis (可选，生产环境建议使用)

### 安装依赖

```bash
go mod download
```

### 配置环境变量

1. 复制 `.env.example` 为 `.env` 并填写实际值：
```bash
cp .env.example .env
```

2. 编辑 `.env` 文件：
```env
# 数据库连接
DATABASE_DSN=root:password@tcp(127.0.0.1:3306)/go_monolith?charset=utf8mb4&parseTime=True&loc=Local

# Redis 密码（如果 Redis 没有设置密码，留空即可）
CACHE_REDIS_PASSWORD=

# JWT 密钥（生产环境请使用强随机密钥）
JWT_SECRET=your-secret-key-change-in-production
```

### 运行服务

```bash
# 开发环境（默认）
go run cmd/server/main.go -env dev

# 生产环境
go run cmd/server/main.go -env prod
```

### 编译运行

```bash
go build -o bin/server cmd/server/main.go
./bin/server -env dev
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包测试
go test ./pkg/mysql/... -v
go test ./pkg/logger/... -v
go test ./internal/repository/... -v

# 检查代码
go vet ./...
go fmt ./...
```

## API 文档

### 基础 URL

```
http://localhost:8080/api/v1
```

### 公开路由（无需认证）

#### 健康检查
```bash
curl http://localhost:8080/health
```

#### 用户登录
```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456"
  }'
```

#### 创建用户
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "123456"
  }'
```

### 认证路由（需要 JWT Token）

所有认证路由需要在请求头中携带 JWT Token：
```
Authorization: Bearer {your-jwt-token}
```

#### 获取用户列表（分页）
```bash
curl http://localhost:8080/api/v1/users?page=1&page_size=10&keyword=test \
  -H "Authorization: Bearer {token}"
```

#### 获取用户详情
```bash
curl http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer {token}"
```

#### 更新用户
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newusername",
    "email": "newemail@example.com"
  }'
```

#### 删除用户
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer {token}"
```

#### 修改密码
```bash
curl -X PATCH http://localhost:8080/api/v1/users/1/password \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "123456",
    "new_password": "newpassword"
  }'
```

## 项目结构

```
go-monolith-frame/
├── cmd/                  # 应用入口
│   └── server/
│       └── main.go       # 主程序入口
├── configs/              # 配置文件
│   ├── config.dev.yaml   # 开发环境配置
│   └── config.prod.yaml  # 生产环境配置
├── internal/             # 内部代码（业务逻辑）
│   ├── config/           # 配置加载与结构定义
│   ├── errcode/          # 业务错误码
│   ├── handler/          # HTTP 处理器（Controller 层）
│   ├── middleware/       # Gin 中间件
│   │   ├── auth.go       # JWT 认证中间件
│   │   ├── cors.go       # 跨域处理
│   │   ├── logger.go     # 请求日志
│   │   └── recovery.go   # Panic 恢复
│   ├── model/            # 数据模型
│   │   ├── entity/       # 数据库实体定义
│   │   ├── request/      # 请求 DTO（Request Body）
│   │   └── response/     # 响应 DTO（Response Body）
│   ├── repository/       # 数据访问层（DAO）
│   └── service/          # 业务逻辑层（Service）
├── pkg/                  # 可重用包（基础设施）
│   ├── cache/            # 缓存抽象层（Memory/Redis）
│   ├── errcode/          # 通用错误码定义
│   ├── logger/           # 日志封装（zap + lumberjack）
│   ├── mysql/            # 数据库连接管理
│   ├── response/         # 统一响应格式
│   └── utils/            # 工具函数
│       ├── jwt.go        # JWT 工具
│       └── password.go   # 密码加密/验证
├── .env.example          # 环境变量示例
├── go.mod                # Go 模块定义
├── go.sum                # Go 依赖校验
├── CLAUDE.md             # Claude Code 项目指南
└── README.md
```

## 配置说明

### 数据库配置

数据库连接信息通过环境变量 `DATABASE_DSN` 配置，或在 `configs/config.{env}.yaml` 中设置：

```yaml
database:
  dsn: "root:password@tcp(127.0.0.1:3306)/go_monolith?charset=utf8mb4&parseTime=True&loc=Local"
  max_open_conns: 10        # 最大连接数
  max_idle_conns: 5         # 最大空闲连接数
  conn_max_lifetime: 1h     # 连接最大存活时间
  conn_max_idle_time: 10m    # 连接最大空闲时间
  log_level: "info"          # 日志级别: warn / info
  slow_threshold: 500ms      # 慢查询阈值
  ignore_record_not_found_error: false  # 是否忽略记录未找到错误
```

### 缓存配置

支持两种缓存模式：

**Memory 缓存**（开发环境）：
```yaml
cache:
  type: "memory"
  memory:
    cleanup_interval: 10m
```

**Redis 缓存**（生产环境）：
```yaml
cache:
  type: "redis"
  redis:
    addr: "127.0.0.1:6379"
    password: ""             # 通过环境变量 CACHE_REDIS_PASSWORD 设置
    db: 0
```

### 日志配置

```yaml
log:
  level: "debug"            # 日志级别: debug / info / warn / error
  filename: "logs/app.log"
  max_size: 100             # 单文件最大 MB
  max_backups: 3            # 保留备份文件数
  max_age: 7                # 保留天数
  compress: true             # 是否压缩
  console: true              # 是否输出到控制台
  enable_color: true         # 是否启用颜色
```

### JWT 配置

```yaml
jwt:
  secret: ""                 # 通过环境变量 JWT_SECRET 设置
  expire_time: 24h           # Token 有效期
```

### 服务器配置

```yaml
server:
  addr: ":8080"              # 监听地址
  mode: "debug"              # Gin 模式: debug / release
  read_timeout: 60s          # 读取超时
  write_timeout: 60s         # 写入超时
```

## 开发指南

### 添加新的 API

1. 在 `internal/model/entity/` 定义实体（数据库表结构）
2. 在 `internal/model/request/` 定义请求 DTO（请求体参数）
3. 在 `internal/model/response/` 定义响应 DTO（响应体结构）
4. 在 `internal/repository/` 实现数据访问接口（继承 `*gorm.DB`）
5. 在 `internal/service/` 实现业务逻辑（注入 Repository 和 Cache）
6. 在 `internal/handler/` 实现 HTTP 处理器（注入 Service）
7. 在 `cmd/server/main.go` 中注册路由

示例依赖注入流程：
```go
// main.go 中
db, _ := mysql.New(&cfg.Database)
cacheClient, _ := cache.New(&cfg.Cache)

// 依赖注入链
userRepo := repository.NewUserRepository(db)
userService := service.NewUserService(userRepo, cacheClient)
userHandler := handler.NewUserHandler(userService)

// 注册路由
users := api.Group("/users")
users.POST("", userHandler.CreateUser)
users.GET("", userHandler.ListUsers)
```

### 添加新的中间件

在 `internal/middleware/` 目录下创建新的中间件文件，实现 `gin.HandlerFunc` 类型：

```go
package middleware

import "github.com/gin-gonic/gin"

func YourMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 前置处理
        c.Next()
        // 后置处理
    }
}
```

### 自定义错误码

业务错误码按模块分类，通用错误码在 `pkg/errcode/`，模块特定错误码在 `internal/errcode/`：

**pkg/errcode/common.go** - 通用错误码（0xxxx）：
```go
const (
    CodeSuccess        = 0
    CodeInvalidParam   = 40001
    CodeInternalError  = 50001
)
```

**internal/errcode/user.go** - 用户模块错误码（2xxxx）：
```go
const (
    CodeUserNotFound    = 20001
    CodeUserExists      = 20002
    CodeInvalidPassword = 20003
)
```

## 架构设计

### 分层架构

项目采用清晰的分层架构，每一层只依赖下层：

```
┌─────────────────────────────────────────────┐
│            Handler Layer (Controller)         │
│  处理 HTTP 请求/响应，参数绑定，路由注册      │
└─────────────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────┐
│            Service Layer (Business)          │
│  业务逻辑处理，调用 Repository 和 Cache      │
│  使用 Cache-Aside 缓存模式                   │
└─────────────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────┐
│         Repository Layer (Data Access)       │
│  数据访问层，使用 GORM 操作数据库            │
└─────────────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────┐
│              Infrastructure                 │
│  MySQL, Redis, Logger, Config, Utils      │
└─────────────────────────────────────────────┘
```

### 依赖注入模式

项目采用构造函数依赖注入，所有依赖通过构造函数显式传递：

```go
// Repository 层
func NewUserRepository(db *gorm.DB) UserRepository { ... }

// Service 层
func NewUserService(repo UserRepository, cache Cache) UserService { ... }

// Handler 层
func NewUserHandler(service UserService) *UserHandler { ... }
```

### 错误码约定

- `0xxxx` - 通用/系统错误
- `1xxxx` - 预留给认证模块
- `2xxxx` - 用户模块错误
- `3xxxx` - 订单模块错误
- ...

### 缓存策略

采用 Cache-Aside 模式：
- **读取**: 先查缓存，未命中则查数据库并写入缓存
- **写入**: 更新/删除时直接操作数据库，同时删除/失效缓存
- **缓存键格式**: `{entity}:{id}`，例如 `user:1`

## VS Code 调试

项目包含 VS Code 调试配置（`.vscode/launch.json`），可以直接按 `F5` 启动调试。

## 常见问题

### 1. 数据库连接失败

确保 MySQL 服务正在运行，并且 `.env` 中的 `DATABASE_DSN` 配置正确。

### 2. Redis 连接失败

如果使用 Redis 缓存，确保 Redis 服务正在运行，并检查 `CACHE_REDIS_PASSWORD` 配置。

### 3. 日志不显示

检查 `configs/config.{env}.yaml` 中的 `log.console` 和 `log.level` 配置。

### 4. JWT Token 无效

确保 `JWT_SECRET` 环境变量已设置，并且 Token 未过期。

## License

MIT License
