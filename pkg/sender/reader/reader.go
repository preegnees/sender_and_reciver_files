package reader

import (
	"context"
	"os"
	"io"
	"sync"
	"log"
	"bufio"
	"fmt"

	common "files/pkg/common"
)

func Read(ctx context.Context, bufferSize int64, optionsSize int64, fileSize int64, file *os.File, errCh chan<- error, writer io.Writer) {
	var storage = sync.Pool{New: func() interface{} {
		lines := make([]byte, bufferSize + optionsSize)
		return lines
	}}
	log.Println("Пулл с памятю инициализоварован")

	r := bufio.NewReader(file)
	var currentOffset int64
	const FIRST int16 = 1
	const LAST int16 = 2
	const OTHER int16 = 0
	var index int16
	for currentOffset = 0; currentOffset <= fileSize; currentOffset += bufferSize {
		// можно еще сюда добавить отмену пользователем
		
		select {
		case <-ctx.Done():
			errCh <- fmt.Errorf("$Был отменен контекст, read() завершило выполнение")
			return
		default:
			nSeek, err := file.Seek(currentOffset, 0)
			if err != nil {
				errCh <- fmt.Errorf("$Ошибка при сдвиге в файле=%s, err:=%v", file.Name(), err)
				return
			}
			log.Println("nSeek:", nSeek)

			poolStorage := storage.Get().([]byte)
			nBytes, err := r.Read(poolStorage[:len(poolStorage)-int(optionsSize)])
			log.Println("Сдвиг в файле=", file.Name(), " состовляет n:", nSeek)

			if nBytes == 0 {
				errCh <- fmt.Errorf("$Файл весит 0 bytes или был достигнут конец файла")
				return
			} else if err == io.EOF {
				errCh <- fmt.Errorf("$Конец файла=%s", file.Name())
				return
			} else if err != nil {
				errCh <- fmt.Errorf("$Ошибка при чтении файла=%s, err:=%v", file.Name(), err)
				return
			}
			log.Println("Прочиталось", nBytes, " bytes, с файла=", file.Name(),
				", при этом сдвиг составил=", currentOffset)

			buf := poolStorage[:nBytes]
			storage.Put(poolStorage)

			if currentOffset == 0 {
				index = FIRST
			} else if currentOffset + bufferSize > fileSize {
				index = LAST
			} else {
				index = OTHER
			}
			options := common.ConstructOptions(file.Name(), currentOffset, index)
			buf = append(buf, options...)
			writer.Write(buf)
			log.Println("Было отправлено в канал:", len(poolStorage), ", это:", currentOffset + int64(nBytes), "/", fileSize)
		}
	} 
}