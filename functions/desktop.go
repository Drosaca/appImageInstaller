package functions

import (
	"appinstall/desktopFile"
	"appinstall/structs"
	"appinstall/utils"
	"fmt"
	"log"
)

func FindInternalDesktop(extractPath string) string {
	desktopPath, err := utils.FindFile(extractPath, []string{".desktop"})
	if err != nil {
		log.Fatal("fail to find desktop file:", err)
	}
	return desktopPath
}

func GenerateDesktopFile(desktopPath string) (*desktopFile.DesktopFile, error) {

	desktop := desktopFile.New()
	err := desktop.FromFile(desktopPath)
	if err != nil {
		return nil, fmt.Errorf("fail to open desktop file: %s", err)
	}
	return desktop, nil
}

func EditDesktop(desktop *desktopFile.DesktopFile, config structs.Config) {
	err := desktop.Category("Desktop Entry").Set("Exec", config.ExecPath)
	if err != nil {
		log.Fatal("fail to edit Desktop: ", err)
	}
}
