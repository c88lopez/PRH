package main

import (
	"sync"
)

const csvsFolderPath = "csvs"
const skipFileRegex = ".gitkeep|converted|.xls$"

var waitGroup sync.WaitGroup

func main() {
	Start()
}
