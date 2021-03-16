package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// JSONFile stores json objects in a file
type JSONFile struct {
	p Path
}

// NewJSONFile returns an initialzed JSONCache object
func NewJSONFile(path Path) *JSONFile {
	return &JSONFile{p: path}
}

// Delete deletes the related file
func (jc *JSONFile) Delete() error {
	return os.Remove(jc.p.Path())
}

// Save stores the passed object in json format to a file
func (jc *JSONFile) Save(obj interface{}) error {
	p := filepath.Dir(jc.p.Path())
	if _, err := os.Stat(p); os.IsNotExist(err) {
		if err = os.MkdirAll(p, 0700); err != nil {
			return err
		}
	}
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(jc.p.Path(), data, 0700)
}

// Exists returns true if the related file exists
func (jc *JSONFile) Exists() bool {
	_, err := os.Stat(jc.p.Path())
	return err != nil
}

// Load loads the json file into the passed object
func (jc *JSONFile) Load(obj interface{}) error {
	var data []byte
	data, err := ioutil.ReadFile(jc.p.Path())
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &obj)
}
