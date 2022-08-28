package readfile

import (
	"context"
	"fmt"
	"log"
	"os"


	models "files/pkg/models"
	reader "files/pkg/sender/reader"
)

type ReaderOfFile struct {
	currentOffset     int64
	numberOfDivisions int64
	bufferSize        int64
	end bool
}


func ReadFile(settings models.SettingsReaderOfFile) error {
	log.Print("Настройки: ", settings)

	conf, err := settings.IReaderOfConf.Get()
	if err != nil {
		return fmt.Errorf("$Ошибка при получении конфига, err:=%v", err)
	}
	log.Println("conf:", conf)

	file, err := os.Open(conf.Path)
	if err != nil {
		return fmt.Errorf("$Ошибка при открытии файла, path=%s, err:=%v", conf.Path, err)
	}
	defer file.Close()
	log.Println("Файл успешно открылся, file:", file)

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("$Ошикба при получении информации о файле, err:=%v", err)
	}
	log.Println("Размер файла состовляет, size:", stat.Size(), "bytes")

	numberOfDivisions := (stat.Size() / settings.BufferSize) + 1 // должнг быть проверено на тот случай, что buffer не кратно 8

	log.Println("Колличество делений файла состовляет, numberOfDivisions=", numberOfDivisions)

	ctxForRead, cansel := context.WithCancel(context.Background())
	errCh := make(chan error)
	go reader.Read(ctxForRead, settings.BufferSize, settings.OptionsSize, stat.Size(), file, errCh, settings.Writer)

	for {
		select {
		case <-settings.Ctx.Done():
			cansel()
			return fmt.Errorf("$Была вызвана отмена контекста, завершилось выполнение функции read()")
		case err:= <-errCh:
			cansel()
			return fmt.Errorf("$Ошибка при работы метода read(), err:=%v", err)
		}
	}
}
