package parser

import (
	"bytes"
	"encoding/json"
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

type vgmStreamOutput struct {
	StreamInfo struct {
		Name  string  `json:"name"`
		Total float64 `json:"total"`
	} `json:"StreamInfo"`
}

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
		if list, err := RunVgmStream(savePath); err == nil && list != nil {
			os.Remove(savePath)

			head := strings.ReplaceAll(filepath.Dir(fileName), "\\", "/")
			saveHead := strings.ReplaceAll(filepath.Dir(savePath), "\\", "/")

			for _, d := range list {
				wav := saveHead + "/" + d + ".wav"
				upload := saveHead + "/" + d + ".mp3"
				key := head + "/" + d + ".mp3"

				fmt.Println(upload)
				if err := RunFFMPeg(wav); err != nil {
					fmt.Println(err)
					continue
				}
				os.Remove(wav)
				data, err = ioutil.ReadFile(upload)
				if err != nil {
					return err
				}
				if err := awsutil.PutFile(key, bytes.NewReader(data)); err != nil {
					fmt.Println(err)
				}
			}
			return nil
		} else {
			fmt.Println(err)
		}
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

	if conf.Get().Aws.BucketName != "" && strings.Contains(fileName, "ondemand_live2d_") {
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
				awsutil.PutFile(key, bytes.NewReader(file))
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
	case strings.Contains(fileName, "card_chara_transparent_"):
		targetFile = strings.Replace(strings.Replace(savePath, "iOS", "images", 1), "ondemand_", "", 1) + ".png"
		key = strings.Replace(strings.Replace(fileName, "iOS", "images", 1), "ondemand_", "", 1) + ".png"
	case strings.Contains(fileName, "AssetBundles"):
		return nil
	}

	if conf.Get().Aws.BucketName != "" {
		data, err = ioutil.ReadFile(targetFile)
		if err != nil {
			return err
		}
		if err := awsutil.PutFile(key, bytes.NewReader(data)); err != nil {
			return err
		}
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
func RunVgmStream(filePath string) ([]string, error) {
	saveName := strings.TrimRight(filePath, ".acb")

	base := filepath.Base(saveName)
	dir := filepath.Dir(saveName)

	bytes, err := exec.Command(conf.Get().VgmStreamPath, "-m", "-I", filePath).Output()
	if err != nil {
		return nil, err
	}
	var data vgmStreamOutput
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}

	if data.StreamInfo.Total <= 1 {
		saveName += ".wav"
		c := exec.Command(conf.Get().VgmStreamPath, filePath, "-i", "-F", "-o", saveName)
		return []string{base}, c.Run()
	}

	bytes, err = exec.Command(conf.Get().VgmStreamPath, filePath, "-i", "-I", "-F", "-S", "0", "-o", path.Join(dir, base+"-?n.wav")).Output()
	if err != nil {
		return nil, err
	}
	list := []string{}
	splited := strings.Split(string(bytes), "\n")

	for _, d := range splited {
		var data vgmStreamOutput
		if err := json.Unmarshal([]byte(d), &data); err != nil {
			continue
		}
		list = append(list, base+"-"+data.StreamInfo.Name)
	}

	return list, nil
}

// RunFFMPeg convert .wav file to .mp3
func RunFFMPeg(filePath string) error {
	c := exec.Command(conf.Get().FfmpegPath, "-y", "-i", filePath, "-codec:a", "libmp3lame", strings.Replace(filePath, "wav", "mp3", 1))

	c.Stderr = os.Stderr
	return c.Run()
}
