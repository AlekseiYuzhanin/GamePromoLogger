package logger

type Format string

const (
	JSONFormat Format = "json"
	TextFormat Format = "text"
)

type Config struct {
	Level       string `env:"LOG_LEVEL" default:"info"`
	Format      `env:"LOG_FORMAT" default:"json"`
	OutputPath  string `env:"OUTPUT_PATH" default:""`
	ServiceName string `env:"SERVICE_NAME" required:"true"`
}

func NewConfig(serviceName, level, outputPath string) *Config {
	return &Config{
		ServiceName: serviceName,
		Level:       level,
		OutputPath:  outputPath,
	}
}
