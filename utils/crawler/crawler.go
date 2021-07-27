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

var ModifiedDate map[string]time.Time

func init() {
	ModifiedDate = make(map[string]time.Time)
}

func Start() {
	if modified, err := isModified("iOSResourceList.msgpack"); !modified || err != nil {
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}
	lastList, err := openListFile()
	if err != nil && err != errFileNotFound {
		fmt.Println(err.Error())
		return
	}
	c := make(chan *status)
	go do("iOSResourceList.msgpack", c)
	if result := <-c; !result.IsSuccess {
		fmt.Println(result)
		return
	}

	curList, err := openListFile()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	list := []string{}
	for k := range curList {
		if _, ok := lastList[k]; !ok || strings.HasPrefix(k, "Master") {
			list = append(list, k)
		}
	}

	p, _ := ants.NewPoolWithFunc(conf.Get().CrawlerPool, func(i interface{}) {
		data := i.(string)
		do(data, c)
	})
	defer p.Release()
	go func() {
		for _, d := range list {
			p.Invoke(d)
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

func getDownloadPath(path string) string {
	downloadPath := conf.Get().AssetServerPath + path
	if !strings.HasSuffix(path, "acb") {
		downloadPath += ".enc"
	}
	return downloadPath
}

func isModified(path string) (bool, error) {
	resp, err := http.Head(getDownloadPath(path))
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 200 {
		return false, errors.New("error request")
	}
	parsedTime, _ := time.Parse(time.RFC1123, resp.Header.Get("Last-Modified"))
	lastModified, ok := ModifiedDate[path]
	if !ok || lastModified.Before(parsedTime) {
		ModifiedDate[path] = parsedTime
		return true, nil
	}
	return false, nil
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
	if strings.HasPrefix(file, "Master") {
		if modified, _ := isModified(file); !modified {
			c <- &status{IsSuccess: true, FileName: file, ErrorMessage: "file is not modified"}
			return
		}
	}
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
	savePath := path.Join(conf.Get().AssetPath, file)

	// skip save if file exists
	// if _, err := os.Stat(savePath); os.IsExist(err) && !strings.HasPrefix(file, "Master") {
	// 	c <- &status{IsSuccess: true, FileName: file, ErrorMessage: "file is already exists"}
	// 	return
	// }

	// If forder not exists, create folder and save file
	if err := ioutil.WriteFile(savePath, decrypt, 0644); err != nil {
		dir, _ := filepath.Split(savePath)
		if err := os.MkdirAll(dir, 0766); err != nil {
			c <- &status{IsSuccess: false, FileName: file, ErrorMessage: "Error create directory: " + err.Error()}
		}
		if err := ioutil.WriteFile(savePath, decrypt, 0644); err != nil {
			c <- &status{IsSuccess: false, FileName: file, ErrorMessage: "Error save file: " + err.Error()}
		}
	}
	if strings.HasSuffix(savePath, "msgpack") || strings.HasPrefix(savePath, "chart_") {
		if err := msgpackToJSON(savePath); err != nil {
			c <- &status{IsSuccess: false, FileName: file, ErrorMessage: err.Error()}
			return
		}
		if err := os.Remove(savePath); err != nil {
			fmt.Println(err)
		}
	}

	c <- &status{IsSuccess: true, FileName: file, ErrorMessage: ""}
}

func downlaod(path string) ([]byte, error) {
	resp, err := http.Get(getDownloadPath(path))
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
	c := exec.Command("dotnet", conf.Get().ToolPath, filePath)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
