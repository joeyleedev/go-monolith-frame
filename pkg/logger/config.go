package logger

// Config 日志配置
type Config struct {
	Level       string `yaml:"level"`        // 日志级别: debug, info, warn, error
	Filename    string `yaml:"filename"`     // 日志文件路径
	MaxSize     int    `yaml:"max_size"`     // 单个日志文件最大大小(MB)
	MaxBackups  int    `yaml:"max_backups"`  // 保留的旧日志文件最大数量
	MaxAge      int    `yaml:"max_age"`      // 保留旧日志文件的最大天数
	Compress    bool   `yaml:"compress"`     // 是否压缩日志文件
	Console     bool   `yaml:"console"`      // 是否同时输出到控制台
	EnableColor bool   `yaml:"enable_color"` // 控制台输出是否启用颜色
}
