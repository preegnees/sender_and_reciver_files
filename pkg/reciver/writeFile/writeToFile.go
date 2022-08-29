package writefile

import (
	"io"
	"log"

	distribution "files/pkg/reciver/distribution"
	controller "files/pkg/reciver/writer"
	models "files/pkg/models"
)

func WriteFile(settings models.SettingsWriterToFile) {
	// dir := "C:\\Users\\secrr\\Desktop\\work_with_files\\example2"
	cntrl := controller.New(settings.ParentDir)
	fd := distribution.New(cntrl)
	for {
		select {
		default:
			if _, err := io.Copy(fd, settings.Reader); err != nil {
				log.Println(err)
				return
			}
		}
	}
}