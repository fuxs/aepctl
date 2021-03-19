package util

import (
	"bytes"
	"io"
	"io/ioutil"
)

type Restreamable struct {
	data []byte
}

func NewRestreamable(r io.Reader) (*Restreamable, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &Restreamable{data: data}, nil
}

func (r *Restreamable) Reader() io.Reader {
	return bytes.NewReader(r.data)
}
