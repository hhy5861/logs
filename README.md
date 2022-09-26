### 日志采用 [uber zap](https://github.com/uber-go/zap) 进行二次封装，包含日志自动清理，切割等功能。自动切割与清理采用[lumberjack](https://github.com/natefinch/lumberjack) 。

#### 快速使用

1. config.yaml 配置

```yaml
logger:
  level: debug            # 日志级别(debug|info|warn|error|dpanic|panic|fatal)
  stackLevel: error       # 打印栈级别，默认是error
  lumberjack:
    filename: ./app.log   # 日志写入目录（包含日志文件名）
    maxsize: 10           # 每个日志文件大小，单位M
    maxage: 1             # 最大保留天数
    maxbackups: 1         # 最多保存日志文件数
    compress: true        # 是否压缩日志文件；（gzip压缩格式）
```

2. 添加配置文件解析，配置结构体放到项目结构体解析。配置结构体已经定义，可直接使用

```go
type(
    Config struct {
        Logger *logs.StoreConfig `yaml:"logger"`
    }
)
```

3. 使用示例
```go
logger := logs.NewFactory(logs.NewStore(cfg.Logger).JsonEncoder())

1. 带context使用方式, 会添加traceId字段
logger.For(ctx.Request.Context()).Error("[ServiceCodeException]", zap.String("uri", ctx.Request.URL.Path), zap.Error(err))

2. 无context使用方式
logger.Bg().Error("[ServiceCodeException]", zap.String("uri", ctx.Request.URL.Path), zap.Error(err))
```

4. 日志实现接口
```go
type Logger interface {
    Info(msg string, fields ...zapcore.Field)
    
    Debug(msg string, fields ...zapcore.Field)
    
    Warn(msg string, fields ...zapcore.Field)
    
    Error(msg string, fields ...zapcore.Field)
    
    Fatal(msg string, fields ...zapcore.Field)
    
    With(fields ...zapcore.Field) Logger

    Println(v ...interface{})
}
```