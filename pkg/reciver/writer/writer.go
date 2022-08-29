package writer

import (
	"fmt"
	"log"
	"os"
	"strings"
	"path/filepath"
)

type ControllerOfFiles struct {
	storage   map[string](*os.File)
	ParentDir string
}

func New(parDir string) *ControllerOfFiles {

	return &ControllerOfFiles{
		storage: make(map[string]*os.File),
		ParentDir: parDir,
	} 
}


func (cof *ControllerOfFiles) WriteToFile(path string, data []byte) (nb int, err error) {
	file, ok := cof.storage[path]

	if !ok {

		pathArr := strings.Split(path, string(filepath.Separator))	
		dir := strings.Join(pathArr[:len(pathArr)-1], string(filepath.Separator))
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return 0, fmt.Errorf("$Ошибка при создании каталога:%s, err:=%v", dir, err)
		}

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			return 0, fmt.Errorf("$Ошибка при открытии файла, которого не было в хранилище, err:=%v", err)
		}
		cof.storage[path] = f
		file = f
		log.Println("Файл был открыт и закжширован")
	} else {

		log.Printf("Файл по пути:%s, был взят из кэша", path)
	}
	log.Printf("Данные который приходят:%s\n", data)
	
	n, err := file.Write(data)
	if err != nil {
		return 0, fmt.Errorf("$Ошибка при записи в файл=%s, err:=%v", path, err)
	}

	return n, nil
}

func (cof *ControllerOfFiles) CloseAllFiles() {

	for _, file := range cof.storage {
		file.Close()
	}
}

func (cof *ControllerOfFiles) CloseFile(path string) error {
	file, ok := cof.storage[path]

	if ok {
		file.Close()
		return nil
	} else {
		return fmt.Errorf("$Ошибка при закрытии файла=%s, видимо, такого файла нет", path)
	}
}