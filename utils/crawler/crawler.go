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
	"time"

	"github.com/gangjun06/d4dj-info-server/env"
	"github.com/gangjun06/d4dj-info-server/utils/crypto"
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
		// TODO: Save Error Log To database
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

	go func() {
		for _, d := range list {
			go do(d, c)
			time.Sleep(time.Second)
		}
	}()

	for range list {
		result := <-c
		fmt.Println(result)
	}
}

func openListFile() (map[string]interface{}, error) {
	listFilePath := path.Join(string(env.Get(env.KeyAssetPath)), "iOSResourceList.json")
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
	savePath := path.Join(env.Get(env.KeyAssetPath), strings.ReplaceAll(file, ".enc", ""))
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
	c <- &status{IsSuccess: true, FileName: file, ErrorMessage: ""}
}

func downlaod(path string) ([]byte, error) {
	downloadPath := path
	if !strings.HasSuffix(path, "acb") {
		downloadPath += ".enc"
	}
	resp, err := http.Get("https://resources.d4dj-groovy-mix.com/1161b98bd529f32da32e631f1504b928c4f3961f/" + downloadPath)
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
	c := exec.Command(env.Get(env.KeyEnvToolPath), filePath)
	return c.Run()
}
