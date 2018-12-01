package inputs

import "strings"

type ArgsInput struct {
	cmdLines []string
	Type     InputType
	Alias    string
}

func (input *ArgsInput) GetItems() []string {
	return input.cmdLines
}

func (input *ArgsInput) GetInputType() InputType {
	return input.Type
}

func (input *ArgsInput) GetSummary() string {
	return input.Alias
}

func GetArgsInput(args []string, alias string) (*ArgsInput, error) {
	i := 0
	for i = range args {
		if strings.HasSuffix(args[i], alias) {
			break
		}
	}
	line := strings.Join(args[i+1:], " ")
	input := ArgsInput{Type: STORE, cmdLines: []string{line}, Alias: alias}
	return &input, nil
}
