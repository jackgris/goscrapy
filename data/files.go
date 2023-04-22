package data

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jackgris/goscrapy/database"
)

func ReadCSV(path string) []database.Product {

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(file)

	var products []database.Product

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if len(record) > 1 {
			number := record[4]
			number = strings.ReplaceAll(number, ",", ".")
			price, err := strconv.ParseFloat(number, 64)
			if err != nil {
				log.Println("ReadCSV: ", err)
			} else {
				prices := []database.Value{}
				value := database.Value{
					Price: price,
				}
				prices = append(prices, value)

				product := database.Product{
					Name:  record[0],
					Price: prices,
				}
				products = append(products, product)
			}
		}
	}
	return products
}
