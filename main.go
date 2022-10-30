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

func createDirectories(config structs.Config) error {
	err := os.MkdirAll(config.ExtractDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(config.ExecDir, os.ModePerm)
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

func main() {

	if len(os.Args) != 2 || os.Args[1] == "-h" {
		fmt.Println("usage :appinstall path")
	}
	path, _ := filepath.Abs(os.Args[1])
	config := setConfig(path)
	err := os.RemoveAll(config.ExtractDir)
	if err != nil {
		log.Fatal("removing", err)
	}
	install(config)
	err = os.RemoveAll(config.ExtractDir)
	if err != nil {
		log.Fatal("removing: ", err)
	}
}
