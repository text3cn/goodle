package goodlog

func Trace(out ...interface{}) {
	instance()
	goodlogSvc.Trace(out...)
}

func Debug(out ...interface{}) {
	instance()
	goodlogSvc.Debug(out...)
}

func Info(out ...interface{}) {
	instance()
	goodlogSvc.Info(out...)
}

func Warn(out ...interface{}) {
	instance()
	goodlogSvc.Warn(out...)
}

func Error(out ...interface{}) {
	instance()
	goodlogSvc.Error(out...)
}

func Fatal(out ...interface{}) {
	instance()
	goodlogSvc.Fatal(out...)
}

func Tracef(out ...interface{}) {
	instance()
	goodlogSvc.Tracef(out...)
}

func Debugf(out ...interface{}) {
	instance()
	goodlogSvc.Debugf(out...)
}

func Infof(out ...interface{}) {
	instance()
	goodlogSvc.Infof(out...)
}

func Warnf(out ...interface{}) {
	instance()
	goodlogSvc.Warnf(out...)
}

func Errorf(out ...interface{}) {
	instance()
	goodlogSvc.Errorf(out...)
}

func Fatalf(out ...interface{}) {
	instance()
	goodlogSvc.Fatalf(out...)
}
