package manager

import (
	"appImageInstaller/desktopFile"
	"appImageInstaller/structs"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Manager struct {
	config structs.Config
}

func New(config structs.Config) *Manager {
	manager := new(Manager)
	manager.config = config
	return manager
}

func (m *Manager) IsGeneratedDesktop(deskFile *desktopFile.DesktopFile) (bool, error) {
	isGood, err := m.IsValidDesktop(deskFile)
	if !isGood {
		return isGood, err
	}
	execPath, _ := deskFile.Category("Desktop Entry").Get("Exec")
	if !strings.Contains(execPath, m.config.ExecDir) {
		return false, fmt.Errorf("external exec path")
	}
	return true, nil
}

func (m *Manager) IsValidDesktop(deskFile *desktopFile.DesktopFile) (bool, error) {
	if !deskFile.HasValues("Desktop Entry", []string{"Get", "Name", "Exec"}) {
		return false, fmt.Errorf("missing basic values")
	}
	path, _ := deskFile.Category("Desktop Entry").Get("Exec")
	_, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("fail to find binary")
	}
	return true, nil
}

func (m *Manager) List() []*desktopFile.DesktopFile {
	var appList []*desktopFile.DesktopFile
	entries, err := os.ReadDir(m.config.GnomeDesktopDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range entries {
		deskFile := desktopFile.New()
		deskFilePath := filepath.Join(m.config.GnomeDesktopDir, e.Name())
		err := deskFile.FromFile(deskFilePath)
		if err != nil {
			fmt.Println("Listing Error: ", "file: ", deskFilePath, ":", err)
			continue
		}
		_, err = m.IsValidDesktop(deskFile)
		if err != nil {
			fmt.Println("Listing Error for file:", deskFilePath, " ", err)
			continue
		}
		isGenerated, err := m.IsGeneratedDesktop(deskFile)
		if !isGenerated {
			continue
		}
		appList = append(appList, deskFile)
	}
	return appList
}

func (m *Manager) Delete(appName string) error {
	entries, err := os.ReadDir(m.config.GnomeDesktopDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range entries {
		deskFile := desktopFile.New()
		deskFilePath := filepath.Join(m.config.GnomeDesktopDir, e.Name())
		err := deskFile.FromFile(deskFilePath)
		if err != nil {
			fmt.Println("Listing Error: ", "file: ", deskFilePath, ":", err)
			continue
		}
		isGenerated, err := m.IsGeneratedDesktop(deskFile)
		if !isGenerated {
			continue
		}
		ExecPath, _ := deskFile.Category("Desktop Entry").Get("Exec")
		name, _ := deskFile.Category("Desktop Entry").Get("Name")
		if name != appName {
			continue
		}
		err = os.Remove(ExecPath)
		if err != nil {
			return err
		}
		err = os.Remove(deskFilePath)
		if err != nil {
			return err
		}
		fmt.Println("removed binary file", ExecPath)
		fmt.Println("removed Desktop file", deskFilePath)
		return nil
	}
	return fmt.Errorf("App not found")
}
