package common

import (
	"encoding/binary"
)

func ConstructOptions(name string, currentOffset int64, idx int16, sizeFile int64, buffer int64) (options []byte) {
	index := make([]byte, 2)
	binary.LittleEndian.PutUint16(index, uint16(idx))

	offset := make([]byte, 8)
	binary.LittleEndian.PutUint64(offset, uint64(currentOffset))

	size := make([]byte, 8)
	binary.LittleEndian.PutUint64(size, uint64(sizeFile))

	bufferSize := make([]byte, 8)
	binary.LittleEndian.PutUint64(bufferSize, uint64(buffer))

	sizeName := make([]byte, 8)
	binary.LittleEndian.PutUint64(sizeName, uint64(len(name)))

	fileName := make([]byte, len(name))[:0]
	fileName = append(fileName, []byte(name)...)

	option := make([]byte, 2+8+8+8+8+len(fileName))[:0]
	option = append(option, fileName...)
	option = append(option, sizeName[:]...)
	option = append(option, bufferSize[:]...)
	option = append(option, size[:]...)
	option = append(option, offset[:]...)
	option = append(option, index[:]...)
	return option
}

func DeconstructOptions(data []byte) (offset int64, name string, idx int16, sizeOpt int64, size int64, buffer int64) {
	indexData := data[len(data)-2:]
	index := int16(binary.LittleEndian.Uint16(indexData))

	currenOffsetData := data[len(data)-8-2 : len(data)-2]
	currentOffset := int64(binary.LittleEndian.Uint64(currenOffsetData))

	sizeFileData := data[len(data)-8-8-2 : len(data)-2-8]
	sizeFile := int64(binary.LittleEndian.Uint64(sizeFileData))

	bufferSizeData := data[len(data)-8-8-8-2 : len(data)-2-8-8]
	bufferSize := int64(binary.LittleEndian.Uint64(bufferSizeData))

	sizeFileNameData := data[len(data)-8-8-8-8-2 : len(data)-8-8-8-2]
	sizeFileName := int64(binary.LittleEndian.Uint64(sizeFileNameData))

	fileNameData := data[len(data)-8-8-8-8-2-int(sizeFileName): len(data)-8-8-8-8-2]
	fileName := string(fileNameData)

	sizeOtion := 2 + 8 + 8 + 8 + 8 + sizeFileName
	return currentOffset, fileName, index, sizeOtion, sizeFile, bufferSize
}
