package util

type JSONValueIterator struct {
	*JSONIterator
	filter []string
	f      bool
}

func NewJSONValueIterator(c *JSONCursor, filter []string) *JSONValueIterator {
	return &JSONValueIterator{JSONIterator: NewJSONIterator(c), filter: filter, f: len(filter) > 0}
}

func (j *JSONValueIterator) Next() (*Query, error) {
	if j.f {
		return j.c.NextValueF(j.filter)
	}
	return j.c.NextValue()
}

func (j *JSONValueIterator) More() bool {
	return j.c.MoreTokens()
}
