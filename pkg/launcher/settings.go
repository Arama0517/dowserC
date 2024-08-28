package launcher

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/apex/log"
)

var (
	CWD     = "" // 当前目录
	DataDir = "" // 游戏的数据目录(存放mod引导文件和存档等等)
	ModDir  = "" // 存放Mod的目录
)

var LauncherSettingsPath = "" // 启动器配置文件路径

//nolint:tagliatelle
type alternativeExecutable struct {
	Label   map[string]string `json:"label"`
	ExePath string            `json:"exePath"`
	ExeArgs []string          `json:"exeArgs"`
}

//nolint:tagliatelle
type launcherSettings struct {
	GameID                   string                  `json:"gameId"`
	DisplayName              string                  `json:"displayName"`
	Version                  string                  `json:"version"`
	RawVersion               string                  `json:"rawVersion"`
	DistPlatform             string                  `json:"distPlatform"`
	IngameSettingsPath       string                  `json:"ingameSettingsPath"`
	GameDataPath             string                  `json:"gameDataPath"`
	DLCPath                  string                  `json:"dlcPath"`
	IngameSettingsLayoutPath string                  `json:"ingameSettingsLayoutPath"`
	ThemeFile                string                  `json:"themeFile"`
	BrowserDlcURL            string                  `json:"browserDlcUrl"`
	BrowserModURL            string                  `json:"browserModUrl"`
	ExePath                  string                  `json:"exePath"`
	ExeArgs                  []string                `json:"exeArgs"`
	AlternativeExecutables   []alternativeExecutable `json:"alternativeExecutables"`
	GameCachePaths           []string                `json:"gameCachePaths"`
}

func init() {
	// 定义当前目录
	if c, err := os.Getwd(); err == nil {
		CWD = c
	} else {
		log.WithError(err).Fatal("获取当前目录失败")
	}

	// 定义启动器配置文件路径
	LauncherSettingsPath = filepath.Join(CWD, "launcher-settings.json")

	// 定义Mod目录
	ModDir = filepath.Join(CWD, "mod")
	_ = os.MkdirAll(ModDir, 0o755)

	// 检查启动器配置文件是否存在
	if file, err := os.Open(LauncherSettingsPath); os.IsNotExist(err) {
		log.Error("启动器配置文件不存在")
		log.Error("请确认本程序是否在游戏的根目录下运行...")
		log.Fatalf("当前路径: %s", CWD)
	} else if err != nil {
		log.WithError(err).Fatal("读取启动器配置文件失败, 可能是因为权限不足")
	} else {
		// 读取启动器配置文件
		byteData, err := io.ReadAll(file)
		if err != nil {
			log.WithError(err).Fatal("读取启动器配置文件为字节失败, 可能是因为文件损坏")
		}
		var settings launcherSettings
		if err := json.Unmarshal(byteData, &settings); err != nil {
			log.Fatal("解析启动器配置文件为JSON失败, 可能是因为格式错误或者为非法内容")
		}
		DataDir = filepath.Join(CWD, "Paradox Interactive", settings.DisplayName)
		_ = os.MkdirAll(DataDir, 0o755)
		file.Close()

		// 修改数据目录
		if settings.GameDataPath != DataDir {
			settings.GameDataPath = DataDir
			byteData, err = json.MarshalIndent(settings, "", "    ")
			if err != nil {
				log.WithError(err).Fatal("生成启动器配置文件失败, 可能是因为编码失败")
			}
			if err := os.WriteFile(LauncherSettingsPath, byteData, 0o755); err != nil {
				log.WithError(err).Fatal("写入启动器配置文件失败, 可能是因为权限错误")
			}
		}
	}
}
