package functions

import (
	"appinstall/desktopFile"
	"appinstall/structs"
	"appinstall/utils"
	"fmt"
)

func CopyImage(desktop *desktopFile.DesktopFile, config structs.Config) error {
	types := []string{".svg", ".png", ".jpg", ".jpeg", ".bmp", ".ico", ".cur", ".tif", ".tiff", ".webp", ".jfif", ".pjpeg", ".pjp", ".gif", ".avif", ".apng"}
	imgName, err := desktop.Category("Desktop Entry").Get("Icon")
	if err != nil {
		fmt.Println("no image name in desktop file: ", err)
	}
	for i := 0; i < len(types); i++ {
		types[i] = imgName + types[i]
	}
	imgPath, err := utils.FindFile(config.AppExtractDir, types)
	if err != nil {
		return fmt.Errorf("fail to find image: %s", err)
	}
	err = utils.Copy(imgPath, config.ImgPath)
	if err != nil {
		return fmt.Errorf("fail to copy image: %s", err)
	}
	return nil
}
