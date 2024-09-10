package main

import (
	_ "embed"
	"os"

	"github.com/Arama0517/dowserC/cmd"
	goversion "github.com/caarlos0/go-version"
	"github.com/caarlos0/log"
)

var (
	version   = ""
	commit    = ""
	treeState = ""
	date      = ""
	builtBy   = ""
	website   = "https://github.com/Arama0517/dowserC"
	//go:embed art.txt
	asciiArt string
)

func main() {
	if err := cmd.Execute(os.Args[1:], buildVersion(version, commit, date, builtBy, treeState)); err != nil {
		log.Info("按下回车键退出...")
		_, _ = os.Stdin.Read(make([]byte, 1))
		os.Exit(1)
	}
}

func buildVersion(version, commit, date, builtBy, treeState string) goversion.Info {
	return goversion.GetVersionInfo(
		goversion.WithAppDetails("dowserC", "根据游戏根目录下的mod文件夹自动生成mod定位文件", website),
		goversion.WithASCIIName(asciiArt),
		func(i *goversion.Info) {
			if commit != "" {
				i.GitCommit = commit
			}
			if treeState != "" {
				i.GitTreeState = treeState
			}
			if date != "" {
				i.BuildDate = date
			}
			if version != "" {
				i.GitVersion = version
			}
			if builtBy != "" {
				i.BuiltBy = builtBy
			}
		},
	)
}
