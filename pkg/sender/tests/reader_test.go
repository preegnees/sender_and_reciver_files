package tests

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	common "files/pkg/common"
	reader "files/pkg/sender/reader"
	models "files/pkg/models"
)

func TestRead_TrueSettings(t *testing.T) {
	start := time.Now().UnixMilli()

	ctx, cancel := context.WithCancel(context.Background())
	pr, pw := io.Pipe()
	var wg sync.WaitGroup
	wg.Add(1)

	errCh := make(chan error)
	path := "C:\\Users\\secrr\\Desktop\\work_with_files\\example\\send.txt"
	file, _ := os.Open(path)
	stat, _ := file.Stat()
	size := stat.Size()

	go func() {
		for {
			select {
			case err := <-errCh:
				log.Println(err)
			default:
				fd := FileDistributor{}
				if _, err := io.Copy(&fd, pr); err != nil {
					cancel()
					log.Println(err)
				}
				wg.Done()
				return
			}
		}
	}()
	paretDir := "C:\\Users\\secrr\\Desktop\\work_with_files\\example"
	pathToFile := "C:\\Users\\secrr\\Desktop\\work_with_files\\example\\send.txt"
	relativePath := common.GetRelativePath(paretDir, pathToFile)
	cnf := models.ConfForReader{
		Ctx: ctx,
		BufferSize: 256,
		OptionsSize: 256,
		FileSize: size,
		File: file,
		RelativePath: relativePath,
		ErrCh: errCh,
		Writer: pw,
	}
	reader.Read(cnf) // 100 * 1024 * 1024 == 100 mb
	pw.Close()
	wg.Wait()
	pr.Close()

	stop := time.Now().UnixMilli()
	log.Println("время выполнения:", stop-start)
}

type controllerOfFiles struct {
	storage   map[string](*os.File)
	parentDir string
}

func (cof *controllerOfFiles) writeToFile(path string, data []byte) (nb int, err error) {
	file, ok := cof.storage[path]
	if !ok {
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			return 0, fmt.Errorf("$Ошибка при открытии файла (файла не было в хронилище), err:=%v", err)
		}
		cof.storage[path] = f
		file = f

	}
	
	n, err := file.Write(data)
	if err != nil {
		return 0, fmt.Errorf("$Ошибка при записи в файл=%s, err:=%v", path, err)
	}

	return n, nil
}

func (cof *controllerOfFiles) closeAllFiles() {
	for _, file := range cof.storage {
		file.Close()
	}
}

func (cof *controllerOfFiles) closeFile(path string) error {
	file, ok := cof.storage[path]
	if ok {
		file.Close()
		return nil
	} else {
		return fmt.Errorf("$Ошибка при закрытии файла=%s, видимо, такого файла нет", path)
	}
}

func (cof *controllerOfFiles) New() *controllerOfFiles {
	return &controllerOfFiles{
		storage: make(map[string]*os.File),
	} 
}

type Reader interface {
	Write(buf []byte) (n int, err error)
}
type FileDistributor struct{}

var cof controllerOfFiles = controllerOfFiles{}
var controller = cof.New() 

var _ io.Writer = (*FileDistributor)(nil)

func (r *FileDistributor) Write(buf []byte) (n int, err error) {
	fmt.Printf("%s\n", buf)
	offset, path, index, sizeOpt, sizeFile, sizeBuffer := common.DeconstructOptions(buf)
	log.Println("offset:", offset)
	log.Println("path:", path)
	log.Println("index:", index)
	log.Println("sizeOpt:", sizeOpt)
	log.Println("sizeFile:", sizeFile)
	log.Println("sizeBuffer:", sizeBuffer)
	log.Printf("message: %s\n", buf[:len(buf)-int(sizeOpt)])
	n, err = controller.writeToFile(path, buf[:len(buf)-int(sizeOpt)])
	log.Println(err)
	if err != nil {
		return n, err
	}
	return len(buf), nil
}