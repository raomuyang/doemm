package sync

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type MapToPath func(fileName string) (path string)

func PullFilesById(token string, gistId string, f MapToPath) error {
	gistInfo, err := GetGistInfo(token, gistId)
	if err != nil {
		return err
	}
	return PullFiles(token, gistInfo, f)
}

func PullFiles(token string, gistInfo *GistInfo, f MapToPath) error {
	files := gistInfo.Files

	wg := sync.WaitGroup{}

	errChannel := make(chan error, 100)
	for name, fileInfo := range files {
		path := f(name)
		if len(fileInfo.Content) == fileInfo.Size {
			err := writeFile(path, fileInfo.Content)
			if err != nil {
				return err
			}
		} else {
			wg.Add(1)
			go func() {
				defer wg.Done()
				content, e := GetRawFile(token, fileInfo.RawUrl)
				if e != nil {
					log.Warnf("goroutine: can not download raw file: %s", path)
					return
				}
				e = writeFile(path, *content)
				if e != nil {
					log.Warnf("goroutine: write file failed, cause: %v", e)
					select {
					case errChannel <- e:
						log.Debug("put error")
					default:
						break
					}
				}
			}()
		}
	}

	wg.Wait()
	failed := len(errChannel)
	if failed > 0 {
		return errors.New(fmt.Sprintf("sync (download raw files) failed (%d files download failed)", failed))
	}
	return nil
}

func PullSingleFile(token string, targetFile string, possibleSuffix string, gistInfo *GistInfo, f MapToPath) error {
	files := gistInfo.Files

	targetFileWithSuffix := targetFile + possibleSuffix
	for name, fileInfo := range files {
		if strings.Compare(targetFile, name) != 0 &&
			strings.Compare(targetFileWithSuffix, name) != 0 {
			continue
		}

		path := f(name)
		if len(fileInfo.Content) == fileInfo.Size {
			err := writeFile(path, fileInfo.Content)
			if err != nil {
				return err
			}
		} else {
			content, e := GetRawFile(token, fileInfo.RawUrl)
			if e != nil {
				return e
			}
			e = writeFile(path, *content)
			if e != nil {
				return e
			}
		}
		return nil
	}

	return errors.Errorf("not found such item: %s", targetFile)
}

func PushLocalFiles(token string, pathList []string, gistId string) (gistInfo *GistInfo, err error) {
	gist := Gist{Files: map[string]*FileContent{}}
	for _, path := range pathList {
		content, e := readFile(path)
		if e != nil {
			err = e
			return
		}
		name := filepath.Base(path)
		gist.Files[name] = content
	}
	gistInfo, err = UpdateGistInfo(token, &gist, gistId)
	if err != nil {
		log.Warnf("This error maybe cause set a nil content to a not exists file: %v", err)
	}
	return
}

func PushSingleFile(token string, path string, gistId string) (gistInfo *GistInfo, err error) {
	gist := Gist{Files: map[string]*FileContent{}}
	content, err := readFile(path)
	if err != nil {
		return
	}
	name := filepath.Base(path)
	gist.Files[name] = content
	gistInfo, err = UpdateGistInfo(token, &gist, gistId)
	return
}

// 文件不存在时返回一个空的内容
func readFile(path string) (content *FileContent, err error) {
	stat, e := os.Stat(path)

	// 文件为相对路径时，可能会读取到同名的目录，不取任何内容，将content置为空
	if e != nil || stat.IsDir() {
		log.Infof("delete file content: %s", path)
		return
	}

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	content = &FileContent{Content: string(data)}
	return
}

func writeFile(path, content string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Warnf("Failed to create file: %s", path)
	}
	defer file.Close()
	n, err := file.WriteString(content)
	log.Infof("Update file %s, size %d", path, n)
	return err
}
