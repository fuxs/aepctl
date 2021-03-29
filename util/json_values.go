package util

import "io"

type JSONValueIterator struct {
	*JSONIterator
}

func NewJSONValueIterator(stream io.ReadCloser) *JSONValueIterator {
	return &JSONValueIterator{JSONIterator: NewJSONIterator(stream)}
}

func (j *JSONValueIterator) Next() (*Query, error) {
	return j.c.NextValue()
}
