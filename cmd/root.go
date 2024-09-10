package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/Arama0517/dowserC/pkg/launcher"
	goversion "github.com/caarlos0/go-version"
	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
)

func Execute(args []string, version goversion.Info) error {
	// 重置 Log 的缩进
	log.ResetPadding()

	// 初始化
	if err := launcher.InitSettings(); err != nil {
		return err
	}

	// 运行
	cmd := newRootCmd(version)
	cmd.SetArgs(args)
	return cmd.Execute()
}

func newRootCmd(version goversion.Info) *cobra.Command {
	var debug bool
	cmd := &cobra.Command{
		Use:     "dowserC",
		Short:   "根据游戏根目录下的mod文件夹自动生成mod定位文件",
		Version: version.String(),
		Run: func(*cobra.Command, []string) {
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
			cmd := exec.Command(launcher.DowserPath)
			cmd.Dir = launcher.CWD
			if err := cmd.Start(); err != nil {
				log.WithError(err).Fatal("启动客户端失败, 可能是因为权限不足")
			}
			time.Sleep(3 * time.Second)
		},
		PersistentPreRun: func(*cobra.Command, []string) {
			if debug {
				log.SetLevel(log.DebugLevel)
			}
		},
	}
	cmd.SetVersionTemplate("{{.Version}}")
	cmd.PersistentFlags().BoolVar(&debug, "debug", false, "开启调试模式")
	return cmd
}
