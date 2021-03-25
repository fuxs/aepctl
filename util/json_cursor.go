package util

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type jsonState int

const (
	JSONS_UNDEFINED jsonState = iota
	JSONS_OPEN
	JSONS_DONE
	JSONS_O  // object attribute, expecting string or }
	JSONS_OV // object value, expecting {, [ VALUE or }
	JSONS_A  // array value, expecting {, [ VALUE or ]
)

type jsonStateStack []jsonState

func (jss *jsonStateStack) Push(js jsonState) {
	*jss = append(*jss, js)
}

func (jss *jsonStateStack) Peek() jsonState {
	l := len(*jss)
	if l == 0 {
		return JSONS_UNDEFINED
	}
	return (*jss)[l-1]
}

func (jss *jsonStateStack) Pop() jsonState {
	l := len(*jss)
	if l == 0 {
		return JSONS_UNDEFINED
	}
	result := (*jss)[l-1]
	*jss = (*jss)[:l-1]
	return result
}

type jsonPath []string

func (ps *jsonPath) Push(p string) {
	*ps = append(*ps, p)
}

func (ps *jsonPath) Peek() string {
	l := len(*ps)
	if l == 0 {
		return ""
	}
	result := (*ps)[l-1]
	*ps = (*ps)[:l-1]
	return result
}

func (ps *jsonPath) Pop() string {
	l := len(*ps)
	if l == 0 {
		return ""
	}
	result := (*ps)[l-1]
	*ps = (*ps)[:l-1]
	return result
}

type JSONCursor struct {
	dec    *json.Decoder
	stream io.ReadCloser
	jss    jsonStateStack
	jp     jsonPath
}

func NewJSONCursor(stream io.ReadCloser) *JSONCursor {
	jss := make(jsonStateStack, 0, 16)
	jss.Push(JSONS_OPEN)
	jp := make(jsonPath, 0, 16)
	return &JSONCursor{dec: json.NewDecoder(stream), stream: stream, jss: jss, jp: jp}
}

func (j *JSONCursor) PathInfo() (string, string) {
	l := len(j.jp)
	if l == 0 {
		return "", ""
	}
	name := j.jp[l-1]
	if l == 1 {
		return name, ""
	}
	return name, Concat(j.jp[:l-2], ".")
}

func (j *JSONCursor) More() bool {
	return j.dec.More()
}

func (j *JSONCursor) Offset() int64 {
	return j.dec.InputOffset()
}

func (j *JSONCursor) Token() (json.Token, error) {
	state := j.jss.Peek()
	if state == JSONS_DONE {
		return nil, io.EOF
	}
	t, err := j.dec.Token()
	if err != nil {
		return nil, err
	}
	switch state {
	case JSONS_OPEN:
		d, ok := t.(json.Delim)
		if !ok || !(d == '{' || d == '[') {
			return nil, fmt.Errorf("expecting [ or { at position %v", j.dec.InputOffset())
		}
		j.jss.Pop()
		j.jss.Push(JSONS_DONE)
		if d == '{' {
			j.jss.Push(JSONS_O)
		} else {
			j.jss.Push(JSONS_A)
		}
	case JSONS_O:
		str, ok := t.(string)
		if ok {
			j.jp.Push(str)
			j.jss.Push(JSONS_OV)
			return t, nil
		}
		d, ok := t.(json.Delim)
		if !ok || !(d == '}') {
			return nil, fmt.Errorf("expecting } position %v", j.dec.InputOffset())
		}
		j.jss.Pop()
		if j.jss.Peek() == JSONS_O {
			j.jp.Pop()
		}
	case JSONS_OV:
		d, ok := t.(json.Delim)
		if ok {
			switch d {
			case '{':
				j.jss.Push(JSONS_O)
			case '[':
				j.jss.Push(JSONS_A)
			case '}':
				j.jss.Pop()
				if j.jss.Peek() == JSONS_O {
					j.jp.Pop()
				}
			default:
				return nil, fmt.Errorf("expecting [,{ or } at position %v", j.dec.InputOffset())
			}
		}
	case JSONS_A:
		d, ok := t.(json.Delim)
		if ok {
			switch d {
			case '{':
				j.jss.Push(JSONS_O)
			case '[':
				j.jss.Push(JSONS_A)
			case ']':
				j.jss.Pop()
			default:
				return nil, fmt.Errorf("expecting [,{ or ] at position %v", j.dec.InputOffset())
			}
		}
	default:
		return nil, errors.New("state error")
	}
	return t, nil
}

func (j *JSONCursor) Decode(v interface{}) error {
	state := j.jss.Peek()
	if state == JSONS_DONE {
		return io.EOF
	}
	if state != JSONS_OV && state != JSONS_A && state != JSONS_OPEN {
		return errors.New("state error")
	}
	if err := j.dec.Decode(v); err != nil {
		return err
	}
	j.jss.Pop()
	if j.jss.Peek() == JSONS_O {
		j.jp.Pop()
	}
	return nil
}

// Close closes the underlying ReaderCloser stream
func (j *JSONCursor) Close() error {
	return j.stream.Close()
}

// PrintRaw copies the raw data to standard out
func (j *JSONCursor) PrintRaw() error {
	bout := bufio.NewWriter(os.Stdout)
	defer bout.Flush()
	_, err := io.Copy(bout, j.stream)
	return err
}

// PrintPretty prints the raw data with indention to standard out
func (j *JSONCursor) PrintPretty() error {
	bout := bufio.NewWriter(os.Stdout)
	defer bout.Flush()
	return JSONPrintPretty(j.dec, bout)
}
