package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/hahamx/utools/utils/env"
	"github.com/hahamx/utools/utils/timeutil"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	// DefaultLevel the default log level
	DefaultLevel = zapcore.InfoLevel

	// DefaultTimeLayout the default time layout;
	DefaultTimeLayout = time.RFC3339
)

var (
	BaseLog     = "./logs/"
	ProjectName = "strtech-api"

	ProjectAccessLogFile = BaseLog + ProjectName + "-access.log"
	MonthDate            = map[string]int{
		"December":  12,
		"November":  11,
		"October":   10,
		"September": 9,
		"August":    8,
		"July":      7,
		"June":      6,
		"May":       5,
		"April":     4,
		"March":     3,
		"February":  2,
		"January":   1,
	}
)

// Option custom setup config
type Option func(*option)

type option struct {
	level          zapcore.Level
	fields         map[string]string
	file           io.Writer
	timeLayout     string
	disableConsole bool
}

// WithDebugLevel only greater than 'level' will output
func WithDebugLevel() Option {
	return func(opt *option) {
		opt.level = zapcore.DebugLevel
	}
}

// WithInfoLevel only greater than 'level' will output
func WithInfoLevel() Option {
	return func(opt *option) {
		opt.level = zapcore.InfoLevel
	}
}

// WithWarnLevel only greater than 'level' will output
func WithWarnLevel() Option {
	return func(opt *option) {
		opt.level = zapcore.WarnLevel
	}
}

// WithErrorLevel only greater than 'level' will output
func WithErrorLevel() Option {
	return func(opt *option) {
		opt.level = zapcore.ErrorLevel
	}
}

// WithField add some field(s) to log
func WithField(key, value string) Option {
	return func(opt *option) {
		opt.fields[key] = value
	}
}

// WithFileP write log to some file
func WithFileP(file string) Option {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0766); err != nil {
		panic(err)
	}

	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0766)
	if err != nil {
		panic(err)
	}

	return func(opt *option) {
		opt.file = zapcore.Lock(f)
	}
}

// WithFileRotationP write log to some file with rotation
func WithFileRotationP(file string) Option {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0766); err != nil {
		panic(err)
	}

	return func(opt *option) {
		opt.file = &lumberjack.Logger{ // concurrent-safed
			Filename:   file, // ????????????
			MaxSize:    128,  // ??????????????????????????????????????? M
			MaxBackups: 300,  // ???????????? 300 ?????????
			MaxAge:     30,   // ??????????????????????????? day
			LocalTime:  true, // ??????????????????
			Compress:   true, // ???????????? disabled by default
		}
	}
}

// WithTimeLayout custom time format
func WithTimeLayout(timeLayout string) Option {
	return func(opt *option) {
		opt.timeLayout = timeLayout
	}
}

// WithDisableConsole WithEnableConsole write log to os.Stdout or os.Stderr
func WithDisableConsole() Option {
	return func(opt *option) {
		opt.disableConsole = true
	}
}

// NewJSONLogger return a json-encoder zap logger,
func NewJSONLogger(opts ...Option) (*zap.Logger, error) {
	opt := &option{level: DefaultLevel, fields: make(map[string]string)}
	for _, f := range opts {
		f(opt)
	}

	timeLayout := DefaultTimeLayout
	if opt.timeLayout != "" {
		timeLayout = opt.timeLayout
	}

	// similar to zap.NewProductionEncoderConfig()
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger", // used by logger.Named(key); optional; useless
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace", // use by zap.AddStacktrace; optional; useless
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // ???????????????
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(timeLayout))
		},
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // ??????????????????
	}

	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// lowPriority usd by info\debug\warn
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= opt.level && lvl < zapcore.ErrorLevel
	})

	// highPriority usd by error\panic\fatal
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= opt.level && lvl >= zapcore.ErrorLevel
	})

	stdout := zapcore.Lock(os.Stdout) // lock for concurrent safe
	stderr := zapcore.Lock(os.Stderr) // lock for concurrent safe

	core := zapcore.NewTee()

	if !opt.disableConsole {
		core = zapcore.NewTee(
			zapcore.NewCore(jsonEncoder,
				zapcore.NewMultiWriteSyncer(stdout),
				lowPriority,
			),
			zapcore.NewCore(jsonEncoder,
				zapcore.NewMultiWriteSyncer(stderr),
				highPriority,
			),
		)
	}

	if opt.file != nil {
		core = zapcore.NewTee(core,
			zapcore.NewCore(jsonEncoder,
				zapcore.AddSync(opt.file),
				zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
					return lvl >= opt.level
				}),
			),
		)
	}

	logger := zap.New(core,
		zap.AddCaller(),
		zap.ErrorOutput(stderr),
	)

	for key, value := range opt.fields {
		logger = logger.WithOptions(zap.Fields(zapcore.Field{Key: key, Type: zapcore.StringType, String: value}))
	}
	return logger, nil
}

var _ Meta = (*meta)(nil)

// Meta key-value
type Meta interface {
	Key() string
	Value() interface{}
	meta()
}

type meta struct {
	key   string
	value interface{}
}

func (m *meta) Key() string {
	return m.key
}

func (m *meta) Value() interface{} {
	return m.value
}

func (m *meta) meta() {}

// NewMeta create meat
func NewMeta(key string, value interface{}) Meta {
	return &meta{key: key, value: value}
}

// WrapMeta wrap meta to zap fields
func WrapMeta(err error, metas ...Meta) (fields []zap.Field) {
	capacity := len(metas) + 1 // namespace meta
	if err != nil {
		capacity++
	}

	fields = make([]zap.Field, 0, capacity)
	if err != nil {
		fields = append(fields, zap.Error(err))
	}

	fields = append(fields, zap.Namespace("meta"))
	for _, meta := range metas {
		fields = append(fields, zap.Any(meta.Key(), meta.Value()))
	}

	return
}

// ?????????????????????????????????
func NewLogger(TName string) *zap.SugaredLogger {

	y, month, d := time.Now().Date()
	minthInt := MonthDate[month.String()]
	days := fmt.Sprintf("%v%v%v/", y, minthInt, d)
	if TName == "" {
		TName = "default"
	}
	logFileName := BaseLog + days + TName + "_access.log"
	accessLogger, err := NewJSONLogger(
		WithDisableConsole(),
		WithField("configs", fmt.Sprintf("%s[%s]", TName, env.Active().Value())),
		WithTimeLayout(timeutil.CSTLayout),
		WithFileP(logFileName),
	)
	if err != nil {
		panic(err)
	}
	return accessLogger.Sugar()
}
