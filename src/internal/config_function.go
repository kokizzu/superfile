package internal

import (
	"embed"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/barasher/go-exiftool"
	"github.com/pelletier/go-toml/v2"
	varibale "github.com/yorukot/superfile/src/config"
)

func initialConfig(dir string) (toggleDotFileBool bool, firstFilePanelDir string) {
	var err error

	logOutput, err = os.OpenFile(varibale.LogFilea, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error while opening superfile.log file: %v", err)
	}

	loadConfigFile()

	loadHotkeysFile()

	loadThemeFile()

	toggleDotFileData, err := os.ReadFile(varibale.ToggleDotFilea)
	if err != nil {
		outPutLog("Error while reading toggleDotFile data error:", err)
	}
	if string(toggleDotFileData) == "true" {
		toggleDotFileBool = true
	} else if string(toggleDotFileData) == "false" {
		toggleDotFileBool = false
	}
	LoadThemeConfig()

	if Config.Metadata {
		et, err = exiftool.NewExiftool()
		if err != nil {
			outPutLog("Initial model function init exiftool error", err)
		}
	}

	firstFilePanelDir = varibale.HomeDir
	if dir != "" {
		firstFilePanelDir, err = filepath.Abs(dir)
		if err != nil {
			firstFilePanelDir = varibale.HomeDir
		}
	}
	return toggleDotFileBool, firstFilePanelDir
}

func loadConfigFile() {

	_ = toml.Unmarshal([]byte(ConfigTomlString), &Config)
	tempForCheckMissingConfig := ConfigType{}

	data, err := os.ReadFile(varibale.ConfigFilea)
	if err != nil {
		log.Fatalf("Config file doesn't exist: %v", err)
	}

	_ = toml.Unmarshal(data, &tempForCheckMissingConfig)
	err = toml.Unmarshal(data, &Config)
	if err != nil {
		log.Fatalf("Error decoding config file ( your config file may have misconfigured ): %v", err)
	}

	if !reflect.DeepEqual(Config, tempForCheckMissingConfig) {
		tomlData, err := toml.Marshal(Config)
		if err != nil {
			log.Fatalf("Error encoding config: %v", err)
		}

		err = os.WriteFile(varibale.ConfigFilea, tomlData, 0644)
		if err != nil {
			log.Fatalf("Error writing config file: %v", err)
		}
	}
	if Config.FilePreviewWidth > 10 || Config.FilePreviewWidth == 1 {
		log.Fatalf("Config file file_preview_width invalidation")
	}
}

func loadHotkeysFile() {

	_ = toml.Unmarshal([]byte(HotkeysTomlString), &hotkeys)
	tempForCheckMissingConfig := HotkeysType{}
	data, err := os.ReadFile(varibale.HotkeysFilea)

	if err != nil {
		log.Fatalf("Config file doesn't exist: %v", err)
	}

	_ = toml.Unmarshal(data, &tempForCheckMissingConfig)
	err = toml.Unmarshal(data, &hotkeys)
	if err != nil {
		log.Fatalf("Error decoding hotkeys file ( your config file may have misconfigured ): %v", err)
	}

	if !reflect.DeepEqual(hotkeys, tempForCheckMissingConfig) {
		tomlData, err := toml.Marshal(hotkeys)
		if err != nil {
			log.Fatalf("Error encoding hotkeys: %v", err)
		}

		err = os.WriteFile(varibale.HotkeysFilea, tomlData, 0644)
		if err != nil {
			log.Fatalf("Error writing hotkeys file: %v", err)
		}
	}

}

func loadThemeFile() {
	data, err := os.ReadFile(varibale.ThemeFoldera + "/" + Config.Theme + ".toml")
	if err != nil {
		data = []byte(DefaultThemeString)
	}

	err = toml.Unmarshal(data, &theme)
	if err != nil {
		log.Fatalf("Error while decoding theme file( Your theme file may have errors ): %v", err)
	}
}

func LoadAllDefaultConfig(content embed.FS) {

	temp, err := content.ReadFile("src/superfileConfig/hotkeys.toml")
	if err != nil {
		return
	}
	HotkeysTomlString = string(temp)

	temp, err = content.ReadFile("src/superfileConfig/config.toml")
	if err != nil {
		return
	}
	ConfigTomlString = string(temp)

	temp, err = content.ReadFile("src/superfileConfig/theme/catppuccin.toml")
	if err != nil {
		return
	}
	DefaultThemeString = string(temp)

	currentThemeVersion, err := os.ReadFile(varibale.ThemeFileVersiona)

	if err != nil && !os.IsNotExist(err) {
		outPutLog("Error reading from file:", err)
		return
	}

	_, err = os.Stat(varibale.ThemeFoldera)

	if os.IsNotExist(err) {
		err := os.MkdirAll(varibale.ThemeFoldera, 0755)
		if err != nil {
			outPutLog("error create theme direcroty", err)
			return
		}
	} else if string(currentThemeVersion) == varibale.CurrentVersion {
		return
	}

	files, err := content.ReadDir("src/superfileConfig/theme")
	if err != nil {
		outPutLog("error read theme directory from embed", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		src, err := content.ReadFile(filepath.Join("src/superfileConfig/theme", file.Name()))
		if err != nil {
			outPutLog("error read theme file from embed", err)
			return
		}

		file, err := os.Create(filepath.Join(varibale.ThemeFoldera, file.Name()))
		if err != nil {
			outPutLog("error create theme file from embed", err)
			return
		}
		file.Write(src)
		defer file.Close()
	}

	os.WriteFile(varibale.ThemeFileVersiona, []byte(varibale.CurrentVersion), 0644)
}
