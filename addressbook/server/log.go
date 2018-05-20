package main

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initLog() {
	// For console log
	consoleLogLevelStr := viper.GetString("log.console.level")
	consoleLogLevel, err := log.ParseLevel(consoleLogLevelStr)
	if err != nil {
		log.WithFields(log.Fields{
			"level": consoleLogLevelStr,
			"err":   err,
		}).Fatal("Failed to read log level for console")
	}
	log.SetLevel(consoleLogLevel)
	log.WithField("level", consoleLogLevelStr).Info("Set console log level")
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	// For file log
	path := viper.GetString("log.file.path")
	rotationTime := time.Duration(viper.GetInt64("log.file.rotationTime"))
	maxAge := time.Duration(viper.GetInt64("log.file.maxAge"))
	writer, err := rotatelogs.New(
		path+".%Y%m%d%H%M%S",
		rotatelogs.WithRotationTime(rotationTime*time.Second),
		rotatelogs.WithMaxAge(maxAge*time.Second),
	)
	if err != nil {
		log.WithError(err).Info("Failed to create logrotate writer")
	}

	fileLogLevelStr := viper.GetString("log.file.level")
	fileLogLevel, err := log.ParseLevel(fileLogLevelStr)
	if err != nil {
		log.WithFields(log.Fields{
			"level": fileLogLevelStr,
			"err":   err,
		}).Fatal("Failed to read log level for file")
	}

	writerMap := lfshook.WriterMap{}
	switch fileLogLevel {
	case log.DebugLevel:
		writerMap[log.DebugLevel] = writer
		fallthrough
	case log.InfoLevel:
		writerMap[log.InfoLevel] = writer
		fallthrough
	case log.WarnLevel:
		writerMap[log.WarnLevel] = writer
		fallthrough
	case log.ErrorLevel:
		writerMap[log.ErrorLevel] = writer
		fallthrough
	case log.FatalLevel:
		writerMap[log.FatalLevel] = writer
		fallthrough
	case log.PanicLevel:
		writerMap[log.PanicLevel] = writer
	default:
		log.WithField("level", fileLogLevelStr).Fatal("Invalid log level for file log")
	}
	log.WithField("level", fileLogLevelStr).Info("Set file log level")

	log.AddHook(lfshook.NewHook(
		writerMap,
		&log.JSONFormatter{},
	))
}
