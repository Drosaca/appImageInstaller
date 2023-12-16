package desktopFile

import (
	"bufio"
	"fmt"
	"os"
)

type DesktopFile struct {
	data             map[string]map[string]string
	selectedCategory string
	source           string
}

func (d DesktopFile) Category(name string) DesktopFile {
	d.selectedCategory = name
	return d
}

func (d DesktopFile) Get(name string) (string, error) {
	data, exists := d.data[d.selectedCategory]
	if !exists {
		return "", fmt.Errorf("category not found")
	}
	value, exists := data[name]
	if !exists {
		return "", fmt.Errorf("value not found")
	}
	d.selectedCategory = "root"
	return value, nil
}

func (d DesktopFile) Set(name string, value string) error {
	d.data[d.selectedCategory][name] = value
	d.selectedCategory = "root"
	return nil
}

func New() *DesktopFile {
	desktop := new(DesktopFile)
	desktop.data = make(map[string]map[string]string)
	desktop.selectedCategory = "root"
	desktop.source = "self generated"
	return desktop
}
func (d DesktopFile) FromMap(desktopMap map[string]map[string]string) {
	d.data = desktopMap
	d.source = "self-generated"
}

func (d DesktopFile) FromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("fail to open file: %s", err)
	}
	scanner := bufio.NewScanner(file)
	d.parseFile(scanner)
	d.source = path
	return nil
}

func (d DesktopFile) ToFile(path string) error {
	file, err := os.Create(path)

	if err != nil {
		return fmt.Errorf("fail to create desktop file: %s", err)
	}
	for category, parameter := range d.data {
		_, err := file.WriteString("[" + category + "]\n")
		if err != nil {
			return err
		}
		for name, value := range parameter {
			_, err := file.WriteString(name + "=" + value + "\n")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (d DesktopFile) HasValues(category string, values []string) bool {
	for _, value := range values {
		_, err := d.Category(category).Get(value)
		if err != nil {
			return false
		}
	}
	return true
}

func (d DesktopFile) GetSource() string {
	return d.source
}
