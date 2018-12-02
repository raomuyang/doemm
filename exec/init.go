package exec

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/user"
	"path/filepath"
)

const (
	logName                 = "emm..its.log"
	defaultLogLevel         = log.Level(3)
	defaultGithubOAuthToken = "<set oauth token here>"
)

var (
	userHome      string
	appHome       string
	configPath    string
	bucketDir     string
	logLevel      log.Level
	configuration Configuration
)

func init() {
	log.SetLevel(defaultLogLevel)

	current, err := user.Current()
	if err != nil {
		panic(err)
	}

	userHome = current.HomeDir
	appHome = userHome + "/.do_emm"
	configPath = appHome + "/configuration.yml"
	bucketDir = appHome + "/bucket"

	_, err = os.Stat(bucketDir)
	if err != nil {
		os.MkdirAll(bucketDir, 0755)
	}

	_, err = os.Stat(configPath)

	createNew := false
	if err != nil {
		createNew = true
	}

	initConfigYml(createNew)

	initLogger()
}

func initConfigYml(createNew bool) {
	configuration = getDefaultConfig()
	if createNew {
		err := dumpConfiguration()
		if err != nil {
			exit("Error: create default yaml file failed, cause: %s", err)
		}
		log.Info("Initial configuration file by default")
	} else {
		err := loadConfiguration()
		if err != nil {
			exit("Error: read configuration failed: %v", err)
		}
	}
}

func getDefaultConfig() Configuration {
	return Configuration{
		LogLevel:       3,
		DefaultEncrypt: false,
		GistToken:      defaultGithubOAuthToken}
}

func initLogger() {
	logLevel = log.Level(configuration.LogLevel)
	log.SetLevel(logLevel)

	filePath := filepath.Join(appHome, logName)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0766)
	if err != nil {
		exit("Error: create log file failed, cause: %s", err)
	}
	log.SetOutput(file)
}
