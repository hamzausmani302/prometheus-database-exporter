package utils

import (
	"bytes"
	"fmt"
	"os"

	"github.com/go-gota/gota/dataframe"
)

func DataFrameToCSVBytes(df dataframe.DataFrame) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := df.WriteCSV(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DataFrameFromCSVBytes(data []byte) dataframe.DataFrame {
	return dataframe.ReadCSV(bytes.NewReader(data))
}


func SetEnvironmentVariable(key string, value string) { 
	err := os.Setenv(key, value)
	if err != nil {
		fmt.Printf("Error setting environment variable: %v\n", err)
		panic(err)
	}

}