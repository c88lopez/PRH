package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const csvsFolderPath = "csvs"

func main() {
	csvsFolder, err := os.Open(csvsFolderPath)
	if err != nil {
		log.Fatal(err)
	}

	csvFiles, err := csvsFolder.Readdir(0)
	if err != nil {
		log.Fatal(err)
	}

	for _, csvFile := range csvFiles {
		file, err := os.Open(fmt.Sprintf("%s%c%s", csvsFolderPath, os.PathSeparator, csvFile.Name()))
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		createConvertedFile(file)
	}
}

func createConvertedFile(file *os.File) {
	var newFileContent [][]string

	newFileContent = append(newFileContent, []string{"Clase de pedido", "Org. de Vtas.", "Canal de Dist.", "Solicitante", "Dest. de Merc", "Nro. Ord. Comp. Cli.", "Fe. Doc.", "Cond. de Pago", "Cond. de Exp.", "Motivo", "Clase Pedido Cli.", "Fe. Venc.", "Fe. creac. Ord. Comp.", "Fe. Ent.", "Moneda", "Material", "Material del Cliente", "PVP", "Descuento", "Cantidad"})

	csvReader := csv.NewReader(file)

	record, err := csvReader.Read()
	if err != nil {
		log.Fatal(err)
	}

	orderNumbers := record[2:]

	record, err = csvReader.Read() // i do nothing with this
	// record, err = csvReader.Read()
	if err != nil {
		log.Fatal(err)
	}

	clientAddresses := record[2:]

	var isbnOrders [][]string

	isbnCount := 0
	for {
		record, err = csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatal(err)
		}

		isbnCount++
		isbnOrders = append(isbnOrders, record)
	}

	for orderNumberIndex, orderNumber := range orderNumbers {
		for i := 0; i < isbnCount; i++ {

			if strings.TrimSpace(isbnOrders[i][2+orderNumberIndex]) == "" { // no order for the current isbn
				continue
			}

			newFileContent = append(newFileContent, []string{"", "", "", "", clientAddresses[orderNumberIndex],
				orderNumber, "", "", "", "", "", "", "", "", "", isbnOrders[i][0], "", "", "",
				strings.TrimSpace(isbnOrders[i][2+orderNumberIndex])})
		}
	}

	convertedFile, err := os.Create(strings.Replace(file.Name(), ".csv", " converted.csv", 1))
	if err != nil {
		log.Fatal(err)
	}

	w := csv.NewWriter(convertedFile)
	w.WriteAll(newFileContent) // calls Flush internally

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
}
