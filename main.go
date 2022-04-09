package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/go-gota/gota/dataframe"
)

type User struct {
	Name     string
	Age      int
	Accuracy float64
	ignored  bool
}

func main() {
	users := []User{
		{"Aram", 17, 0.2, true},
		{"Juan", 18, 0.8, true},
		{"Ana", 22, 0.5, true},
	}

	df := dataframe.LoadStructs(users)

	var buf bytes.Buffer
	err := df.Filter(
		dataframe.F{
			Colname:    "Age",
			Comparator: ">=",
			Comparando: 20,
		},
	).Select(
		[]string{"Age", "Name"},
	).WriteJSON(&buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(buf.String())

	sample1()
}

func sample1() {
	csvStr := `
Country,Date,Age,Amount,Id
"United States",2012-02-01,50,112.1,01234
"United States",2012-02-01,32,321.31,54320
"United Kingdom",2012-02-01,17,18.2,12345
"United States",2012-02-01,32,321.31,54320
"United Kingdom",2012-02-01,NA,18.2,12345
"United States",2012-02-01,32,321.31,54320
"United States",2012-02-01,32,321.31,54320
Spain,2012-02-01,66,555.42,00241
`
	df := dataframe.ReadCSV(strings.NewReader(csvStr))

	var buf bytes.Buffer
	err := df.Select(
		// TODO "Age"にNaNが入っているがこちらをどのように処理すると良いかを考察する
		[]string{"Country", "Date", "Amount"},
	).GroupBy(
		"Country", "Date",
	).Aggregation(
		[]dataframe.AggregationType{dataframe.Aggregation_MEAN},
		[]string{"Amount"},
	).Select(
		[]string{"Country", "Amount_MEAN"},
	).WriteCSV(&buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(buf.String())

}
