package data

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/jackgris/goscrapy/database"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

func ReadCSV(path string, log *logrus.Logger) []database.Product {

	var products []database.Product

	file, err := os.Open(path)
	if err != nil {
		log.Println("ReadCSV: ", err)
		return products
	}
	r := csv.NewReader(file)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("ReadCSV, while try to a read rows: ", err)
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

func WriteXlsx(path string, log *logrus.Logger, products []database.Product) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet("prices_to_update")
	if err != nil {
		log.Println(err)
		return
	}
	// Set value of a cell.
	_ = f.SetCellValue("prices_to_update", "A1", "Name")
	_ = f.SetCellValue("prices_to_update", "B1", "Old Price")
	_ = f.SetCellValue("prices_to_update", "C1", "New Price")

	for p, product := range products {
		price := strconv.FormatFloat(product.Price[len(product.Price)-1].Price, 'f', 2, 64)
		position := strconv.Itoa(p + 2)

		err := f.SetCellValue("prices_to_update", "A"+position, product.Name)
		if err != nil {
			log.Printf("WriteXlsx, while save Name of product %s provoque, error: %v", product.Name, err)
		}

		err = f.SetCellValue("prices_to_update", "B"+position, price)
		if err != nil {
			log.Printf("WriteXlsx, while save old price of product %s provoque, error: %v", product.Name, err)
		}

		err = f.SetCellValue("prices_to_update", "C"+position, price)
		if err != nil {
			log.Printf("WriteXlsx, while save new price of product %s provoque, error: %v", product.Name, err)
		}
	}

	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err := f.SaveAs("NewPrices.xlsx"); err != nil {
		log.Println("WriteXlsx, while saved file occur: ", err)
	}
}
