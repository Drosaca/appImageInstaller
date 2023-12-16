package main

import (
	f "appinstall/functions"
	"appinstall/manager"
	"appinstall/structs"
	"appinstall/utils"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func extractApp(appPath string) {
	cmd := exec.Command(appPath, "--appimage-extract")
	err := cmd.Run()
	if err != nil {
		log.Fatal("fail to extract: ", err)
	}
}

func createDirectories(config structs.Config) error {
	err := os.MkdirAll(config.ExtractDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(config.ExecDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(config.ImgPath, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.Chmod(config.ExecDir, 0777)
	if err != nil {
		return err
	}
	return nil
}

func setConfig(path string) structs.Config {
	config := structs.Config{
		AppExtractDir:   "squashfs-root",
		ExtractDir:      "/tmp/appInstaller",
		ExecDir:         "/usr/share/appImages/",
		GnomeDesktopDir: "/usr/share/applications/",
		Debug:           false,
		ImgPath:         "/usr/share/pixmaps/",
		InputPath:       path,
	}
	config.ExecPath = filepath.Join(config.ExecDir, filepath.Base(config.InputPath))
	config.AppExtractDir = filepath.Join(config.ExtractDir, "squashfs-root")
	config.InputDir = filepath.Dir(config.InputPath)
	config.InputFileName = filepath.Base(config.InputPath)
	return config
}

func preInstall(config structs.Config) error {
	err := createDirectories(config)
	if err != nil {
		log.Fatal("fail to create directories: ", err)
	}
	err = os.Chmod(config.InputPath, 0777)
	if err != nil {
		return err
	}
	err = utils.Copy(config.InputPath, config.ExecDir)
	if err != nil {
		return err
	}
	return nil
}

func install(config structs.Config) {
	err := preInstall(config)
	if err != nil {
		log.Fatal("setup failed: ", err)
	}
	os.Chdir(config.ExtractDir)
	extractApp(config.InputPath)
	desktopPath := f.FindInternalDesktop(config.AppExtractDir)
	desktop, err := f.GenerateDesktopFile(desktopPath)
	if err != nil {
		log.Fatal("fail to parse Desktop file: ", err)
	}
	f.EditDesktop(desktop, config)
	err = f.CopyImage(desktop, config)
	if err != nil {
		fmt.Println(err)
	}
	err = desktop.ToFile(filepath.Join(config.GnomeDesktopDir, filepath.Base(desktopPath)))
	if err != nil {
		log.Fatal("fail to write Desktop file ", err)
	}

}

func runInstallScript(appPath string) error {
	path, _ := filepath.Abs(appPath)
	_, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("install", appPath, ":", err)
	}
	config := setConfig(path)
	err = os.RemoveAll(config.ExtractDir)
	if err != nil {
		return fmt.Errorf("removing extract directory :", err)
	}
	install(config)
	err = os.RemoveAll(config.ExtractDir)
	if err != nil {
		return fmt.Errorf("removing extract dir: ", err)
	}
	return nil
}

func help() {
	fmt.Println("usage: sudo appinstall path/to/app.AppImage")
	fmt.Println("(after install use sudo update-desktop-database to reload gnome icons)")
	fmt.Println("other options: ")
	fmt.Println("-l          #to list installed apps (from this tool only)")
	fmt.Println(" ")
	fmt.Println("-d appName  #to delete the app (installed by this tool)")

}

func listing() error {
	config := setConfig("")
	m := manager.New(config)
	entries := m.List()
	fmt.Println("Apps:")
	for _, entry := range entries {
		name, _ := entry.Category("Desktop Entry").Get("Name")
		fmt.Println(name)
	}
	return nil
}

func deleteApp(appName string) error {
	config := setConfig("")
	m := manager.New(config)
	return m.Delete(appName)
}

func chooseScript() error {
	if os.Args[1] == "-l" {
		return listing()
	}
	if os.Args[1] == "-d" && len(os.Args) >= 3 {
		return deleteApp(os.Args[2])
	}
	if len(os.Args) == 2 {
		return runInstallScript(os.Args[1])
	}
	help()
	return nil
}

func main() {

	if os.Args[1] == "-h" {
		help()
		os.Exit(0)
	}
	if os.Args[1] == "-v" {
		fmt.Println("1.0")
		os.Exit(0)
	}
	err := chooseScript()
	if err != nil {
		fmt.Println(err)
	}
}
