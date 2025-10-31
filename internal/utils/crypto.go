/*
All the functions related to cryptographic functions are defined here
so they can be managed from a single place
*/
package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func HashDeprecated(params ...string) string {
	s := md5.New()
	result := []byte("")
	for _, param := range params {
		result = s.Sum([]byte(param))
	}
	return fmt.Sprintf("%x", result)
}

func Hash(params ...string) string {
	var payload string = ""
	for _, param := range params {
		payload += param
	} 
	hash := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(hash[:])
}
