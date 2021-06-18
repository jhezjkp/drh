package main

import (
	"drh/internal/drh"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func main() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b, err := ioutil.ReadFile(path + "/test/prices.json")
	if err != nil {
		log.Fatal(err)
	}

	p := drh.ParseData(b)

	year := 3
	tlds := "date,audio"
	registrars := "Porkbun"
	ps := p.FindCheapestRegistrar(year, false, false, tlds, registrars, 10)

	// 输出
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetAutoIndex(true)
	t.AppendHeader(table.Row{"Tld", "Total", "Registrar", "Register", "Renew", "Transfer", "FreeWhois"})
	for _, s := range ps {
		t.AppendRow(table.Row{s.Tld, Decimal(s.TotalFee), s.Registrar,
			Decimal(s.RegisterPrice), Decimal(s.RenewPrice),
			Decimal(s.TransferPrice), s.FreeWhois})
	}
	t.Render()
}
