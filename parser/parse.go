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
	"strconv"
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
		if total, err := RunVgmStream(savePath); err == nil {
			RunFFMPeg(total, savePath)
			trimed := strings.TrimSuffix(savePath, ".acb")
			files, err := filepath.Glob(trimed + "*" + ".wav")
			if err == nil {
				for _, f := range files {
					if err := os.Remove(f); err != nil {
						fmt.Println(err)
					}
				}
			}
			if conf.Get().Aws.BucketName != "" && total > 1 {
				if err := os.Remove(savePath); err != nil {
					fmt.Println(err)
				}
				dir := strings.ReplaceAll(filepath.Dir(key), "\\", "/")
				files, err := filepath.Glob(trimed + "*" + ".mp3")
				if err == nil {
					for _, f := range files {
						data, err = ioutil.ReadFile(f)
						if err != nil {
							return err
						}
						base := filepath.Base(f)
						if err := awsutil.PutFile(dir+"/"+base, bytes.NewReader(data)); err != nil {
							return err
						}
					}
				}
				return nil
			}

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
	case strings.Contains(fileName, "acb"):
		targetFile = strings.Replace(savePath, "acb", "mp3", 1)
		key = strings.Replace(fileName, "acb", "mp3", 1)
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
func RunVgmStream(filePath string) (int, error) {
	bytes, err := exec.Command(conf.Get().VgmStreamPath, "-I", filePath).Output()
	if err != nil {
		return -1, err
	}

	saveName := strings.TrimRight(filePath, ".acb")

	var data map[string]interface{}
	json.Unmarshal(bytes, &data)

	var total = 1

	streamInfo, ok := data["streamInfo"].(map[string]interface{})
	if ok {
		totalTemp, ok := streamInfo["total"].(float64)
		if ok {
			total = int(totalTemp)
		}
	}

	var cmd []string

	ext := ".wav"

	if total > 1 {
		cmd = []string{filePath, "-i", "-F", "-S", "0", "-o", saveName + "-?s" + ext}
	} else {
		cmd = []string{filePath, "-o", saveName + ext}

	}

	c := exec.Command(conf.Get().VgmStreamPath, cmd...)

	c.Stderr = os.Stderr
	return total, c.Run()
}

// RunFFMPeg convert .wav file to .mp3
func RunFFMPeg(total int, filePath string) error {
	if total > 1 {
		fileName := strings.TrimRight(filePath, ".acb")
		for i := 0; i < total; i++ {
			base := fileName + "-" + strconv.Itoa(i+1)
			originName := base + ".wav"
			saveName := base + ".mp3"

			c := exec.Command(conf.Get().FfmpegPath, "-y", "-i", originName, "-codec:a", "libmp3lame", saveName)
			c.Run()
		}
		return nil
	}

	c := exec.Command(conf.Get().FfmpegPath, "-y", "-i", strings.Replace(filePath, "acb", "wav", 1), "-codec:a", "libmp3lame", strings.Replace(filePath, "acb", "mp3", 1))

	c.Stderr = os.Stderr
	return c.Run()
}
