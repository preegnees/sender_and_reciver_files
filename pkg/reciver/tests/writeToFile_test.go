package test

import (
	"context"
	"log"
	"io"
	"sync"
	"testing"

	models "files/pkg/models"
	wf "files/pkg/reciver/writeFile"
	rf "files/pkg/sender/readFile"
)

func TestWriteToFile(t *testing.T) {
	ctx := context.Background()
	pr, pw := io.Pipe()

	settingForReciver := models.SettingsWriterToFile{
		Ctx:       ctx,
		ParentDir: "C:\\Users\\secrr\\Desktop\\work_with_files\\example2",
		Reader:    pr,
	}

	settingForSender := models.SettingsReaderOfFile{
		Ctx:           ctx,
		PathToFile:    "C:\\Users\\secrr\\Desktop\\work_with_files\\example\\me\\блин как жить.txt",
		ParentDir:     "C:\\Users\\secrr\\Desktop\\work_with_files\\example",
		BufferSize:    5,
		OptionsSize:   256,
		Writer:        pw,
		CurrentOffset: 0,
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		err := wf.WriteFile(settingForReciver)
		if err != nil {
			log.Println(err)
		}
		pw.Close()
		wg.Done()
	}()
	go func() {
		err := rf.ReadFile(settingForSender)
		if err != nil {
			log.Println("Err from ReadFile:", err)
		}
		pr.Close()
		wg.Done()
	}()
	wg.Wait()	
}
