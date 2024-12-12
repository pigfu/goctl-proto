package proto

const (
	Version2 = "proto2"
	Version3 = "proto3"
)

// File todo: group services
type File struct {
	Syntax   string
	Package  string
	Options  []*Option
	Messages []*Message
	*Service
}

type Option struct {
	Name  string
	Value string
}

type Message struct {
	Name   string
	Descs  []string
	Fields []*MessageField
}

type MessageField struct {
	Name          string
	Descs         []string
	TypeName      string
	Repeated      bool
	BaseTypeNames []string
}

type Service struct {
	Name  string
	Descs []string
	Rpcs  []*ServiceRpc
}

type ServiceRpc struct {
	Name     string
	Descs    []string
	Request  *Message
	Response *Message
}
