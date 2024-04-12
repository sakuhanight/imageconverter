package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/gographics/imagick.v3/imagick"
	"imageconverter/controller"
)

func main() {
	loggerSetup()

	imagick.Initialize()
	defer imagick.Terminate()

	r := controller.Router()
	r.Run()

}

func loggerSetup() {
	config := zap.NewProductionConfig()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig = encoderConfig
	config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	logger, _ := config.Build(zap.AddCallerSkip(1))
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
}
