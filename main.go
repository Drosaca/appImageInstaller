package main

import (
	f "appImageInstaller/functions"
	"appImageInstaller/structs"
	"appImageInstaller/utils"
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

func createDirectories(config structs.Config) {
	os.MkdirAll(config.ExtractDir, os.ModePerm)
	os.MkdirAll(config.ExecDir, os.ModePerm)
}

func setConfig(path string) structs.Config {
	config := structs.Config{
		AppExtractDir:   "squashfs-root",
		ExtractDir:      "/tmp/appInstaller",
		ExecDir:         "/usr/share/appImages/",
		GnomeDesktopDir: "/usr/share/applications/",
		ImgPath:         "/usr/share/pixmaps/",
		InputPath:       path,
	}
	config.ExecPath = filepath.Join(config.ExecDir, filepath.Base(config.InputPath))
	config.AppExtractDir = filepath.Join(config.ExtractDir, "squashfs-root")
	createDirectories(config)
	return config
}

func install(config structs.Config) {
	os.Chdir(config.ExtractDir)
	utils.Copy(config.InputPath, config.ExecPath)
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

func main() {

	if len(os.Args) != 2 || os.Args[1] == "-h" {
		fmt.Println("usage :appinstall path")
	}
	path, _ := filepath.Abs(os.Args[1])
	config := setConfig(path)
	os.RemoveAll(config.ExtractDir)
	install(config)
	os.RemoveAll(config.ExtractDir)
}
