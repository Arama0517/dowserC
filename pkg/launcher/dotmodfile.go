package launcher

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/apex/log"
)

var ErrNoNameFound = errors.New("没有找到Mod的名称")

func GenerateDotModFile(path string) (string, error) {
	byteData, err := os.ReadFile(filepath.Join(path, "descriptor.mod"))
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`name\s*=\s*"([^"]+)"`)
	match := re.FindStringSubmatch(string(byteData))
	if match == nil {
		return "", ErrNoNameFound
	}
	modName := match[1]
	dotModPath := filepath.Join(DataDir, "mod", fmt.Sprintf("%s.mod", modName))
	if _, err := os.Stat(dotModPath); os.IsNotExist(err) {
		file, err := os.Create(dotModPath)
		if err != nil {
			return "", err
		}
		_, err = file.Write(byteData)
		if err != nil {
			return "", err
		}
		_, err = file.WriteString(fmt.Sprintf("\npath=\"%s\"", filepath.ToSlash(path)))
		if err != nil {
			return "", err
		}
		file.Close()
	} else if err != nil {
		log.WithError(err).Fatal("无法读取mod文件")
	}
	return modName, nil
}
