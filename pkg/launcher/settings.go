package launcher

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/caarlos0/log"
)

var (
	CWD     = "" // 当前目录
	DataDir = "" // 游戏的数据目录(存放mod引导文件和存档等等)
	ModDir  = "" // 存放Mod的目录
)

var (
	LauncherSettingsPath = "" // 启动器配置文件路径
	DowserPath           = "" // dowser
)

//nolint:tagliatelle
type launcherSettings struct {
	GameID      string `json:"gameId"`
	DisplayName string `json:"displayName"`
	Version     string `json:"version"`
	RawVersion  string `json:"rawVersion"`
	// DistPlatform             string   `json:"distPlatform"`
	// IngameSettingsPath       string   `json:"ingameSettingsPath"`
	GameDataPath string `json:"gameDataPath"`
	// DLCPath                  string   `json:"dlcPath"`
	// IngameSettingsLayoutPath string   `json:"ingameSettingsLayoutPath"`
	// ThemeFile                string   `json:"themeFile"`
	// BrowserDlcURL            string   `json:"browserDlcUrl"`
	// BrowserModURL            string   `json:"browserModUrl"`
	// ExePath                  string   `json:"exePath"`
	// ExeArgs                  []string `json:"exeArgs"`
	// AlternativeExecutables   []struct {
	// 	Label   map[string]string `json:"label"`
	// 	ExePath string            `json:"exePath"`
	// 	ExeArgs []string          `json:"exeArgs"`
	// } `json:"alternativeExecutables"`
	// GameCachePaths []string `json:"gameCachePaths"`
}

func InitSettings() error {
	// 定义当前目录
	if c, err := os.Getwd(); err == nil {
		CWD = c
	} else {
		log.WithError(err).Error("获取当前目录失败")
		return err
	}

	// 定义启动器配置文件路径
	LauncherSettingsPath = filepath.Join(CWD, "launcher-settings.json")
	DowserPath = filepath.Join(CWD, "dowser.exe")

	// 检查dowser是否存在
	switch _, err := os.Stat(DowserPath); err {
	case os.ErrNotExist: // 不存在
		log.Error("dowser.exe 不存在")
		log.Error("请确认本程序是否在游戏的根目录下运行...")
		log.Errorf("当前路径: %s", CWD)
		return err
	case os.ErrPermission: // 权限错误
		log.Error("你没有读取 dowser.exe 的权限, 请检查当前账户是否有访问游戏文件的权限...")
		log.Errorf("当前路径: %s", CWD)
		return err
	case nil: // 没有错误
		break
	default: // 其他情况
		log.WithError(err).Error("获取 dowser.exe 的信息失败")
		return err
	}

	// 检查启动器配置文件是否存在
	file, err := os.Open(LauncherSettingsPath)
	// 检查文件是否存在
	// if os.IsNotExist(err) {
	// 	log.Error("启动器配置文件不存在")
	// 	log.Error("请确认本程序是否在游戏的根目录下运行...")
	// 	log.Errorf("当前路径: %s", CWD)
	// 	return err
	// } else if err != nil { // 其他情况
	// 	log.WithError(err).Error("读取启动器配置文件失败, 可能是因为权限不足")
	// 	return err
	// }
	switch err {
	case os.ErrNotExist: // 不存在
		log.Error("启动器配置文件不存在")
		log.Error("请确认本程序是否在游戏的根目录下运行...")
		log.Errorf("当前路径: %s", CWD)
		return err
	case os.ErrPermission: // 权限错误
		log.Error("你没有读取启动器配置文件的权限, 请检查当前账户是否有访问游戏文件的权限...")
		log.Errorf("当前路径: %s", CWD)
		return err
	case nil: // 没有错误
		break
	default: // 其他情况
		log.WithError(err).Error("读取启动器配置文件失败, 可能是因为权限不足")
	}

	// 读取启动器配置文件
	byteData, err := io.ReadAll(file)
	if err != nil {
		log.WithError(err).Error("读取启动器配置文件为字节失败, 可能是因为文件损坏")
		return err
	}
	var settings launcherSettings
	if err := json.Unmarshal(byteData, &settings); err != nil {
		log.WithError(err).Error("解析启动器配置文件为JSON失败, 可能是因为格式错误或者为非法内容")
		return err
	}
	DataDir = filepath.Join(CWD, "Paradox Interactive", settings.DisplayName)
	_ = os.MkdirAll(DataDir, 0o755)
	file.Close()

	// 修改数据目录
	if settings.GameDataPath != DataDir {
		settings.GameDataPath = DataDir
		byteData, err = json.MarshalIndent(settings, "", "    ")
		if err != nil {
			log.WithError(err).Error("生成启动器配置文件失败, 可能是因为编码失败")
			return err
		}
		if err := os.WriteFile(LauncherSettingsPath, byteData, 0o755); err != nil {
			log.WithError(err).Error("写入启动器配置文件失败, 可能是因为权限错误")
			return err
		}
	}

	// 定义Mod目录
	ModDir = filepath.Join(CWD, "mod")
	_ = os.MkdirAll(ModDir, 0o755)
	return nil
}
