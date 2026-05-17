package log

type noopLogger struct{}

func (n *noopLogger) Debug(msg string, fields ...Field)            {}
func (n *noopLogger) Info(msg string, fields ...Field)             {}
func (n *noopLogger) Warn(msg string, fields ...Field)             {}
func (n *noopLogger) Error(msg string, err error, fields ...Field) {}
func (n *noopLogger) Fatal(msg string, fields ...Field)            {}
func (n *noopLogger) With(fields ...Field) Logger                  { return n }

func NewNoop() Logger { return &noopLogger{} }
