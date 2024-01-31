package commoninterface

type Logger interface {
	Info(msc string, args ...any)
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
}
