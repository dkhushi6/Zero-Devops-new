package logger

import "go.uber.org/zap"


func New(env string) *zap.Logger {
	if env == "production"{
		l,_ := zap.NewProduction()
		return l
	}
	l,_ := zap.NewDevelopment()
	return l
}