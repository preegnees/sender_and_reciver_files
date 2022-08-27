package common

import (
	models "files/pkg/models"
)

type Conf struct {
	models.Conf
}

var _ models.IReaderOfConf = (*Conf)(nil)

func (c *Conf) Get() (*models.Conf, error) {
	return nil, nil
}