package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/gangjun06/d4dj-crawler/awsutil"
	"github.com/gangjun06/d4dj-crawler/conf"
)

func Parse(fileName string, data []byte) error {
	savePath := path.Join(conf.Get().AssetPath, fileName)
	if err := ioutil.WriteFile(savePath, data, 0644); err != nil {
		dir, _ := filepath.Split(savePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error create directory: %s", err.Error())
		}
		if err := ioutil.WriteFile(savePath, data, 0644); err != nil {
			return fmt.Errorf("error save file: %s", err.Error())
		}

	}
	// base := filepath.Base(fileName)
	if strings.HasSuffix(savePath, "msgpack") || strings.HasPrefix(filepath.Base(savePath), "chart_") {
		if err := RunD4DJTool(savePath); err != nil {
			return err
		}
		if err := os.Remove(savePath); err != nil {
			fmt.Println(err)
		}

	} else if strings.Contains(fileName, "ondemand_card_chara_transparent") || strings.Contains(fileName, "ondemand_live2d_") {
		if err := RunExtractor(savePath); err != nil {
			return err
		}
		if err := os.Remove(savePath); err != nil {
			fmt.Println(err)
		}
	}
	var err error
	key := fileName
	switch {
	case strings.Contains(fileName, "msgpack"):
		targetFile := strings.Replace(savePath, "msgpack", "json", 1)
		key = strings.Replace(fileName, "msgpack", "json", 1)
		data, err = ioutil.ReadFile(targetFile)
	case strings.Contains(fileName, "chart"):
		targetFile := savePath + ".json"
		data, err = ioutil.ReadFile(targetFile)
	case strings.Contains(fileName, "ondemand_live2d_"):
		targetDirectory := strings.Replace(strings.Replace(savePath, "iOS", "Live2D", 1), "ondemand_", "", 1)
		filepath.Walk(targetDirectory,
			func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return err
				}
				name := info.Name()
				key := strings.Replace(strings.Replace(fileName, "iOS", "Live2D", 1), "ondemand_", "", 1) + "/"
				if strings.HasSuffix(name, "exp3.json") {
					key += "expressions/"
				} else if strings.Contains(name, "motion3.json") {
					key += "motions/"
				} else if strings.HasPrefix(name, "texture") {
					key += "textures/"
				}
				file, err := ioutil.ReadFile(path)
				if err != nil {
					return nil
				}
				awsutil.PutFile(key+name, bytes.NewReader(file))
				return nil
			})
		return nil
	case strings.Contains(fileName, "card_chara_transparent_"):
		name := strings.Replace(strings.Replace(savePath, "iOS", "images", 1), "ondemand_", "", 1) + ".png"
		keyName := strings.Replace(strings.Replace(fileName, "iOS", "images", 1), "ondemand_", "", 1) + ".png"
		data, err := ioutil.ReadFile(name)
		if err != nil {
			return err
		}
		awsutil.PutFile(keyName, bytes.NewReader(data))
		return nil
	case strings.Contains(fileName, "AssetBundles"):
		return nil
	}

	if err != nil {
		return err
	}
	if err := awsutil.PutFile(key, bytes.NewReader(data)); err != nil {
		return err
	}
	return nil

}

// RunD4DJTool convert Msgpack and chart_ to .json
func RunD4DJTool(filePath string) error {
	c := exec.Command("dotnet", conf.Get().ToolPath, filePath)
	// c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

// RunExtractor extract UnityAsset(live2d, character card image) to normal file
func RunExtractor(filePath string) error {
	c := exec.Command("dotnet", conf.Get().ExtractorPath, filePath)
	// c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
