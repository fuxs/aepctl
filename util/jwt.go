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
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// Claim contains the three regisitered claims required for authentication.
// See https://tools.ietf.org/html/rfc7519 for more details.
type Claim struct {
	Iss string // iss (issuer):identifies the principal that issued the JWT
	Sub string // sub (subject): identifies the principal that is the subject of the JWT
	Aud string // aud (audience): identifies the recipients that the JWT is intended for
}

func (claim *Claim) header() string {
	header := `{"alg": "RS256", "typ": "JWT"}`
	return base64.RawStdEncoding.EncodeToString([]byte(header))
}

func (claim *Claim) payload() string {
	var str bytes.Buffer
	exp := time.Now().Add(time.Hour * 24).Unix() // exp (expiration): measuring the absolute time since 01/01/1970 GMT
	str.WriteString(`{"exp":`)
	str.WriteString(strconv.FormatInt(exp, 10))
	str.WriteString(`,"iss":"`)
	str.WriteString(claim.Iss)
	str.WriteString(`","sub":"`)
	str.WriteString(claim.Sub)
	str.WriteString(`","https://ims-na1.adobelogin.com/s/ent_dataservices_sdk": true`)
	str.WriteString(`,"aud":"`)
	str.WriteString(claim.Aud)
	str.WriteString(`"}`)
	return base64.RawStdEncoding.EncodeToString(str.Bytes())
}

// JWT calculates the signature of the claim and returns the JSON Web Token.
func (claim *Claim) JWT(key *rsa.PrivateKey) (string, error) {
	message := claim.header() + "." + claim.payload()
	hashed := sha256.Sum256([]byte(message))
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	signature, err := rsa.SignPKCS1v15(rng, key, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("Could not sign JWT: %v", err)
	}
	result := message + "." + base64.RawStdEncoding.EncodeToString(signature)
	return result, nil
}
