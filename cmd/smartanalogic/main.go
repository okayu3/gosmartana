package main

import (
	"fmt"

	"github.com/okayu3/gosmartana/internal/logic"
	"github.com/pelletier/go-toml"
)

func main() {
	_, outDir, settingDir := loadSettings()
	logic.RunLogic(outDir, settingDir)
}

func loadSettings() (string, string, string) {
	settings, _ := toml.LoadFile("./settings.toml")
	if settings == nil {
		fmt.Println("cant load settings.toml file")
		return "C:/task/prj/YG01/mst/", "C:/Users/woodside3/go/output/",
			"C:/Users/woodside3/go/settings/"
	}
	mstDir := settings.Get("MasterPath.MST_DIR").(string)
	outDir := settings.Get("OutputPath.OUT_DIR").(string)
	setDir := settings.Get("SettingsPath.SETTING_DIR").(string)
	return mstDir, outDir, setDir
}
