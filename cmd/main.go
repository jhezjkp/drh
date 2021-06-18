package main

import (
	"drh/internal/drh"
	"flag"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func downloadAndRead(url string, path string) ([]byte, error) {
	f := url[strings.LastIndex(url, "/")+1:]
	fullPath := path + f
	fi, err := os.Stat(fullPath)
	if os.IsNotExist(err) || (err == nil && time.Now().Sub(fi.ModTime()).Hours() > 2) {
		out, err := os.Create(fullPath)
		if err != nil {
			return nil, err
		}
		defer out.Close()

		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return nil, err
		}
	}

	return ioutil.ReadFile(fullPath)
}

func main() {
	year := flag.Int("year", 1, "specific year")
	count := flag.Int("count", 10, "result size, default 10")
	renew := flag.Bool("renew", false, "domain renew")
	transfer := flag.Bool("transfer", false, "domain transfer")
	tlds := flag.String("tld", "", "specific tlds, comma seperated, eg: com,net")
	registrars := flag.String("registrar", "", "specific registrars, comma seperated, eg: namecheap,porkbun")

	flag.Parse()

	url := "https://www.domcomp.com/prices.json"
	path := os.TempDir()

	b, err := downloadAndRead(url, path)
	if err != nil {
		log.Fatal(err)
	}

	p := drh.ParseData(b)

	ps := p.FindCheapestRegistrar(*year, *renew, *transfer, *tlds, *registrars, *count)

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
