/*
All the functions related to cryptographic functions are defined here
so they can be managed from a single place
*/
package utils

import (
	"crypto/md5"
	"fmt"
)

func Hash(params ...string) string {
	s := md5.New()
	result := []byte("")
	for _, param := range params {
		result = s.Sum([]byte(param))
	}
	return fmt.Sprintf("%x", result)
}