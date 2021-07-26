package file

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gangjun06/d4dj-info-server/conf"
)

type Explorer struct {
	Path       string
	TargetPath string
	IsDir      bool
	NotFound   bool
}
type FileList struct {
	Name  string
	Path  string
	IsDir bool
	Size  string
	Ext   string
}

func NewExplorer(targetPath string) *Explorer {
	targetPath = strings.ReplaceAll(targetPath, "..", "")
	realPath := path.Join(conf.Get().AssetPath, targetPath)
	explorer := &Explorer{
		Path:       realPath,
		TargetPath: targetPath,
		NotFound:   false,
	}
	info, err := os.Stat(realPath)
	if err != nil {
		explorer.NotFound = true
	} else {
		explorer.IsDir = info.IsDir()
	}
	return explorer
}

func (e *Explorer) FileList() []*FileList {
	info, _ := ioutil.ReadDir(e.Path)
	var list []*FileList
	for _, d := range info {
		name := d.Name()
		if d.IsDir() {
			name += "/"
		}
		list = append(list, &FileList{
			Name:  name,
			Path:  path.Join(e.TargetPath, d.Name()),
			Size:  strconv.Itoa(int(d.Size())),
			IsDir: d.IsDir(),
			Ext:   filepath.Ext(d.Name())})
	}
	return list
}

func (e *Explorer) FileData() ([]byte, error) {
	return ioutil.ReadFile(e.Path)
}
