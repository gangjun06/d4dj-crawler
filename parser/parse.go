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

func Parse(fileName string, data []byte, extSavePath ...string) error {

	dirPath := conf.Get().AssetPath
	if len(extSavePath) > 0 {
		dirPath = extSavePath[0]
	}
	savePath := path.Join(dirPath, fileName)
	if err := ioutil.WriteFile(savePath, data, 0644); err != nil {
		dir, _ := filepath.Split(savePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error create directory: %s", err.Error())
		}
		if err := ioutil.WriteFile(savePath, data, 0644); err != nil {
			return fmt.Errorf("error save file: %s", err.Error())
		}
	}

	var err error
	key := fileName
	usedExternalTool := true

	if strings.HasSuffix(savePath, "msgpack") || strings.HasPrefix(filepath.Base(savePath), "chart_") {
		err = RunD4DJTool(savePath)
	} else if strings.Contains(fileName, "ondemand_card_chara_transparent") || strings.Contains(fileName, "ondemand_live2d_") {
		err = RunExtractor(savePath)
	} else if strings.HasSuffix(fileName, "acb") {
		err = RunVgmStream(savePath)
	} else {
		usedExternalTool = false
	}

	if usedExternalTool {
		if err != nil {
			return err
		}
		if err := os.Remove(savePath); err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(fileName, "ondemand_live2d_") {
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
					return err
				}
				key += name
				awsutil.PutFile(key+name, bytes.NewReader(file))
				return nil
			})
		return nil
	}

	targetFile := savePath
	switch {
	case strings.Contains(fileName, "msgpack"):
		targetFile = strings.Replace(savePath, "msgpack", "json", 1)
		key = strings.Replace(fileName, "msgpack", "json", 1)
	case strings.Contains(fileName, "chart"):
		targetFile = savePath + ".json"
	case strings.Contains(fileName, "acb"):
		targetFile = strings.Replace(savePath, "acb", "wav", 1)
		key = strings.Replace(fileName, "acb", "wav", 1)
	case strings.Contains(fileName, "card_chara_transparent_"):
		targetFile = strings.Replace(strings.Replace(savePath, "iOS", "images", 1), "ondemand_", "", 1) + ".png"
		key = strings.Replace(strings.Replace(fileName, "iOS", "images", 1), "ondemand_", "", 1) + ".png"
	case strings.Contains(fileName, "AssetBundles"):
		return nil
	}

	data, err = ioutil.ReadFile(targetFile)
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
	c.Stderr = os.Stderr
	return c.Run()
}

// RunExtractor extract UnityAsset(live2d, character card image) to normal file
func RunExtractor(filePath string) error {
	c := exec.Command("dotnet", conf.Get().ExtractorPath, filePath)
	c.Stderr = os.Stderr
	return c.Run()
}

// RunVGMStream convert .acb file to .wav
func RunVgmStream(filePath string) error {
	c := exec.Command(conf.Get().VgmStreamPath, "-o", strings.Replace(filePath, "acb", "wav", 1), filePath)
	c.Stderr = os.Stderr
	return c.Run()
}
