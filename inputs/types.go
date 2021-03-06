package inputs

type InputType int

const (
	STORE InputType = iota
	CONFIG
	PULL
	PUSH
	SWITCH
	LIST
	PRINT
	DELETE
)

type Input interface {
	// input items or command(s)
	GetItems() []string
	GetInputType() InputType
	// description of input or alias of command(s)
	GetSummary() string
}
