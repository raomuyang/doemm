package exec

import (
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

func run(alias string) error {
	items, err := show(alias)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return errors.Errorf("not found: %s", alias)
	} else if len(items) == 1 {
		return executeCommands(items)
	} else {
		// 多行命令写入临时文件中运行
		switch runtime.GOOS {
		case "windows":
			return runBat(alias, items)
		default:
			return runShell(alias, items)
		}
	}
}

func runShell(alias string, items []string) error {
	tmpPath := path.Join(bucketDir, alias+".emm.tmp"+uuid.NewV1().String())

	content := "#!/usr/bin/env bash\n" + strings.Join(items, "\n")

	defer func() {
		e := os.Remove(tmpPath)
		if e != nil {
			log.Warnf("delete tmp script failed: %v", e)
		}
	}()

	err := ioutil.WriteFile(tmpPath, []byte(content), 0777)
	if err != nil {
		log.Warnf("write file %s failed: %v", tmpPath, err)
		return err
	}

	cmd := getCommand([]string{tmpPath})
	return cmd.Run()
}

func runBat(alias string, items []string) error {
	tmpPath := path.Join(bucketDir, alias+".emm.tmp."+uuid.NewV1().String()+".bat")

	content := strings.Join(items, "\r\n")

	defer func() {
		e := os.Remove(tmpPath)
		if e != nil {
			log.Warnf("delete tmp script failed: %v", e)
		}
	}()

	err := ioutil.WriteFile(tmpPath, []byte(content), 0777)
	if err != nil {
		log.Warnf("write file %s failed: %v", tmpPath, err)
		return err
	}

	cmd := getCommand([]string{tmpPath})
	return cmd.Start()
}

// 多条命令可能涉及到上下文的逻辑，所以此函数暂时只做一条命令的执行操作
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
