/*
Package util util consists of general utility functions and structures.

Copyright 2021 Michael Bungenstock

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the
License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.
*/
package util

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

// LoadPrivateKeyPEM reads the private key from the passed file in PEM format.
func LoadPrivateKeyPEM(path string) (*rsa.PrivateKey, error) {
	var (
		dat []byte
		err error
	)
	//var pkey *rsa.PrivateKey
	if dat, err = ioutil.ReadFile(path); err != nil {
		return nil, fmt.Errorf("Could not load private key: %v", err)
	}
	block, _ := pem.Decode(dat)
	if block == nil {
		return nil, fmt.Errorf("Failed to decode PEM block in file %v", path)
	}
	var key interface{}
	key, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse private key in file %v: %v", path, err)
	}
	return key.(*rsa.PrivateKey), nil
}
