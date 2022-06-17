package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func findFile(path string, patterns []string) (string, error) {
	foundPath := ""
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			for _, pattern := range patterns {
				if strings.Contains(filepath.Base(path), pattern) {
					foundPath = path
					return io.EOF
				}
			}
		}
		return nil
	})
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return foundPath, err
	}
	return foundPath, nil
}
func findFiles(path string, patterns []string) (string, error) {
	foundPath := ""
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			for _, pattern := range patterns {
				if strings.Contains(filepath.Base(path), pattern) {
					foundPath = path
					return io.EOF
				}
			}
		}
		return nil
	})
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return foundPath, err
	}
	return foundPath, nil
}

func extractApp(path string) string {
	cmd := exec.Command(path, "--appimage-extract")
	err := cmd.Run()
	if err != nil {
		log.Fatal("fail to extract: ", err)
	}
	extractPath := filepath.Join("/tmp/", "squashfs-root")
	return extractPath
}

func parseCategory(desktop map[string]map[string]string, line string) (bool, string) {
	rCategory, _ := regexp.Compile(`\[(.+)\]`)
	if rCategory.MatchString(line) {
		category := rCategory.FindStringSubmatch(line)[1]
		desktop[category] = map[string]string{}
		return true, category
	}
	return false, ""
}
func parseParameter(desktop map[string]map[string]string, line string, category string) {
	rCategory, _ := regexp.Compile(`.+=.+`)
	if rCategory.MatchString(line) {
		tuple := strings.Split(line, "=")
		desktop[category][tuple[0]] = tuple[1]
	}
}

func parseInternalDesktop(scanner *bufio.Scanner) map[string]map[string]string {
	desktop := make(map[string]map[string]string)
	category := "root"
	for scanner.Scan() {
		line := scanner.Text()
		isCategory, cat := parseCategory(desktop, line)
		if isCategory {
			category = cat
		}
		parseParameter(desktop, line, category)
	}
	return desktop
}

func findInternalDesktop(extractPath string) map[string]map[string]string {
	desktopPath, err := findFile(filepath.Join(extractPath), []string{".desktop"})
	if err != nil {
		log.Fatal("fail to find desktop file:", err)
	}
	file, err := os.Open(desktopPath)
	if err != nil {
		log.Fatal("fail to open desktop file:", err)
	}
	scanner := bufio.NewScanner(file)
	return parseInternalDesktop(scanner)
}

func findImage(desktop map[string]map[string]string, path string) string {
	types := []string{".svg", ".png", ".jpg", ".jpeg", ".bmp", ".ico", ".cur", ".tif", ".tiff", ".webp", ".jfif", ".pjpeg", ".pjp", ".gif", ".avif", ".apng"}
	category := "Desktop Entry"
	if _, key := desktop[category]; !key {
		log.Fatal("fail to find image ", desktop)
	}
	path, err := findFile(path, types)
	if err != nil {
		log.Fatal("fail to find image: ", err)
	}
	return path
}

func copyImage(imagePath string) {
	cmd := exec.Command("cp", imagePath, "/usr/share/pixmaps/")
	err := cmd.Run()
	if err != nil {
		log.Fatal("fail to copy image ", err)
	}
}

func dumpToFile(path string, desktop map[string]map[string]string) {
	file, err := os.Create(path)

	if err != nil {
		log.Fatal("fail to create desktop file: ", err)
	}
	for category, parameter := range desktop {
		file.WriteString("[" + category + "]\n")
		for name, value := range parameter {
			file.WriteString(name + "=" + value + "\n")
		}
	}
}

func createDesktop(imagePath string, execPath string, oldDesktop map[string]map[string]string) {
	newDesktop := map[string]map[string]string{
		"Desktop Entry": {
			"Name": oldDesktop["Desktop Entry"]["Name"],
			"Icon": strings.TrimSuffix(filepath.Base(imagePath), filepath.Ext(imagePath)),
			"Exec": execPath,
			"Type": "Application",
		},
	}
	dumpToFile(filepath.Join("/usr/share/applications/", oldDesktop["Desktop Entry"]["Name"]+".desktop"), newDesktop)
}

func copyApp(path string, destPath string) {
	cmd := exec.Command("mkdir", filepath.Dir(destPath))
	err := cmd.Run()
	cmd = exec.Command("cp", path, destPath)
	err = cmd.Run()
	if err != nil {
		log.Fatal("fail to copy binary ", err)
	}
	cmd = exec.Command("cp", path, "/tmp")
	err = cmd.Run()
}
func clean(path string) {
	cmd := exec.Command("rm", "-rf", filepath.Join(path, "squashfs-root"))
	err := cmd.Run()
	if err != nil {
		log.Fatal("fail to clean: ", err)
	}
}

func main() {
	if len(os.Args) != 2 || os.Args[1] == "-h" {
		fmt.Println("usage :appinstall path")
	}
	path, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal("Path resolution failed: ", err)
		return
	}
	os.Chdir("/tmp")
	clean("/tmp")
	execPath := filepath.Join("/usr/share/appImages/", filepath.Base(path))
	copyApp(path, execPath)
	extractPath := extractApp(filepath.Join("/tmp", filepath.Base(path)))
	desktop := findInternalDesktop(extractPath)
	imgPath := findImage(desktop, extractPath)
	copyImage(imgPath)
	createDesktop(imgPath, execPath, desktop)
	clean("/tmp")
}
