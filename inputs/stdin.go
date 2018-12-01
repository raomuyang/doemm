package inputs

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	Prefix      = "emm...> "
	Done        = "emm.done"
	AliasPrefix = "emm.alias."
)

type Stdin struct {
	CmdLines []string
	Alias    string
}

func (stdin *Stdin) GetItems() []string {
	return stdin.CmdLines
}

func (stdin *Stdin) GetInputType() InputType {
	return STORE
}

func (stdin *Stdin) GetSummary() string {
	return stdin.Alias
}

func GetStdin() (*Stdin, error) {
	input := Stdin{}

	err := input.read()
	if err != nil {
		return nil, err
	}

	return &input, nil
}

func (stdin *Stdin) read() error {
	fmt.Println(
		"Input multipart commands in interactive \n" +
			"1. type emm.done to exit \n" +
			"2. type emm.alias.ALIAS_NAME to set alias of commands")
	alias := ""
	for {
		fmt.Print(Prefix)
		line, err := readLineFromStdin()
		if err != nil {
			return err
		}
		if strings.HasPrefix(*line, AliasPrefix) {
			alias = (*line)[len(AliasPrefix):]
			continue
		} else if strings.Compare(*line, Done) == 0 {
			if len(alias) == 0 {
				fmt.Print(Prefix)
				fmt.Println("please input alias via emm.alias.ALIAS_NAME")
				continue
			}
			break
		}
		stdin.CmdLines = append(stdin.CmdLines, *line)
	}
	stdin.Alias = alias
	return nil
}

func readLineFromStdin() (*string, error) {

	inputReader := bufio.NewReader(os.Stdin)
	var segments []string
	for {
		input, isPrefix, err := inputReader.ReadLine()
		if err != nil {
			return nil, err
		}
		segments = append(segments, string(input))
		if !isPrefix {
			break
		}
	}

	var line string
	if len(segments) == 1 {
		line = segments[0]
	} else {
		line = strings.Join(segments, "")
	}

	return &line, nil
}
