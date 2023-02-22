package desktopFile

import (
	"bufio"
	"regexp"
	"strings"
)

func (d DesktopFile) parseCategory(line string) (bool, string) {
	rCategory, _ := regexp.Compile(`^( )*\[(.+)\]$`)
	if rCategory.MatchString(line) {
		category := rCategory.FindStringSubmatch(line)[2]
		d.data[category] = map[string]string{}
		return true, category
	}
	return false, ""
}
func (d DesktopFile) parseParameter(line string, category string) {
	rCategory, _ := regexp.Compile(`.+=.+`)
	if rCategory.MatchString(line) {
		tuple := strings.Split(line, "=")
		d.data[category][tuple[0]] = tuple[1]
	}
}

func (d DesktopFile) parseFile(scanner *bufio.Scanner) {
	category := "root"
	for scanner.Scan() {
		line := scanner.Text()
		isCategory, cat := d.parseCategory(line)
		if isCategory {
			category = cat
		}
		d.parseParameter(line, category)
	}
}
