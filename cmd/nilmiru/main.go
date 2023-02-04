package main

import (
	"github.com/speed1313/nilmiru"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(nilmiru.Analyzer) }
