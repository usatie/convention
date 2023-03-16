package main

import (
	"convention"

	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(convention.Analyzer) }
