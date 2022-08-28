package test

import (
	"log"
	"testing"

	common "files/pkg/common"
)

func TestGetRelativePath(t *testing.T) {
	paretDir := "C:\\Users\\secrr\\Desktop\\work_with_files\\example"
	pathToFile := "C:\\Users\\secrr\\Desktop\\work_with_files\\example\\send.txt"
	relP := common.GetRelativePath(paretDir, pathToFile)

	log.Println(relP)
}