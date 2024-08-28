package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/Arama0517/dowserC/pkg/launcher"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

func init() {
	log.SetHandler(cli.Default)
	log.SetLevel(log.DebugLevel)
}

func main() {
	entries, err := os.ReadDir(launcher.ModDir)
	if err != nil {
		log.WithError(err).Fatal("获取Mod失败, 可能是因为权限不足")
	}
	modNames := []string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		modName, err := launcher.GenerateDotModFile(filepath.Join(launcher.ModDir, entry.Name()))
		if err != nil {
			continue
		}
		modNames = append(modNames, modName)
	}
	log.Infof("共计加载: %d个Mod", len(modNames))
	log.Info("他们分别是:")
	for _, modName := range modNames {
		log.Infof("    %s", modName)
	}
	log.Info("即将启动游戏, 3秒后自动退出.. .")
	if err := launcher.RunDowser(); err!= nil {
		log.WithError(err).Fatal("启动客户端失败, 可能是因为权限不足")
	}
	time.Sleep(3 * time.Second)
}
