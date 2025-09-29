package utils

import (
	"bytes"

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
