package test

import (
	"log"
	"testing"

	assert "github.com/stretchr/testify/assert"

	workWithOptions "files/pkg/common"
)

func TestDeconstructOptions_NormalWork(t *testing.T) {
	name := "/file/me/test/dir"
	offset := 55555
	options := workWithOptions.ConstructOptions(name, int64(offset), 1, 256, 8)
	log.Println(options)
	offsetNew, nameNew, index, sizeOpt, sizeFile, bufferSize := workWithOptions.DeconstructOptions(options)
	assert.True(t, name == nameNew)
	assert.True(t, offset == int(offsetNew))
	assert.True(t, index == 1)
	log.Println("nameNew:", nameNew)
	log.Println("offsetNew:", offsetNew)
	log.Println("index:", index)
	log.Println("sizeOpt:", sizeOpt)
	log.Println("sizeFile:", sizeFile)
	log.Println("bufferSize:", bufferSize)
}