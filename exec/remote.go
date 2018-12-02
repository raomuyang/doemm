package exec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/raomuyang/doemm/sync"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

const (
	defaultGistDescription = "This gist was created by do.emm, " +
		"it is strongly recommended that do not edit this description " +
		"and keep it unique."
	setGistPrefix = "emm.gist_id."
	initGistCmd   = "emm.init_gist"

	helloWorldFile    = "doemm.hello"
	helloWorldContent = "echo this file was created by do.emm...\n" +
		"echo created time: %v\n" +
		"echo from %s\n"
)

func mapNameToLocal(name string) string {
	fmt.Printf("new: %s\n", name)
	local := path.Join(bucketDir, name)
	old := getConflictName(local)
	e := os.Remove(old)
	log.Debugf("try to remove local old file: %v", e)
	return local
}

func getConflictName(name string) string {
	var conflictName string
	if strings.HasSuffix(name, encryptSuffix) {
		conflictName = name[:len(name)-len(encryptSuffix)]
	} else {
		conflictName = name + encryptSuffix
	}
	return conflictName
}

func initGist() (*sync.GistInfo, error) {

	content := fmt.Sprintf(helloWorldContent, time.Now(), runtime.GOOS)
	fileContent := sync.FileContent{Content: content}
	gist := sync.Gist{
		Description: defaultGistDescription,
		Public:      false,
		Files:       map[string]*sync.FileContent{helloWorldFile: &fileContent}}
	return sync.CreateGist(&gist, configuration.GistToken)
}

func dumpGist(gistInfo *sync.GistInfo) {
	data, _ := json.Marshal(gistInfo)
	e := ioutil.WriteFile(path.Join(appHome, ".gist"), data, 0644)
	log.Debugf("Dump gist info, err: %v", e)
}

func checkAndGetGistId() (gistId string, err error) {
	gistId = configuration.GistId

	if len(gistId) == 0 {
		fmt.Print("please check gist existing and set gist-id or init gist,\n" +
			"> type \"" + setGistPrefix + "\" to set gist id\n" +
			"> or type \"" + initGistCmd + "\" to auto-init\n" +
			"> ")
		cmd := ""
		fmt.Scanln(&cmd)

		if strings.Compare(cmd, initGistCmd) == 0 {
			fmt.Println("wait a minute...")
			var gistInfo *sync.GistInfo
			var e error
			gistInfo, e = sync.FindGistByDesc(configuration.GistToken, defaultGistDescription)
			if e != nil {
				err = e
				return
			} else if gistInfo == nil {
				fmt.Println("create a new gist...")
				gistInfo, e = initGist()
				if e != nil {
					err = e
					return
				}
				gistId = gistInfo.Id
				fmt.Println("gist created!!!")
				fmt.Printf("gist url: %s\n", gistInfo.URL)
			} else {
				fmt.Printf("found existing gist: %s\n", gistInfo.URL)
				gistId = gistInfo.Id
			}

			dumpGist(gistInfo)

		} else if strings.HasPrefix(cmd, setGistPrefix) {
			gistId = cmd[len(setGistPrefix):]
		} else {
			err = errors.New("unknown input: " + cmd)
		}
		configuration.GistId = gistId
		dumpConfiguration()
		fmt.Printf("apply gist-id: %s \n", gistId)
	}

	return

}

func pull(fileName string) (err error) {

	if len(configuration.GistToken) == 0 {
		return errors.New("please config gist token via `doemm config -gist <gist token>`")
	}
	gistId, err := checkAndGetGistId()
	if err != nil {
		return
	}
	if len(fileName) == 0 {
		err = sync.PullFilesById(configuration.GistToken, gistId, mapNameToLocal)
	} else {
		var gistInfo *sync.GistInfo
		gistInfo, err = sync.GetGistInfo(configuration.GistToken, gistId)
		if err != nil {
			return
		}
		err = sync.PullSingleFile(configuration.GistToken, fileName, encryptSuffix, gistInfo, mapNameToLocal)
	}

	return
}

func push(itemName string) (err error) {

	if len(configuration.GistToken) == 0 {
		return errors.New("please config gist token via `doemm config -gist <gist token>`")
	}
	gistId, err := checkAndGetGistId()
	if err != nil {
		return
	}
	if len(itemName) == 0 {
		var fInfoList []os.FileInfo
		fInfoList, err = ioutil.ReadDir(bucketDir)
		if err != nil {
			return
		}

		var gistInfo *sync.GistInfo

		gistInfo, err = sync.GetGistInfo(configuration.GistToken, gistId)
		if err != nil {
			return
		}

		var pathList []string
		for _, info := range fInfoList {

			local := path.Join(bucketDir, info.Name())

			// append local
			var f *sync.FileInfo
			f = gistInfo.Files[info.Name()]
			if f == nil || len(f.Content) != f.Size {
				// not exists
				fmt.Printf("new:    %s\n", info.Name())
				pathList = append(pathList, local)
			} else {
				// not equal
				data, _ := ioutil.ReadFile(local)
				equal := bytes.Compare(data, []byte(f.Content)) == 0
				if !equal {
					fmt.Printf("new: %s\n", info.Name())
					pathList = append(pathList, local)
				}
				log.Infof("check equal with local: %v", equal)
			}

			// command-alias.e 和 command-alias 不能同时保留，必须删除其中一个
			// 程序传入一个不存在的相对路径，在同步时若content为空，gist会删除该文件
			uniqueConflict := getConflictName(info.Name())
			f = gistInfo.Files[uniqueConflict]
			if f != nil {
				fmt.Printf("delete: %s\n", uniqueConflict)
				pathList = append(pathList, uniqueConflict)
				log.Debugf("local: %s, unique check: %s", info.Name(), uniqueConflict)
			}
		}

		gistInfo, err = sync.PushLocalFiles(configuration.GistToken, pathList, gistId)
		if err != nil {
			return
		}
		dumpGist(gistInfo)
	} else {
		var gistInfo *sync.GistInfo
		filePath := path.Join(bucketDir, itemName)
		_, e := os.Stat(filePath)
		if e != nil {
			filePath += encryptSuffix
			_, e := os.Stat(filePath)
			if e != nil {
				// not found
				return e
			}
		}

		var g *sync.GistInfo
		g, err = sync.PushSingleFile(configuration.GistToken, filePath, gistId)
		if err != nil {
			return
		}
		gistInfo = g

		duplicatePath := getConflictName(filePath)
		g, err = sync.PushSingleFile(configuration.GistToken, duplicatePath, gistId)
		log.Debugf("Try to delete exists file: %s, result: %v", duplicatePath, err)
		if err == nil {
			gistInfo = g
		}
		dumpGist(gistInfo)
	}

	return
}
