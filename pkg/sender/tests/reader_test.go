package tests

import (
	"bufio"
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
)

func TestRead_TrueSettings(t *testing.T) {
	start := time.Now().UnixMilli()

	ctx, cancel := context.WithCancel(context.Background())
	pr, pw := io.Pipe()
	var wg sync.WaitGroup
	wg.Add(1)

	errCh := make(chan error)
	// path := "D:\\downloaded\\LibreOffice\\LibreOffice_7.3.5_Win_x64.msi"
	path := "C:\\Users\\secrr\\Desktop\\work_with_files\\example\\send.txt"
	file, _ := os.Open(path)
	stat, _ := file.Stat()
	size := stat.Size()

	go func() {
		// нужно так же написать шутку которая бы возварщая нужные открытые файлы, а потом закрывала
		// path := "C:\\Users\\secrr\\Desktop\\work_with_files\\example\\recive.txt"
		// file, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)
		// w := bufio.NewWriter(file)
		defer file.Close()
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

	reader.Read(ctx, 8, 256, size, file, errCh, pw) // 100 * 1024 * 1024 == 100 mb
	pw.Close()
	wg.Wait()

	stop := time.Now().UnixMilli()
	log.Println("время выполнения:", stop-start)
}

type controllerOfFiles struct {
	storage map[string](*os.File)
}

func (fc *controllerOfFiles) writeFile(path string, data []byte) (nb int, err error) {
	file, ok := fc.storage[path]
	if !ok {
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			return 0, fmt.Errorf("$Ошибка при открытии файла (получение файла из хранилища), err:=%v", err)
		}
		fc.storage[path] = file
	}
	w := bufio.NewWriter(file)
	n, err := w.Write(data)
	if err != nil {
		return n, fmt.Errorf("$Ошибка при записи в файл=%s, err:=%v", path, err)
	}

	return n, nil
}

func (fc *controllerOfFiles) closeAllFiles() {
	for _, v := range fc.storage {
		v.Close()
	}
}

func (fc *controllerOfFiles) closeFile(path string) error {
	file, ok := fc.storage[path]
	if ok {
		file.Close()
		return nil
	} else {
		return fmt.Errorf("$Ошибка при закрытии файла=%s", path)
	}
}

type Reader interface {
	Write(buf []byte) (n int, err error)
}
type FileDistributor struct{}

var _ io.Writer = (*FileDistributor)(nil)

var fc controllerOfFiles = controllerOfFiles{
	storage: make(map[string]*os.File),
}

func (r *FileDistributor) Write(buf []byte) (n int, err error) {
	
	fmt.Printf("%s\n", buf)
	offset, path, index, sizeOpt := common.DeconstructOptions(buf)
	log.Println("offset:", offset)
	log.Println("path:", path)
	log.Println("index:", index)
	log.Println("sizeOpt:", sizeOpt)
	log.Printf("message: %s\n", buf[:len(buf)-int(sizeOpt)])
	n, err = fc.writeFile(path, buf[:len(buf)-int(sizeOpt)])
	log.Println(err)
	if err != nil {
		return n, err
	}
	return len(buf), nil
}
