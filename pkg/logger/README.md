# Logger Package

## 快速开始

### 基础使用

```go
package main

import (
    "isvgo/pkg/logger"
    "go.uber.org/zap"
)

func main() {
    // 使用默认配置初始化
    cfg := logger.DefaultConfig()
    err := logger.Init(cfg)
    if err != nil {
        log.Fatal("failed to init logger:", err)
    }
    // 程序结束前同步日志缓冲区
    defer logger.Sync()
    
    // 记录日志
    logger.Info("应用启动", zap.String("version", "1.0.0"))
    logger.Error("发生错误", zap.Error(err))
}
```

### 自定义配置

```go
package main

import (
    "isvgo/pkg/logger"
    "go.uber.org/zap"
)

func main() {
    // 自定义配置
    cfg := &logger.Config{
        Level:            "debug",        // 日志级别
        Format:           "json",         // 输出格式
        Path:             "./logs/app.log", // 日志文件路径
        ConsoleEnabled:   true,           // 启用控制台输出
        EnableCaller:     true,           // 启用调用者信息
        EnableStacktrace: true,           // 启用堆栈跟踪
        MaxSize:          100,            // 单文件最大 100MB
        MaxAge:           7,              // 保留 7 天
        MaxBackups:       5,              // 最多保留 5 个备份文件
    }
    
    err := logger.Init(cfg)
    if err != nil {
        log.Fatal("failed to init logger:", err)
    }
    
    // 使用结构化字段记录日志
    logger.Info("用户登录",
        zap.String("username", "john"),
        zap.String("ip", "192.168.1.100"),
        zap.Duration("response_time", time.Millisecond*50),
    )
}
```

## 配置说明

### Config 结构体

```go
type Config struct {
    // 基础配置
    Level  string // 日志级别: debug, info, warn, error
    Format string // 输出格式: json, console
    
    // 输出配置
    Path           string // 日志文件路径
    ConsoleEnabled bool   // 是否启用控制台输出
    
    // 调试选项
    EnableCaller     bool // 是否启用调用者信息
    EnableStacktrace bool // 是否启用堆栈跟踪
    
    // 文件轮转配置 (lumberjack)
    MaxSize    int // 单个日志文件最大大小 (MB)
    MaxAge     int // 日志文件保留天数
    MaxBackups int // 最大备份文件数量
}
```

### 默认配置

```go
cfg := logger.DefaultConfig()
// 等价于:
cfg := &logger.Config{
    Level:            "info",
    Format:           "console",
    ConsoleEnabled:   true,
    EnableCaller:     true,
    EnableStacktrace: false,
}
```

### 与 Viper 集成

```go
// configs/config.go
package configs

import "isvgo/pkg/logger"

type AppConfig struct {
	Log logger.Config `mapstructure:"log"`
}
```



```go
// main.go
package main

import (
	"flag"
	"log"
	"isvgo/configs"
	"isvgo/pkg/logger"

	"github.com/spf13/viper"
)

func main() {
	var cfgFile string
	flag.StringVar(&cfgFile, "config", "./app.toml", "配置文件(.toml)")
	flag.StringVar(&cfgFile, "c", "./app.toml", "配置文件(.toml)")
	flag.Parse()

	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("failed to read config file:", err)
	}

	var appConfig configs.AppConfig
	if err := viper.Unmarshal(&appConfig); err != nil {
		log.Fatal("failed to unmarshal config:", err)
	}

	err := logger.Init(&appConfig.Log)
    if err != nil {
		log.Fatal("failed to init logger:", err)
	}
    defer logger.Sync()

	logger.Info("logger initialized")
}
```

## API 参考

### 初始化函数

- `Init(config *Config) error` - 使用配置初始化全局日志器
- `GetLogger() *zap.Logger` - 获取全局日志器实例
- `Sync() error` - 同步日志缓冲区

### 日志记录函数

所有日志函数都支持 zap.Field 结构化字段：

- `Debug(msg string, fields ...zap.Field)` - 调试级别日志
- `Info(msg string, fields ...zap.Field)` - 信息级别日志  
- `Warn(msg string, fields ...zap.Field)` - 警告级别日志
- `Error(msg string, fields ...zap.Field)` - 错误级别日志
- `Fatal(msg string, fields ...zap.Field)` - 致命错误日志 (会终止程序)

### 工厂函数

- `NewZapLogger(cfg *Config) (*zap.Logger, error)` - 创建新的 zap 日志器实例

## 使用示例

### 1. 基础日志记录

```go
logger.Debug("调试信息")
logger.Info("信息日志")
logger.Warn("警告信息") 
logger.Error("错误信息")
```

### 2. 结构化日志

```go
logger.Info("received HTTP request",
    zap.String("method", "GET"),
    zap.String("url", "/api/users"),
    zap.Int("status", 200),
    zap.Duration("latency", time.Millisecond*120),
)
```

### 3. 错误日志

```go
if err != nil {
    logger.Error("failed to connect to database",
        zap.Error(err),
        zap.String("database", "mysql"),
        zap.String("host", "localhost:3306"),
    )
}
```

### 4. 仅文件输出

```go
cfg := &logger.Config{
    Level:            "info",
    Format:           "json",
    Path:             "./log/app.log",
    ConsoleEnabled:   false,  // 禁用控制台输出
    EnableCaller:     true,  // 启用调用者信息
    EnableStacktrace: false,  // 隐藏堆栈信息
    MaxSize:          50,
    MaxAge:           30,
    MaxBackups:       10,
}
```

### 5. 仅控制台输出

```go
cfg := &logger.Config{
    Level:            "debug",
    Format:           "console",
    EnableCaller:     true,
    EnableStacktrace: true,
}
```

### 6. 获取原生 zap.Logger

```go
zapLogger := logger.GetLogger()
zapLogger.Info("使用原生 zap API",
    zap.String("component", "auth"),
    zap.Any("config", cfg),
)
```

## 输出格式

### Console 格式

```
2023-10-24 10:30:45.123 INFO    main.go:15      应用启动    {"version": "1.0.0"}
2023-10-24 10:30:45.124 ERROR   main.go:20      数据库连接失败  {"error": "connection refused"}
```

### JSON 格式

```json
{"time":"2023-10-24T10:30:45.123Z","level":"info","caller":"main.go:15","msg":"应用启动","version":"1.0.0"}
{"time":"2023-10-24T10:30:45.124Z","level":"error","caller":"main.go:20","msg":"数据库连接失败","error":"connection refused"}
```

## 文件轮转

当启用文件输出时，日志库会自动处理文件轮转：

- 当日志文件大小超过 `MaxSize` (MB) 时自动轮转
- 保留最近 `MaxBackups` 个备份文件
- 删除超过 `MaxAge` 天的旧日志文件
- 自动压缩轮转的日志文件

轮转后的文件命名格式：`app.log.2023-10-24T10-30-45.123.gz`

## 性能考虑

- 基于 Uber Zap，提供极高的性能
- 使用结构化字段避免字符串格式化开销
- 内置日志缓冲，减少 I/O 操作
- 程序退出前调用 `logger.Sync()` 确保日志完整性

## 线程安全

所有 API 都是线程安全的，可以在多个 goroutine 中安全使用。

## 测试

运行测试：

```bash
go test ./pkg/logger
```

## 依赖

- [go.uber.org/zap](https://github.com/uber-go/zap) - 高性能结构化日志库
- [gopkg.in/natefinch/lumberjack.v2](https://github.com/natefinch/lumberjack) - 日志文件轮转

