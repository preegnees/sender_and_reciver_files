package reader

import (
	"io"
	"sync"
	"log"
	"bufio"
	"fmt"

	common "files/pkg/common"
	models "files/pkg/models"
)

func Read(cnf models.ConfForReader) {

	var storage = sync.Pool{New: func() interface{} {
		lines := make([]byte, cnf.BufferSize + cnf.OptionsSize)
		return lines
	}}
	log.Println("Пулл с памятю инициализоварован")

	defer cnf.File.Close()

	r := bufio.NewReader(cnf.File)

	var currentOffset int64
	const FIRST int16 = 1
	const LAST int16 = 2
	const OTHER int16 = 0
	var index int16

	for currentOffset = 0; currentOffset <= cnf.FileSize; currentOffset += cnf.BufferSize {		
		select {
		case <-cnf.Ctx.Done():
			cnf.ErrCh <- fmt.Errorf("$Был отменен контекст, read() завершило выполнение")
			return
		default:
			nSeek, err := cnf.File.Seek(currentOffset, 0)
			if err != nil {
				cnf.ErrCh <- fmt.Errorf("$Ошибка при сдвиге в файле=%s, err:=%v", cnf.File.Name(), err)
				return
			}
			log.Println("nSeek:", nSeek)

			poolStorage := storage.Get().([]byte)
			nBytes, err := r.Read(poolStorage[:len(poolStorage)-int(cnf.OptionsSize)])
			
			log.Println("Сдвиг в файле=", cnf.File.Name(), " состовляет n:", nSeek)

			if nBytes == 0 {
				cnf.ErrCh <- fmt.Errorf("$Файл весит 0 bytes или был достигнут конец файла")
				return
			} else if err == io.EOF {
				cnf.ErrCh <- fmt.Errorf("$Конец файла=%s", cnf.File.Name())
				return
			} else if err != nil {
				cnf.ErrCh <- fmt.Errorf("$Ошибка при чтении файла=%s, err:=%v", cnf.File.Name(), err)
				return
			}
			log.Println("Прочиталось", nBytes, " bytes, с файла=", cnf.File.Name(),
				", при этом сдвиг составил=", currentOffset)

			buf := poolStorage[:nBytes]
			storage.Put(poolStorage)

			if currentOffset == 0 {
				index = FIRST
			} else if currentOffset + cnf.BufferSize > cnf.FileSize {
				index = LAST
			} else {
				index = OTHER
			}
			options := common.ConstructOptions(cnf.File.Name(), currentOffset, index, cnf.FileSize, cnf.BufferSize)
			buf = append(buf, options...)
			cnf.Writer.Write(buf)
			log.Println("Было отправлено в канал:", len(poolStorage), ", это:", currentOffset + int64(nBytes), "/", cnf.FileSize)
		}
	} 
}