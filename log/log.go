package log

type TypeEnum string

const (
	TypeGet  TypeEnum = "get"
	TypePost          = "post"
	TypeApp           = "app"
)

type Level string

const (
	TraceLevel Level = "trace"
	DebugLevel       = "debug"
	InfoLevel        = "info"
	WarnLevel        = "warn"
	ErrorLevel       = "error"
	FatalLevel       = "fatal"
	PanicLevel       = "panic"
)

type Config struct {
	LogLevel Level
	Debug    bool
	LogMode  uint32
	LogDir   string
}

type Logger interface {
	Errorf(t TypeEnum, format string, args ...interface{})
	Warnf(t TypeEnum, format string, args ...interface{})
	Debugf(t TypeEnum, format string, args ...interface{})
	Infof(t TypeEnum, format string, args ...interface{})
	Fatalf(t TypeEnum, format string, args ...interface{})
	Close()
}

func GetLogTypeByRequestType(rType string) TypeEnum {
	lType := TypeGet
	if rType == "POST" {
		lType = TypePost
	}
	return lType
}

func NewLog(conf *Config) (Logger, error) {
	log := Zerolog{conf: conf}
	err := log.init()
	if err != nil {
		return nil, err
	}
	return &log, nil
}
