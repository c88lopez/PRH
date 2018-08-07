package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strings"
)

// Convert Convert
func Convert(file io.Reader) io.Reader {
	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'

	record, err := csvReader.Read()
	if err != nil {
		log.Fatal("Error reading csv record (first) ", err)
	}

	fmt.Printf("%#v", record)

	orderNumbers := record[2:]

	record, err = csvReader.Read() // i do nothing with this
	if err != nil {
		log.Fatal("Error reading csv record (second) ", err)
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

			log.Fatal("Error reading records in loop ", err)
		}

		isbnCount++
		isbnOrders = append(isbnOrders, record)
	}

	buf := bytes.NewBufferString("")
	csvWriter := csv.NewWriter(buf)

	csvWriter.WriteAll(generateFileContent(
		isbnCount, isbnOrders, orderNumbers, clientAddresses))

	return buf
}

func generateFileContent(isbnCount int, isbnOrders [][]string, orderNumbers []string, clientAddresses []string) [][]string {
	var newFileContent [][]string

	newFileContent = append(newFileContent, []string{"Clase de pedido", "Org. de Vtas.", "Canal de Dist.",
		"Solicitante", "Dest. de Merc", "Nro. Ord. Comp. Cli.", "Fe. Doc.", "Cond. de Pago", "Cond. de Exp.",
		"Motivo", "Clase Pedido Cli.", "Fe. Venc.", "Fe. creac. Ord. Comp.", "Fe. Ent.", "Moneda", "Material",
		"Material del Cliente", "PVP", "Descuento", "Cantidad"})

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

	return newFileContent
}
