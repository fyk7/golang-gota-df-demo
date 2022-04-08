package main

import (
	"bytes"
	"fmt"

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
	df.Filter(
		dataframe.F{
			Colname:    "Age",
			Comparator: ">=",
			Comparando: 20,
		},
	).WriteJSON(&buf)
	// err := df.Select([]string{"Age", "Name"}).WriteJSON(&buf)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println(buf.String())

}
