package log

import (
	"go.uber.org/zap/zapcore"
)

func stringifyErrors(fields []zapcore.Field) {
	for i := range fields {
		if fields[i].Type == zapcore.ErrorType {
			fields[i].Type = zapcore.StringType
			fields[i].String = fields[i].Interface.(error).Error()
			fields[i].Interface = nil
		}
	}
}
