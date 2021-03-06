package main

import (
	"fmt"

	"github.com/okayu3/gosmartana/internal/logic"
	"github.com/pelletier/go-toml"
)

func main() {
	mstDir, outDir, settingDir, tokutaiIki, ninkeiIki := loadSettings()
	logic.RunLogic(mstDir, outDir, settingDir, tokutaiIki, ninkeiIki)
}

func loadSettings() (string, string, string, string, string) {
	settings, _ := toml.LoadFile("./settings.toml")
	if settings == nil {
		fmt.Println("cant load settings.toml file")
		return "C:/task/prj/YG01/mst/", "C:/Users/woodside3/go/output/",
			"C:/Users/woodside3/go/settings/", "", ""
	}
	mstDir := settings.Get("MasterPath.MST_DIR").(string)
	outDir := settings.Get("OutputPath.OUT_DIR").(string)
	setDir := settings.Get("SettingsPath.SETTING_DIR").(string)

	tokutaiInsKigo := settings.Get("Tokutai.TOKUTAI_INS_KIGO").(string)
	ninkeiInsKigo := settings.Get("NiniKeizoku.NINKEI_INS_KIGO").(string)
	return mstDir, outDir, setDir, tokutaiInsKigo, ninkeiInsKigo
}
