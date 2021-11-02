package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/gangjun06/d4dj-crawler/awsutil"
	"github.com/gangjun06/d4dj-crawler/conf"
	"github.com/gangjun06/d4dj-crawler/log"
	"github.com/gangjun06/d4dj-crawler/parser"
	"github.com/gangjun06/d4dj-crawler/parser/crypto"
	"github.com/panjf2000/ants/v2"
)

var errFileNotFound = fmt.Errorf("err file not found")

type status struct {
	IsSuccess    bool
	FileName     string
	ErrorMessage string
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
		// log.Println()

		str := fmt.Sprintf("%d / %d %s", count, listLen, result.FileName)
		if result.IsSuccess {
			if result.ErrorMessage != "" {
				log.Log.WithField("infoMsg", result.ErrorMessage).Info(str)
			}
			log.Log.Info(str)
		} else {
			log.Log.WithField("errorName", result.ErrorMessage).Warn(str)
		}
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
	lastModified, err := awsutil.ModifiedDate(path)
	if err != nil || lastModified == nil {
		return true, nil
	}
	parsedTime, _ := time.Parse(time.RFC1123, resp.Header.Get("Last-Modified"))
	// lastModified, ok := ModifiedDate[path]
	if lastModified.Before(parsedTime) {
		return true, nil
	}
	return false, nil
}

func openListFile() (map[string]interface{}, error) {
	var file []byte
	var err error
	awsData, err := awsutil.GetFile("iOSResourceList.json")
	var lastFileList map[string]interface{}
	if err == nil && awsData != nil {
		file = *awsData
	} else {
		listFilePath := path.Join(conf.Get().AssetPath, "iOSResourceList.json")
		file, err = ioutil.ReadFile(listFilePath)
		if err != nil {
			return map[string]interface{}{}, errFileNotFound
		}
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

	if err := parser.Parse(file, decrypt); err != nil {
		c <- &status{IsSuccess: false, FileName: file, ErrorMessage: err.Error()}
		return
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
