package tests

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	assert "github.com/stretchr/testify/assert"

	models "files/pkg/models"
	readFile "files/pkg/sender/readFile"
)

func TestReadFile_TrueSettings(t *testing.T) {
	start := time.Now().UnixMilli()

	ctx, cancel := context.WithCancel(context.Background())
	pr, pw := io.Pipe()
	var settings models.SettingsReaderOfFile = models.SettingsReaderOfFile{
		Ctx:         ctx,
		PathToFile:  "C:\\Users\\secrr\\Desktop\\work_with_files\\example\\send.txt",
		ParentDir:   "C:\\Users\\secrr\\Desktop\\work_with_files\\example",
		BufferSize:  200 * 1024,
		OptionsSize: 200 * 1024,
		Writer:      pw,
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		path := "C:\\Users\\secrr\\Desktop\\work_with_files\\example\\test.txt"
		file, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)
		w := bufio.NewWriter(file)
		defer file.Close()
		for {
			select {
			default:
				if _, err := io.Copy(w, pr); err != nil {
					cancel()
					log.Println(err)
				}
				wg.Done()
				return
			}
		}
	}()
	err := readFile.ReadFile(settings)
	t.Log(err)
	pw.Close()
	wg.Wait()

	stop := time.Now().UnixMilli()
	log.Println("время выполнения:", stop-start)
}

func TestReadFile_FalseSettings_FalseConf(t *testing.T) {
	_, pw := io.Pipe()
	var settings models.SettingsReaderOfFile = models.SettingsReaderOfFile{
		Ctx:         context.Background(),
		PathToFile:  "C:\\Users\\secrr\\Desktop\\work_with_files\\example\\send.txt",
		ParentDir:   "C:\\Users\\secrr\\Desktop\\work_with_files\\example",
		BufferSize:  200 * 1024,
		OptionsSize: 200 * 1024,
		Writer:      pw,
	}
	err := readFile.ReadFile(settings)
	log.Println(err)
	assert.True(t, err != nil)
}
