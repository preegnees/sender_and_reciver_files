package tests

import (
	"context"
	"io"
	"log"
	"os"
	"sync"
	"testing"
	"time"
	"bufio"

	assert "github.com/stretchr/testify/assert"

	models "files/pkg/models"
	readFile "files/pkg/sender/readFile"
)

type readerOfConfTrue struct{}

func (roc *readerOfConfTrue) Get() (*models.Conf, error) {
	var conf models.Conf = models.Conf{
		// Path: "C:\\Users\\secrr\\Desktop\\work_with_files\\example\\send.txt",
		Path: "D:\\downloaded\\LibreOffice\\LibreOffice_7.3.5_Win_x64.msi",
	}
	return &conf, nil
}

type readerOfConfFalse struct{}

func (roc *readerOfConfFalse) Get() (*models.Conf, error) {
	var conf models.Conf = models.Conf{
		Path: "C:\\Users\\secrr\\Desktop\\work_with_files\\example\\se.txt",
	}
	return &conf, nil
}

func TestReadFile_TrueSettings(t *testing.T) {
	start := time.Now().UnixMilli()

	ctx, cancel := context.WithCancel(context.Background())
	pr, pw := io.Pipe()
	var settings models.SettingsReaderOfFile = models.SettingsReaderOfFile{
		Ctx:       ctx,
		IReaderOfConf: &readerOfConfTrue{},
		BufferSize:    100 * 1024 * 1024, // 200 * 1024
		Writer:        pw,
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		path := "C:\\Users\\secrr\\Desktop\\work_with_files\\example\\libre"
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
		Ctx:       context.Background(),
		IReaderOfConf: &readerOfConfFalse{},
		BufferSize:    8 * 1024,
		Writer:      pw,
	}
	err := readFile.ReadFile(settings)
	log.Println(err)
	assert.True(t, err != nil)
}

func TestReadFile_experiment(t *testing.T) {
	log.Println("helllo world")
}
