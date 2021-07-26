package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/gangjun06/d4dj-info-server/conf"
	"github.com/gangjun06/d4dj-info-server/utils/crypto"
	"github.com/panjf2000/ants/v2"
)

var errFileNotFound = fmt.Errorf("err file not found")

type status struct {
	IsSuccess    bool
	FileName     string
	ErrorMessage string
}

func Start() {
	lastList, err := openListFile()
	if err != errFileNotFound && err != nil {
		// TODO: Save Error Log To database
		return
	}
	c := make(chan *status)
	go do("iOSResourceList.msgpack", c)
	if result := <-c; !result.IsSuccess {
		// TODO: Save Erro Log To database
		fmt.Println(result)
		return
	}

	curList, err := openListFile()
	if err != nil {
		// TODO: Save Error Log To database
		return
	}
	list := []string{}
	for k := range curList {
		if _, ok := lastList[k]; !ok || strings.HasPrefix(k, "Master") {
			list = append(list, k)
		}
	}

	p, _ := ants.NewPool(conf.Get().CrawlerPool)
	defer p.Release()
	go func() {
		for _, d := range list {
			p.Submit(func() { do(d, c) })
		}
	}()

	listLen := len(list)
	count := 1
	for range list {
		result := <-c
		fmt.Println(count, "/", listLen, result)
		count++
	}
}

func openListFile() (map[string]interface{}, error) {
	listFilePath := path.Join(conf.Get().AssetPath, "iOSResourceList.json")
	var lastFileList map[string]interface{}
	file, err := ioutil.ReadFile(listFilePath)
	if err != nil {
		return map[string]interface{}{}, errFileNotFound
	}
	if err := json.Unmarshal(file, &lastFileList); err != nil {
		return map[string]interface{}{}, err
	}
	return lastFileList, err
}

func do(file string, c chan<- *status) {
	data, err := downlaod(file)
	if err != nil {
		c <- &status{IsSuccess: false, FileName: file, ErrorMessage: err.Error()}
		return
	}
	decrypt, err := crypto.New().Decrypt(data)
	if err != nil {
		c <- &status{IsSuccess: false, FileName: file, ErrorMessage: err.Error()}
		return
	}
	savePath := path.Join(conf.Get().AssetPath, strings.ReplaceAll(file, ".enc", ""))
	if _, err := os.Stat(savePath); os.IsExist(err) && !strings.HasPrefix(file, "Master") {
		c <- &status{IsSuccess: true, FileName: file, ErrorMessage: "file is already exists"}
		return
	}
	if err := ioutil.WriteFile(savePath, decrypt, 0644); err != nil {
		dir, _ := filepath.Split(savePath)
		if err := os.MkdirAll(dir, 0766); err != nil {
			c <- &status{IsSuccess: false, FileName: file, ErrorMessage: "Error create directory: " + err.Error()}
		}
		if err := ioutil.WriteFile(savePath, decrypt, 0644); err != nil {
			c <- &status{IsSuccess: false, FileName: file, ErrorMessage: "Error save file: " + err.Error()}
		}
	}
	if err := msgpackToJSON(savePath); err != nil {
		c <- &status{IsSuccess: false, FileName: file, ErrorMessage: err.Error()}
		return
	}
	savePath = strings.ReplaceAll(savePath, "msgpack", "")
	os.Remove(savePath)
	c <- &status{IsSuccess: true, FileName: file, ErrorMessage: ""}
}

func downlaod(path string) ([]byte, error) {
	downloadPath := path
	if !strings.HasSuffix(path, "acb") {
		downloadPath += ".enc"
	}
	resp, err := http.Get(conf.Get().AssetServerPath + downloadPath)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return []byte{}, errors.New("error request")
	}

	return ioutil.ReadAll(resp.Body)
}

func msgpackToJSON(filePath string) error {
	c := exec.Command(conf.Get().ToolPath, filePath)
	return c.Run()
}
