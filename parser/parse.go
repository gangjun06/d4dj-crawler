package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

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

	}
	//  else if strings.HasPrefix(base, "ondemand_card_chara_transparent") || strings.HasPrefix(base, "ondemand_live2d_") {
	// 	if err := RunExtractor(savePath); err != nil {
	// 		return err
	// 	}
	// 	if err := os.Remove(savePath); err != nil {
	// 		fmt.Println(err)
	// 	}
	// }

	return nil

}

// RunD4DJTool convert Msgpack and chart_ to .json
func RunD4DJTool(filePath string) error {
	c := exec.Command("dotnet", conf.Get().ToolPath, filePath)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

// RunExtractor extract UnityAsset(live2d, character card image) to normal file
func RunExtractor(filePath string) error {
	c := exec.Command("dotnet", conf.Get().ToolPath, filePath)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
