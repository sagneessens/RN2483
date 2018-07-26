package rn2483

type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

type NOOPLogger struct{}

func (NOOPLogger) Println(v ...interface{})               {}
func (NOOPLogger) Printf(format string, v ...interface{}) {}

var (
	ERROR Logger = NOOPLogger{}
	WARN  Logger = NOOPLogger{}
	DEBUG Logger = NOOPLogger{}
)
