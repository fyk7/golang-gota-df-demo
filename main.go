package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

func main() {
	sample0()
	sample1()
	sample2()
	sample3()
	sample5()
}

func sample0() {
	type User struct {
		Name     string
		Age      int
		Accuracy float64
		ignored  bool
	}

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
}

func sample1() {
	csvStr := `
Country,Date,Age,Amount,Id
"United States",2012-02-01,50,112.1,01234
"United States",2012-02-01,32,321.31,54320
"United Kingdom",2012-02-01,17,18.2,12345
"United States",2012-02-01,32,321.31,54320
"United Kingdom",2012-02-01,NA,18.2,12345
"United States",2012-02-05,32,321.31,54320
"United States",2012-02-05,32,321.31,54320
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

func sample2() {
	csvStr1 := `
Country,Date,Age,Amount,Id
"United States",2012-02-01,50,112.1,01234
"United States",2012-02-01,32,321.31,54320
"United Kingdom",2012-02-01,17,18.2,12345
"United States",2012-02-01,32,321.31,54320
"United Kingdom",2012-02-01,NA,18.2,12345
"United States",2012-02-05,32,321.31,54320
"United States",2012-02-05,32,321.31,54320
Spain,2012-02-01,66,555.42,00241
`
	csvStr2 := `
Country,Capital
"United States","ワシントンDC"
"United Kingdom","ロンドン"
"Spain","バルセロナ"
`
	df1 := dataframe.ReadCSV(strings.NewReader(csvStr1))
	df2 := dataframe.ReadCSV(strings.NewReader(csvStr2))

	join := df1.InnerJoin(df2, "Country")
	fmt.Println(join)

}

func sample3() {
	f, err := os.Open("sample.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	df := dataframe.ReadCSV(f)
	fmt.Println(df)
	fmt.Println(df.Types())
}

func sample5() {
	f, err := os.Open("sample.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	df := dataframe.ReadCSV(f)

	// 特定の日付よりも新しい行を選択するためのカスタム演算子
	isLater := func(date string) func(el series.Element) bool {
		return func(el series.Element) bool {
			if el.Type() == series.String {
				// 指定したレイアウトを満たせば、time.Time型に変換してくれる
				layout := "2006-01-02"
				thresh, err := time.Parse(layout, date)
				if err != nil {
					fmt.Println(err)
					return false
				}
				if val, ok := el.Val().(string); ok {
					d, err := time.Parse(layout, val)
					if err != nil {
						fmt.Println(err)
						return false
					}
					return d.After(thresh)
				}
			}
			return false
		}
	}

	df2 := df.Filter(
		dataframe.F{
			Colname:    "Date",
			Comparator: series.CompFunc,
			Comparando: isLater("2012-02-03"),
		},
	)
	fmt.Println(df2)

}
