package main

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/gographics/imagick.v3/imagick"
	"imageconverter/constant"
	"imageconverter/controller"
	"os"
)

func main() {
	loggerSetup()

	fileInfo, err := os.Lstat("./")
	if err != nil {
		zap.S().Fatalf("failed to get file info: %v", err)
	}
	fileMode := fileInfo.Mode()
	unixPerms := fileMode & os.ModePerm
	err = os.MkdirAll(constant.UPLOAD_FILE_PATH, unixPerms)
	if err != nil {
		zap.S().Fatalf("failed to create directory: %v", err)
	}
	err = os.MkdirAll(constant.CONVERTED_FILE_PATH, unixPerms)
	if err != nil {
		zap.S().Fatalf("failed to create directory: %v", err)
	}

	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	formats := mw.QueryFormats("*")
	fmt.Printf("supported image format: %v", formats)

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
