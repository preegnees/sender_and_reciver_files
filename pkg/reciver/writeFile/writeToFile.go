package writefile

import (
	"fmt"
	"io"

	models "files/pkg/models"
	distribution "files/pkg/reciver/distribution"
	controller "files/pkg/reciver/writer"
)

func WriteFile(settings models.SettingsWriterToFile) error {
	// не забыть применить контекст
	cntrl := controller.New(settings.ParentDir)
	fd := distribution.New(cntrl)
	if _, err := io.Copy(fd, settings.Reader); err != nil {
		return fmt.Errorf("$Ошибка при чтнении, err:=%v", err)
	}
	return nil
}