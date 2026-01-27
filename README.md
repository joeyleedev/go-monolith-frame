# Go Monolith Frame

一个基于 Gin + GORM 的 Go 单体应用框架，提供开箱即用的企业级功能组件。

## 特性

- **配置管理**: 支持 YAML 配置文件，多环境切换（dev/prod）
- **日志系统**: 基于 zap 的结构化日志，支持日志轮转
- **数据库**: GORM ORM，支持自动迁移、连接池管理
- **缓存**: 抽象缓存层，支持 Memory/Redis 切换
- **认证**: JWT 令牌生成与解析，bcrypt 密码加密
- **中间件**: 请求日志、Panic 恢复、CORS 跨域处理
- **错误处理**: 统一的错误响应格式，业务错误码映射
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
| 配置 | [yaml.v3](https://github.com/go-yaml/yaml) |
| JWT | [golang-jwt/jwt](https://github.com/golang-jwt/jwt) |
| 密码加密 | [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto/bcrypt) |

## 快速开始

### 环境要求

- Go 1.24+
- MySQL 5.7+ / 8.0+
- Redis (可选，生产环境建议使用)

### 安装依赖

```bash
go mod download
```

### 配置数据库

修改 `configs/config.dev.yaml` 中的数据库连接信息：

```yaml
database:
  dsn: "root:password@tcp(127.0.0.1:3306)/go_monolith?charset=utf8mb4&parseTime=True&loc=Local"
```

### 运行服务

```bash
# 开发环境
go run cmd/server/main.go --env=dev

# 生产环境
go run cmd/server/main.go --env=prod
```

### 编译运行

```bash
go build -buildvcs=false -o bin/server cmd/server/main.go
./bin/server --env=dev
```

## API 文档

### 健康检查

```bash
curl http://localhost:8080/health
```

### 创建用户

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "123456"
  }'
```

### 获取用户

```bash
curl http://localhost:8080/api/v1/users/1
```

## 项目结构

```
go-monolith-frame/
├── api/                  # API 定义文件
├── cmd/                  # 应用入口
│   └── server/
│       └── main.go       # 主程序入口
├── configs/              # 配置文件
│   ├── config.dev.yaml   # 开发环境配置
│   └── config.prod.yaml  # 生产环境配置
├── deployments/          # 部署相关文件
├── internal/             # 内部代码
│   ├── config/           # 配置加载
│   ├── handler/          # HTTP 处理器
│   ├── middleware/       # 中间件
│   │   ├── cors.go       # 跨域处理
│   │   ├── logger.go     # 请求日志
│   │   └── recovery.go   # Panic 恢复
│   ├── model/            # 数据模型
│   │   ├── entity/       # 实体定义
│   │   ├── request/      # 请求 DTO
│   │   └── response/     # 响应 DTO
│   ├── repository/       # 数据访问层
│   └── service/          # 业务逻辑层
├── pkg/                  # 可重用包
│   ├── cache/            # 缓存抽象
│   ├── errors/           # 错误定义
│   ├── logger/           # 日志封装
│   ├── mysql/            # 数据库连接
│   ├── response/         # 统一响应
│   └── utils/            # 工具函数
│       ├── jwt.go        # JWT 工具
│       └── password.go   # 密码工具
├── scripts/              # 脚本文件
├── logs/                 # 日志文件目录（自动生成）
└── README.md
```

## 配置说明

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
    password: ""
    db: 0
```

### 日志配置

```yaml
log:
  level: "debug"          # 日志级别: debug / info / warn / error
  filename: "logs/app.log"
  max_size: 100           # 单文件最大 MB
  max_backups: 3          # 保留备份文件数
  max_age: 7              # 保留天数
  compress: true          # 是否压缩
  console: true           # 是否输出到控制台
  enable_color: true      # 是否启用颜色
```

### JWT 配置

```yaml
jwt:
  secret: "your-secret-key"
  expire_time: 24h
```

## 开发指南

### 添加新的 API

1. 在 `internal/model/entity/` 定义实体
2. 在 `internal/model/request/` 定义请求 DTO
3. 在 `internal/model/response/` 定义响应 DTO
4. 在 `internal/repository/` 实现数据访问接口
5. 在 `internal/service/` 实现业务逻辑
6. 在 `internal/handler/` 实现 HTTP 处理器
7. 在 `cmd/server/main.go` 注册路由

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

在 `pkg/errors/errors.go` 中添加新的业务错误码：

```go
const (
    CodeYourError = 40010
)

var ErrYourError = &BizError{
    Code:    CodeYourError,
    Message: "自定义错误信息",
}

func (e *BizError) HTTPStatus() int {
    return 400
}
```

## License

MIT License
