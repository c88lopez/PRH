package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
)

const csvFilePath = "./2018-01-22 PRH-Reposicion en consignacion para sucursales numero de OC 1121430.csv"

func main() {
	var newFileContent [][]string
	newFileContent = append(newFileContent, []string{"Clase de pedido", "Org. de Vtas.", "Canal de Dist.", "Solicitante", "Dest. de Merc", "Nro. Ord. Comp. Cli.", "Fe. Doc.", "Cond. de Pago", "Cond. de Exp.", "Motivo", "Clase Pedido Cli.", "Fe. Venc.", "Fe. creac. Ord. Comp.", "Fe. Ent.", "Moneda", "Material", "Material del Cliente", "PVP", "Descuento", "Cantidad"})

	file, err := os.Open(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}

	csvReader := csv.NewReader(file)

	record, err := csvReader.Read()
	if err != nil {
		log.Fatal(err)
	}

	orderNumbers := record[2:]

	record, err = csvReader.Read() // i do nothing with this
	record, err = csvReader.Read()
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

	w := csv.NewWriter(os.Stdout)
	w.WriteAll(newFileContent) // calls Flush internally

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
}
