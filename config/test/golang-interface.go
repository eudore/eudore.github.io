package main

import (
	"fmt"
)

type (
	Printer interface {
		Print()
	}
	BasePrinter struct {

	}
	StdPrinter struct {
		Printer
	}
)

func (p *BasePrinter) Print() {
	fmt.Println("BasePrinter")
}

func main() {
	p := &StdPrinter{Printer: &BasePrinter{}}
	p.Print()
}