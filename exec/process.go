package exec

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/raomuyang/doemm/inputs"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

var SK = "Bj.Hd.MDYHyDl501"

type Configuration struct {
	Gist           string `yaml:"gist"`
	DefaultEncrypt bool   `yaml:"default_encrypt"`

	// DEBUG(5) INFO(4) WARN(3) ERROR(2) FATAL(1) PANIC(0)
	LogLevel uint32 `yaml:"log_level"`
}

func ProcessInput(input inputs.Input) {
	switch input.GetInputType() {
	case inputs.CONFIG:
		processConfig(input.GetItems())
	case inputs.STORE:
		save(input.GetItems(), input.GetSummary())
		fmt.Println("item stored!")
	case inputs.LIST:
		fmt.Printf("////// Stored list //////\n\n")
		res, err := listAll()
		if err != nil {
			exit("Failed to list, cause: %v", err)
		}
		for i, item := range res {
			fmt.Printf("%d. %s", i, item)
		}
	case inputs.PRINT:
		items, err := show(input.GetSummary())
		if err != nil {
			exit("Failed to load item: %v", err)
		}
		if len(items) == 0 {
			fmt.Println("Not found!")
		} else {
			for _, v := range items {
				fmt.Println(v)
			}
		}
	case inputs.SWITCH:
		fmt.Printf("------ switch alias %s ------ \n\n", input.GetSummary())
		items, err := show(input.GetSummary())
		if err != nil {
			exit("Failed to load item: %v", err)
		}

		if len(items) == 0 {
			fmt.Println("Not found!")
		} else {
			err := executeCommands(items)
			if err != nil {
				fmt.Println("\n ------ done ------")
			}
		}

	case inputs.SYNC:

	}
}

// 0. gist_token, 1. encrypt (true|false)
func processConfig(items []string) {
	encrypt := false
	gistToken := items[0]
	if items[1] == "true" {
		encrypt = true
	}

	configuration.Gist = gistToken
	configuration.DefaultEncrypt = encrypt
	dumpConfiguration()
	fmt.Printf("emm... encrypt: %v \ndone!\n", encrypt)
}

func exit(msg string, err error) {
	fmt.Printf(msg, err)
	log.Warnf(msg, err)
	os.Exit(2)
}

func dumpConfiguration() error {
	out, err := yaml.Marshal(configuration)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(configPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(out)
	return err
}

func loadConfiguration() error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	in, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(in, &configuration)
	return err
}

func decryptText(base64Str, key string) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			exit(fmt.Sprintf("Error: decrypt failed: %v", err), nil)
		}
	}()

	cipherText, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return "", err
	}

	fillKey := make([]byte, 32)
	copy(fillKey, []byte(key))

	block, err := aes.NewCipher([]byte(fillKey))
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", errors.New("cipher text too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	cipher.NewCFBDecrypter(block, iv).XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

func encryptText(srcText, key string) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			exit(fmt.Sprintf("Error: decrypt failed: %v", err), nil)
		}
	}()

	fillKey := make([]byte, 32)
	copy(fillKey, []byte(key))

	block, err := aes.NewCipher([]byte(fillKey))
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(srcText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cipher.NewCFBEncrypter(block, iv).XORKeyStream(cipherText[aes.BlockSize:],
		[]byte(srcText))
	return base64.StdEncoding.EncodeToString(cipherText), nil
}
