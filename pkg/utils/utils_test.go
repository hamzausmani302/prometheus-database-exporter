package utils

import (
	"fmt"
	"testing"

	"github.com/go-gota/gota/dataframe"
)

func TestDataFrameConversionFunction(t *testing.T) {
	// create a sample dataframe
	df := dataframe.LoadRecords(
			[][]string{
				{"Name", "Age", "City"},
				{"Alice", "30", "New York"},
				{"Bob", "25", "Los Angeles"},
				{"Charlie", "35", "Chicago"},
			},
		)
	bData, err := DataFrameToCSVBytes(df)
	if err != nil {
		fmt.Println("Conversion to bytes failed")
	}
	df2 := DataFrameFromCSVBytes(bData)

	if df.Nrow() != df2.Nrow() || df.Ncol() != df2.Ncol() {
		t.Errorf("DataFrame conversion failed, expected %d rows and %d cols, got %d rows and %d cols", df.Nrow(), df.Ncol(), df2.Nrow(), df2.Ncol())
	}
	//chck first value
	if df2.Elem(0, 0).String() != "Alice"{
		t.Errorf("DataFrame conversion failed, expected first row second column to be 'Alice', got %s", df2.Elem(0, 1))
	}
	//chec single value
	if df2.Elem(0, 1).String() != "30"{
		t.Errorf("DataFrame conversion failed, expected first row second column to be 'Alice', got %s", df2.Elem(0, 1))
	}
	// check data frames are equal
	for i :=0;i <df.Nrow();i++{
		for j:=0;j<df.Ncol();j++{
			if df.Elem(i,j).String() != df2.Elem(i,j).String(){
				t.Errorf("DataFrame conversion failed at row %d and column %d, expected %s, got %s", i, j, df.Elem(i,j).String(), df2.Elem(i,j).String())
			}
		}
	}
	fmt.Println("DataFrame conversion test passed")
}