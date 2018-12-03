package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// create: POST  /gists
// edit:   PATCH /gists/:gist_id
// get:    GET   /gists/:gist_id

const (
	endpoint string = "https://api.github.com"
	editGist string = "/gists/%s"
	getGist  string = "/gists/%s"
	gists    string = "/gists"

	tokenFormat         string = "token %s"
	tokenHeader         string = "Authorization"
	defaultIteratePages int    = 100
)

var client = &http.Client{}

// request body of file content
type FileContent struct {
	Content string `json:"content"`
}

// request body of create/edit gist
type Gist struct {
	// the description of gist
	Description string `json:"description"`
	// set the new gist to public/secret , can not apply in an exists gist
	Public bool `json:"public"`
	// files: {"file name": "file content"}
	Files map[string]*FileContent `json:"files"`
}

// response info: include the file info which include in the gist
type FileInfo struct {
	FileName string `json:"filename"`
	Type     string `json:"type"`
	Language string `json:"language"`
	RawUrl   string `json:"raw_url"`
	// check content: size == len(Content)
	Size      int    `json:"size"`
	Truncated bool   `json:"truncated"`
	Content   string `json:"content"`
}

// The major attributes
type GistInfo struct {
	URL         string               `json:"url"`
	Id          string               `json:"id"`
	ForksUrl    string               `json:"forks_url"`
	Files       map[string]*FileInfo `json:"files"`
	NodeId      string               `json:"node_id"`
	CommitsURL  string               `json:"commits_url"`
	GitPullUrl  string               `json:"git_pull_url"`
	GitPushUrl  string               `json:"git_push_url"`
	HtmlUrl     string               `json:"html_url"`
	Public      bool                 `json:"public"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
	Description string               `json:"description"`
	Truncated   bool                 `json:"truncated"`
}

// Get gist info by gist id
func GetGistInfo(token, gistId string) (gistInfo *GistInfo, err error) {
	log.Infof("Get gist by id: %s", gistId)
	url := endpoint + fmt.Sprintf(getGist, gistId)
	method := "GET"
	option := fmt.Sprintf("get gist info(%s)", gistId)
	object := &GistInfo{}

	_, err = doRequest(token, option, method, url, nil, object)

	if err == nil {
		gistInfo = object
	}

	return
}

// List all gists page by page, first page: ""
func ListAllGists(token, page string) (gistList *[]GistInfo, nextPage string, err error) {
	log.Infof("List all gists (page: %s)", page)
	url := page
	if len(page) == 0 {
		url = endpoint + gists
	}

	method := "GET"
	option := "list gists"
	object := &[]GistInfo{}

	resp, err := doRequest(token, option, method, url, nil, object)

	if err != nil {
		return
	}

	h := strings.Split(resp.Header.Get("LINK"), ";")[0]
	nextPage = h[1 : len(h)-1]
	gistList = object

	return
}

func FindGistByDesc(token, desc string) (*GistInfo, error) {
	log.Infof("Find (one) gist by desc: %s", desc)

	next := ""
	var i int
	for i = 0; i < defaultIteratePages; i++ {
		gistList, nextPage, err := ListAllGists(token, next)
		if err != nil {
			return nil, err
		}

		for _, gist := range *gistList {
			if strings.Compare(strings.Trim(gist.Description, " \n"), desc) == 0 {
				return &gist, nil
			}
		}

		if !strings.Contains(nextPage, "page=1") {
			next = nextPage
		} else {
			break
		}

	}
	log.Warnf("Try find gist information with %d times, found: false", i)

	return nil, nil
}

func CreateGist(gist *Gist, token string) (gistInfo *GistInfo, err error) {
	url := endpoint + gists
	body, err := json.Marshal(gist)
	if err != nil {
		log.Warnf("Marshal request body failed: %v", err)
		return
	}

	method := "POST"
	option := "create a new gist"
	object := &GistInfo{}

	_, err = doRequest(token, option, method, url, bytes.NewReader(body), object)

	if err == nil {
		gistInfo = object
	}

	return
}

func UpdateGistInfo(token string, gist *Gist, gistId string) (responseGist *GistInfo, err error) {
	log.Infof("Update gist info: %s", gistId)
	url := endpoint + fmt.Sprintf(editGist, gistId)

	data, err := json.Marshal(gist)
	body := bytes.NewReader(data)

	if err != nil {
		log.Warnf("Marshal request body failed: %v", err)
		return
	}
	method := "PATCH"
	object := &GistInfo{}
	option := fmt.Sprintf("update the information of gist(%s)", gistId)
	_, err = doRequest(token, option, method, url, body, object)
	if err == nil {
		responseGist = object
	}

	return
}

func GetRawFile(token, url string) (content *string, err error) {
	method := "GET"
	option := "get raw file"

	resp, err := doRequest(token, option, method, url, nil, nil)
	if err == nil {
		data, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			log.Warnf("Socket error: %v", e)
			err = e
			return
		}
		str := string(data)
		content = &str
	}
	return
}

func addPublicHeader(request *http.Request, token string) {
	request.Header.Add(tokenHeader, fmt.Sprintf(tokenFormat, token))
}

func checkResponse(response *http.Response) error {
	if response.StatusCode/100 != 2 {
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Warnf("Check response failed, cause: %v", err)
			return errors.New("Unknown")
		}
		return errors.New(string(data))
	}
	return nil
}

func doRequest(token string, option string, method string, url string,
	body io.Reader, object interface{}) (resp *http.Response, err error) {

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Warnf("Failed to create request: %v", err)
		return
	}
	addPublicHeader(request, token)

	resp, err = client.Do(request)
	if err != nil {
		log.Warnf("Failed to request %s, cause: %v", option, err)
	} else if err = checkResponse(resp); err != nil {
		return
	} else if object == nil {
		// process response body by self
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warnf("Socket error: %v", err)
		return
	}
	err = json.Unmarshal(data, object)
	return
}
