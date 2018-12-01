package exec

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

func executeCommands(lines []string) error {
	for _, content := range lines {
		content := strings.Trim(content, " ")
		args := strings.Split(content, " ")

		cmd := getCommand(args)

		err := cmd.Run()
		if err != nil {
			log.Warnf("execute command failure: %s", args[0])
			return err
		}
	}
	return nil
}

func getCommand(args []string) *exec.Cmd {
	var cmd *exec.Cmd
	if len(args) == 1 {
		cmd = exec.Command(args[0])
	} else {
		cmd = exec.Command(args[0], args[1:]...)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd
}
