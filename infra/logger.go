package infra

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (s *Server) setLogger() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	atom := zap.NewAtomicLevel()
	switch s.config.GetString("log.level") {
	case "debug", "debugging", "deb", "debag":
		atom.SetLevel(zap.DebugLevel)
	case "info", "information", "inf":
		atom.SetLevel(zap.InfoLevel)
	case "warn", "warning", "WARN":
		atom.SetLevel(zap.WarnLevel)
	case "err", "error":
		atom.SetLevel(zap.ErrorLevel)
	default:
		atom.SetLevel(zap.InfoLevel)

	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	//encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	),
		zap.AddCaller()).With(
		zap.String("program", s.config.GetString("app.name")),
		zap.String("hostname", hostname))

	s.log = logger.Sugar()
}
