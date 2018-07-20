package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

const csvsFolderPath = "csvs"
const skipFileRegex = ".gitkeep|converted|.xls$"

var waitGroup sync.WaitGroup

func main() {
	log.Println("Starting...")
	start := time.Now()

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

		if skipFile(file) {
			continue
		}

		if err != nil {
			log.Fatal("Error opening csv file", fmt.Sprintf("%s%c%s", csvsFolderPath, os.PathSeparator, csvFile.Name()), err)
		}
		defer file.Close()

		waitGroup.Add(1)
		go createConvertedFile(file)
	}

	waitGroup.Wait()

	log.Printf("End - Elapsed time: %s.\n", time.Since(start))
}

func createConvertedFile(file *os.File) {
	defer waitGroup.Done()

	csvReader := csv.NewReader(file)

	record, err := csvReader.Read()
	if err != nil {
		log.Fatal("Error reading csv record (first) ", file.Name(), err)
	}

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

	generateConvertedFile(
		file,
		generateFileContent(isbnCount, isbnOrders, orderNumbers, clientAddresses))
}

func skipFile(file *os.File) bool {
	return regexp.MustCompile(skipFileRegex).MatchString(file.Name())
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

func generateConvertedFile(file *os.File, content [][]string) {
	convertedFile, err := os.Create(strings.Replace(
		file.Name(), ".csv", " converted.csv", 1))
	if err != nil {
		log.Fatal("Error creating converted csv file ",
			strings.Replace(file.Name(), ".csv", " converted.csv", 1), err)
	}

	w := csv.NewWriter(convertedFile)
	w.WriteAll(content) // calls Flush internally

	if err := w.Error(); err != nil {
		log.Fatalln("error writing converted csv file ", err)
	}
}
