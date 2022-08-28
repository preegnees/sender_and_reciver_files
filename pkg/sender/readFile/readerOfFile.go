package readfile

import (
	"context"
	"fmt"
	"log"
	"os"

	models "files/pkg/models"
	reader "files/pkg/sender/reader"
	common "files/pkg/common"
)

type ReaderOfFile struct {
	currentOffset     int64
	numberOfDivisions int64
	bufferSize        int64
	end bool
}


func ReadFile(settings models.SettingsReaderOfFile) error {
	log.Print("Настройки: ", settings)

	file, err := os.Open(settings.PathToFile)
	if err != nil {
		return fmt.Errorf("$Ошибка при открытии файла, path=%s, err:=%v", settings.PathToFile, err)
	}
	defer file.Close()
	log.Println("Файл успешно открылся, file:", file)

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("$Ошикба при получении информации о файле, err:=%v", err)
	}
	log.Println("Размер файла состовляет, size:", stat.Size(), "bytes")

	ctxForRead, cancel := context.WithCancel(context.Background())
	errCh := make(chan error)
	relativePath := common.GetRelativePath(settings.ParentDir, settings.PathToFile)
	cnf := models.ConfForReader{
		Ctx: ctxForRead,
		BufferSize: settings.BufferSize,
		OptionsSize: settings.OptionsSize,
		FileSize: stat.Size(),
		File: file,
		RelativePath: relativePath,
		ErrCh: errCh,
		Writer: settings.Writer,
	}
	go reader.Read(cnf)

	select {
	case <-settings.Ctx.Done():
		cancel()
		return fmt.Errorf("$Была вызвана отмена контекста, завершилось выполнение функции read()")
	case err:= <-errCh:
		cancel()
		return fmt.Errorf("$Ошибка при работы метода read(), err:=%v", err)
	}
}
