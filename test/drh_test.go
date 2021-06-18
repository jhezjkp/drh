package main

import (
	"drh/internal/drh"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestFib(t *testing.T) {
	type args struct {
		year       int
		renew      bool
		transfer   bool
		tlds       string
		registrars string
		count      int
	}
	tests := []struct {
		name string
		args args
		want drh.PriceSummaries
	}{
		// register
		{"newRegisterFor1Year", args{1, false, false, "", "", 2},
			drh.PriceSummaries{{"eu", "GoDaddy", 0.01, 0.01, 12.99, 12.99,
				false, "", "",
				"<small>An extra 30% ~ 60% discount will be auto applied for first year after transfer</small>"},
				{"com.mx", "GoDaddy", 0.1, 0.1, 29.99, 29.99,
					false, "", "",
					"<small>An extra 30% ~ 60% discount will be auto applied for first year after transfer</small>"}}},
		{"newRegisterFor2Year", args{2, false, false, "", "", 1},
			drh.PriceSummaries{{"bid", "Porkbun", 5.27, 2.53, 2.74, 2.74,
				true, "<center>Free SSL</center>", "<center>Free SSL</center>",
				"<center>Free SSL</center>"}}},
		// renew
		{"renewFor1Year", args{1, true, false, "", "", 1},
			drh.PriceSummaries{{"accountant", "Porkbun", 2.74, 11.48, 2.74, 10.98,
				true, "<center>Free SSL</center>", "<center>Free SSL</center>",
				"<center>Free SSL</center>"}}},
		// transfer
		{"transfer", args{1, false, true, "", "", 1},
			drh.PriceSummaries{{"bid", "Porkbun", 2.74, 2.53, 2.74, 2.74,
				true, "<center>Free SSL</center>", "<center>Free SSL</center>",
				"<center>Free SSL</center>"}}},
		// transferAndRenew
		{"transferAndRenewFor2Year", args{3, false, true, "", "", 1},
			drh.PriceSummaries{{"bid", "Porkbun", 8.22, 2.53, 2.74, 2.74,
				true, "<center>Free SSL</center>", "<center>Free SSL</center>",
				"<center>Free SSL</center>"}}},
		// filter by tld and registrar
		{"filter", args{3, false, false, "date,audio", "Porkbun", 2},
			drh.PriceSummaries{{"date", "Porkbun", 8.01, 2.53, 2.74, 2.74,
				true, "<center>Free SSL</center>", "<center>Free SSL</center>",
				"<center>Free SSL</center>"},
				{"audio", "Porkbun", 279.03000000000003, 70.01, 104.51, 104.51,
					true, "<center>Free SSL</center>", "<center>Free SSL</center>",
					"<center>Free SSL</center>"}}},
	}

	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b, err := ioutil.ReadFile(path + "/prices.json")
	if err != nil {
		log.Fatal(err)
	}

	p := drh.ParseData(b)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := p.FindCheapestRegistrar(tt.args.year, tt.args.renew, tt.args.transfer, tt.args.tlds, tt.args.registrars, tt.args.count); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findCheapestRegistar() = %v, want %v", got, tt.want)
			}
		})
	}

}
