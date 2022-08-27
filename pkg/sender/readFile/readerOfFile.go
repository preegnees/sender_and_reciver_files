package readfile

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"io"

	models "files/pkg/models"
)

type ReaderOfFile struct {
	currentOffset     int64
	numberOfDivisions int64
	bufferSize        int64
	end bool
}

var maxOfDivisions int64 

var _ models.IReaderOfFile = (*ReaderOfFile)(nil)

func (rof *ReaderOfFile) ReadFile(settings models.SettingsReaderOfFile) error {
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

	numberOfDivisions := (stat.Size() / settings.SizeBuffer) + 1 // должнг быть проверено на тот случай, что buffer не кратно 8
	rof.setNumberOfDivisions(numberOfDivisions)
	rof.setBufferSize(settings.SizeBuffer)
	rof.setCurrentOffset()
	log.Println("Колличество делений файла состовляет, numberOfDivisions=", numberOfDivisions)

	//
	maxOfDivisions = numberOfDivisions
	//
	ctxForRead, cansel := context.WithCancel(context.Background())
	errCh := make(chan error)

	go rof.read(ctxForRead, file, errCh, settings.Writer)

	for {
		select {
		case <-settings.Context.Done():
			cansel()
			return fmt.Errorf("$Была вызвана отмена контекста, завершилось выполнение функции read()")
		case err:= <-errCh:
			cansel()
			return fmt.Errorf("$Ошибка при работы метода read(), err:=%v", err)
		}
	}
}

func (rof *ReaderOfFile) read(ctx context.Context, file *os.File, errCh chan<- error, writer io.Writer) {
	var storage = sync.Pool{New: func() interface{} {
		lines := make([]byte, rof.bufferSize + 1024) // учитывать 1024 при чтении конфига, почтиать про r w одновремннных без буфкра
		return lines
	}}
	log.Println("Пулл с памятю инициализоварован")

	r := bufio.NewReader(file)
	for {
		select {
		case <-ctx.Done():
			errCh <- fmt.Errorf("$Был отменен контекст, read() завершило выполнение")
			return
		default:
			offset := rof.getCurrentOffset()
			if rof.end == true {
				errCh <- fmt.Errorf("$Конец файла=%s", file.Name())
				return
			}
			nSeek, err := file.Seek(offset, 0)
			if err != nil {
				errCh <- fmt.Errorf("$Ошибка при сдвиге в файле=%s, err:=%v", file.Name(), err)
				return
			}
			log.Println("Сдвиг в файле=", file.Name(), " состовляет n:", nSeek)

			poolStorage := storage.Get().([]byte)
			nBytes, err := r.Read(poolStorage[:])
			if nBytes == 0 {
				errCh <- fmt.Errorf("$Файл весит 0 bytes")
				return
			} else if err == io.EOF {
				errCh <- fmt.Errorf("$Конец файла=%s", file.Name())
				return
			} else if err != nil {
				errCh <- fmt.Errorf("$Ошибка при чтении файла=%s, err:=%v", file.Name(), err)
				return
			}
			log.Println("Прочиталось", nBytes, " bytes, с файла=", file.Name(),
				", при этом сдвиг составил=", rof.numberOfDivisions*1024, ", это индекс=", rof.numberOfDivisions)

			buf := poolStorage[:nBytes]
			storage.Put(poolStorage)
			options := rof.constructOptions(file.Name())
			buf = append(buf, options...)
			writer.Write(buf)
			log.Println("Было отправлено в канал:", len(buf), ", это:", rof.numberOfDivisions, "/", maxOfDivisions)

			rof.decrNumberOfDivisions()
		}
	}
}

func (rof *ReaderOfFile) constructOptions(name string) []byte {
	indexArr := make([]byte, 8)
	index := append(indexArr, byte(rof.numberOfDivisions))

	sizeNameArr := make([]byte, 8)
	lenFileName := len(name)
	sizeName := append(sizeNameArr, byte(lenFileName))

	fileNameArr := make([]byte, lenFileName)
	fileName := append(fileNameArr, []byte(name)...)

	option := make([]byte, 8+8+lenFileName)
	option = append(option, fileName...)
	option = append(option, sizeName...)
	option = append(option, index...)
	return option
}

func (rof *ReaderOfFile) getCurrentOffset() int64 {
	if rof.numberOfDivisions == 0 {
		rof.end = true
	}
	rof.currentOffset = rof.currentOffset * (1 + rof.bufferSize)
	return rof.currentOffset
}

func (rof *ReaderOfFile) setNumberOfDivisions(number int64) {
	rof.numberOfDivisions = number
}

func (rof *ReaderOfFile) decrNumberOfDivisions() {
	rof.numberOfDivisions = rof.numberOfDivisions - 1
}

func (rof *ReaderOfFile) setBufferSize(bufferSize int64) {
	rof.bufferSize = bufferSize
}

func (rof *ReaderOfFile) setCurrentOffset() {
	rof.currentOffset = 0
}
