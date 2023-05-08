package entity

type SourceCode string

const (
	Call SourceCode = "call"
	Sms  SourceCode = "sms"
)

func (s SourceCode) String() string {
	return string(s)
}

type Message struct {
	Text   string
	Digit  int
	Source SourceCode
}
