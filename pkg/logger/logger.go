package logger

type Interface interface {
	Printf(format string, v ...interface{})
}

type LogLineType string

const (
	LogLineTypeAxe       LogLineType = "axe"
	LogLineTypeContainer LogLineType = "container"
)

type LogLine struct {
	Type      LogLineType
	Namespace string
	Name      string
	Bytes     []byte
}
