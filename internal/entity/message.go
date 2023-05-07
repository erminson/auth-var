package entity

type SourceCode string

const (
	Call SourceCode = "call"
	Sms  SourceCode = "sms"
)

type Message struct {
	Text   string
	Digit  byte
	Source SourceCode
}
