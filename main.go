package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/raomuyang/doemm/exec"
	"github.com/raomuyang/doemm/inputs"
	"os"
)

const (
	CONFIG = "config"
	PULL   = "pull"
	PUSH   = "push"
)

var (
	help       = flag.Bool("help", false, "print help")
	switchTo   = flag.String("s", "", "switch to the specified commands-alias")
	printItems = flag.String("print", "", "print the specified commands-alias")
	alias      = flag.String("alias", "", "set the alias of store command(s)")
	list       = flag.Bool("list", false, "list the saved commands")

	configCommand  = flag.NewFlagSet("config", flag.ExitOnError)
	gistToken      = configCommand.String("gist", "", "gist token")
	defaultEncrypt = configCommand.Bool("encrypt", false, "default encrypt your command items")

	pullCommand    = flag.NewFlagSet("pull", flag.ExitOnError)
	pullSingleFile = pullCommand.String("single", "", "select a single item to pull")

	pushCommand    = flag.NewFlagSet("push", flag.ExitOnError)
	pushSingleFile = pushCommand.String("single", "", "select a single item to push")
)

func main() {
	input, err := getInput()
	if err != nil {
		printDefaultAndExit(err)
	}

	exec.ProcessInput(input)
}

func getInput() (input inputs.Input, err error) {

	if len(os.Args) <= 1 {
		// 从stdin中获取命令行信息
		return inputs.GetStdin()
	}

	switch os.Args[1] {
	case CONFIG:
		fmt.Println("update local configuration")
		err := configCommand.Parse(os.Args[2:])
		if err != nil {
			printDefaultAndExit(err)
		}
		c := inputs.ConfigParams{GistToken: *gistToken, Encrypt: *defaultEncrypt}
		input = &c
	case PULL:
		fmt.Println("pull item(s) from gist!")
		err := pullCommand.Parse(os.Args[2:])
		if err != nil {
			printDefaultAndExit(err)
		}
		s := inputs.PullParams{SingleItem: *pullSingleFile}
		input = &s
	case PUSH:
		fmt.Println("push item(s) to gist!")
		err := pushCommand.Parse(os.Args[2:])
		if err != nil {
			printDefaultAndExit(err)
		}
		s := inputs.PushParams{SingleItem: *pushSingleFile}
		input = &s
	default:
		flag.Parse()
		if *help {
			printDefault()
			os.Exit(0)
		} else if *printItems != "" {
			// switch the specified command(s)
			if len(*alias) != 0 {
				return nil, errors.New("-alias with -print is forbid")
			}
			input = &inputs.PrintParam{Target: *printItems}
			return
		} else if *switchTo != "" {
			// switch the specified command(s)
			if len(*alias) != 0 {
				return nil, errors.New("-alias with -s is forbid")
			}
			input = &inputs.SwitchParam{Target: *switchTo}
			return
		} else if *list {
			input = &inputs.ListParam{}
			return
		} else {
			// store commands
			if len(*alias) == 0 {
				err = errors.New("please specify the alias of command")
				return
			}

			input, err = inputs.GetArgsInput(os.Args[1:], *alias)
			return
		}

		printDefaultAndExit(errors.New("Unknown command: " + os.Args[1]))
	}

	return

}

func printDefaultAndExit(err error) {
	println(fmt.Sprintf("Error: %s.\n", err.Error()))
	printDefault()
	os.Exit(2)
}

func printDefault() {

	print("Usage:\n\n" +
		"save the command:\n" +
		"> doeem -alias <command alias> <target command to save>\n" +
		"switch a specified command:\n" +
		"> doemm -s <commands-alias>\n\n")
	flag.PrintDefaults()

	print("config: \n  write new config item into `config.yml`\n\n")
	configCommand.PrintDefaults()
	print("pull: \n  pull item(s) from gist.github.com\n  eg. pull [item-alias]\n\n")
	pullCommand.PrintDefaults()
	print("push: \n  push local item(s) to gist.github.com\n  eg. push [item-alias]\n\n")
	pushCommand.PrintDefaults()
}
