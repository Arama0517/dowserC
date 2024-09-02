package main

import (
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/Arama0517/dowserC/pkg/launcher"
	"github.com/caarlos0/log"
)

const version = "devel"

func main() {
	// debug模式
	debug := flag.Bool("debug", false, "开启调试模式")
	flag.Parse()
	if *debug {
		log.SetLevel(log.DebugLevel)
	}
	log.Debugf("dowserC 版本: %s", version)

	// 检测无效Mod
	log.Debug("开始检测Mod是否有效")
	log.IncreasePadding()
	files, err := filepath.Glob(filepath.Join(launcher.DataDir, "mod", "*.mod"))
	if err != nil {
		log.WithError(err).Fatal("获取mod定位文件失败")
	}
	for _, file := range files {
		byteData, err := os.ReadFile(file)
		if err != nil {
			log.WithError(err).Fatal("读取mod定位文件失败")
		}
		re := regexp.MustCompile(`path\s*=\s*"([^"]+)"`)
		match := re.FindStringSubmatch(string(byteData))
		if match == nil {
			os.Remove(file)
			log.WithField("path", file).Debug("删除无效的mod定位文件")
			continue
		}
		path := match[1]
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Remove(file)
			log.WithField("path", file).Debug("删除无效的mod定位文件")
		}
	}
	log.ResetPadding()

	entries, err := os.ReadDir(launcher.ModDir)
	if err != nil {
		log.WithError(err).Fatal("获取Mod失败, 可能是因为权限不足")
	}
	modNames := []string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		modName, err := launcher.GenerateModBootFile(filepath.Join(launcher.ModDir, entry.Name()))
		if err != nil {
			continue
		}
		modNames = append(modNames, modName)
	}
	if len(modNames) != 0 {
		log.Infof("共计加载: %d个Mod", len(modNames))
		log.Info("他们分别是:")
		log.IncreasePadding()
		for _, modName := range modNames {
			log.Infof("%s", modName)
		}
		log.ResetPadding()
	}
	log.Info("即将启动游戏, 3秒后自动退出.. .")
	if err := launcher.RunDowser(); err != nil {
		log.WithError(err).Fatal("启动客户端失败, 可能是因为权限不足")
	}
	time.Sleep(3 * time.Second)
}
