package models

import (
	"context"
	"io"
	"os"
)

type SettingsReaderOfFile struct {
	Ctx         context.Context
	PathToFile  string
	ParentDir   string
	BufferSize  int64
	OptionsSize int64
	Writer      io.Writer
}

type IReaderOfFile interface {
	ReadFile(SettingsReaderOfFile) error
}

type ConfForReader struct {
	Ctx context.Context
	BufferSize int64
	OptionsSize int64
	FileSize int64
	File *os.File
	RelativePath string
	ErrCh chan<- error
	Writer io.Writer
}
