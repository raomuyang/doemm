package exec

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const encryptSuffix = ".e"

func save(commands []string, alias string) error {
	filePath := path.Join(bucketDir, alias)
	if configuration.DefaultEncrypt {
		e := os.Remove(filePath) // remove old
		log.Debugf("Remove exists public item %s, result: %v", filePath, e)
		filePath += encryptSuffix
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	text := strings.Join(commands, "\n")
	if configuration.DefaultEncrypt {
		text, err = encryptText(text, SK)
		if err != nil {
			log.Warnf("Failed to encrypt text, cause: %v", err)
			return err
		}
	}
	_, err = file.Write([]byte(text))
	if err != nil {
		log.Warnf("Failed to write content to %s, cause: %v", filePath, err)
	}
	return err
}

// 读取保存的命令行内容，默认先找未加密保存的内容，其次再尝试读取加密保存的内容
func show(alias string) (res []string, err error) {

	filePath := path.Join(bucketDir, alias)
	stat, err := os.Stat(filePath)
	if err != nil {
		filePath += encryptSuffix
		stat, err = os.Stat(filePath)
		if err != nil {
			// not found
			log.Warn("alias not found: %s, %v", alias, err)
			return
		}
	}

	log.Infof("State file: %s: %v", filePath, stat)

	in, err := ioutil.ReadFile(filePath)
	if err != nil {
		return res, err
	}
	content := string(in)
	if strings.HasSuffix(filePath, encryptSuffix) {
		content, err = decryptText(string(content), SK)
		if err != nil {
			log.Warn("Failed to decrypt content, cause: %v", err)
			return res, err
		}
	}

	res = strings.Split(string(content), "\n")
	return
}

func listAll() ([]string, error) {

	fileInfoList, err := ioutil.ReadDir(bucketDir)
	if err != nil {
		return []string{}, err
	}

	var res []string
	for _, v := range fileInfoList {
		name := v.Name()
		time := v.ModTime()
		if strings.HasSuffix(name, encryptSuffix) {
			name = name[:len(name)-len(encryptSuffix)]
		}

		item := fmt.Sprintf("%s  -  %v", name, time)
		res = append(res, item)
	}
	return res, nil
}
