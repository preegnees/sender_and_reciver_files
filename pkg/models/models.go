package models

import (
	"context"
	"io"
) 

type SettingsReaderOfFile struct {
	context.Context
	IReaderOfConf
	SizeBuffer int64
	Writer io.Writer
}

type Conf struct {
	Path string
}

type IReaderOfConf interface {
	Get() (*Conf, error)
}

type IReaderOfFile interface {
	ReadFile(SettingsReaderOfFile) error
}