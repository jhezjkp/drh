package drh

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type TldRegistrarPriceEntry struct {
	Key        string `json:"key"`
	Registrar  string `json:"disp"`
	Price      string `json:"Price"`
	FreeWhois  bool   `json:"free_whois"`
	CouponInfo string `json:"coupon_info"`
}

func (t TldRegistrarPriceEntry) GetPrice() float64 {
	price, err := strconv.ParseFloat(strings.ReplaceAll(t.Price, ",", ""), 64)
	if err == nil {
		return price
	}
	return 0
}

type TldItem struct {
	Tld     string                   `json:"tld"`
	Entries []TldRegistrarPriceEntry `json:"entries"`
}

type Price struct {
	New      []TldItem `json:"new"`
	Renew    []TldItem `json:"renew"`
	Transfer []TldItem `json:"transfer"`
}

type PriceSummary struct {
	Tld                string
	Registrar          string
	TotalFee           float64
	RegisterPrice      float64
	RenewPrice         float64
	TransferPrice      float64
	FreeWhois          bool
	RegisterCouponInfo string
	RenewCouponInfo    string
	TransferCouponInfo string
}

type PriceSummaries []PriceSummary

func (s PriceSummaries) Len() int {
	return len(s)
}

func (s PriceSummaries) Less(i, j int) bool {
	r := s[i].TotalFee - s[j].TotalFee
	if r == 0 {
		c := strings.Compare(s[i].Tld, s[j].Tld)
		if c == 0 {
			return strings.Compare(s[i].Registrar, s[j].Registrar) < 0
		} else {
			return c < 0
		}
	}
	// i<j: true
	return r < 0
}

func (s PriceSummaries) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (p Price) findByTld(tld string) {
	for _, item := range p.New {
		if item.Tld == tld {
			fmt.Println(item)
			break
		}
	}
}

func (p Price) FindCheapestRegistrar(year int, renew bool, transfer bool, strTlds string, strRegistrars string, count int) PriceSummaries {
	priceMap := make(map[string]PriceSummary)
	for _, item := range p.New {
		for _, entry := range item.Entries {
			s := PriceSummary{}
			s.Tld = item.Tld
			s.Registrar = entry.Registrar
			s.FreeWhois = entry.FreeWhois
			s.RegisterPrice = entry.GetPrice()
			s.RegisterCouponInfo = entry.CouponInfo
			if !renew && !transfer {
				s.TotalFee += entry.GetPrice()
			}
			key := item.Tld + "_" + entry.Registrar
			priceMap[key] = s
		}
	}
	for _, item := range p.Transfer {
		for _, entry := range item.Entries {
			key := item.Tld + "_" + entry.Registrar
			s, ok := priceMap[key]
			if !ok {
				if !transfer { // 仅在转入模式下做新增
					continue
				}
				s = PriceSummary{}
				s.Tld = item.Tld
				s.Registrar = entry.Registrar
				s.FreeWhois = entry.FreeWhois
			}
			s.TransferPrice = entry.GetPrice()
			s.TransferCouponInfo = entry.CouponInfo
			if transfer {
				s.TotalFee += entry.GetPrice()
			}
			priceMap[key] = s
		}
	}

	for _, item := range p.Renew {
		for _, entry := range item.Entries {
			key := item.Tld + "_" + entry.Registrar
			s, ok := priceMap[key]
			if !ok {
				if year <= 1 && !renew { // 非续费模式且只1年的，过滤掉
					continue
				}
				s = PriceSummary{}
				s.Tld = item.Tld
				s.Registrar = entry.Registrar
				s.FreeWhois = entry.FreeWhois
			}
			s.RenewPrice = entry.GetPrice()
			s.RenewCouponInfo = entry.CouponInfo
			if renew { // 纯续费
				s.TotalFee += entry.GetPrice() * float64(year)
			} else { // 注册/转移+续费
				s.TotalFee += entry.GetPrice() * float64(year-1)
			}
			priceMap[key] = s
		}
	}
	ps := PriceSummaries{}
	tlds := make(map[string]bool)
	if len(strTlds) > 0 {
		for _, s := range strings.Split(strings.ToLower(strTlds), ",") {
			tlds[s] = true
		}
	}
	registrars := make(map[string]bool)
	if len(strRegistrars) > 0 {
		for _, s := range strings.Split(strings.ToLower(strRegistrars), ",") {
			registrars[s] = true
		}
	}
	for _, value := range priceMap {
		if year <= 1 && transfer && value.TransferPrice <= 0 { // 如果只是转入1年，则过滤掉没有转入价格的
			continue
		}
		// 根据tld过滤
		if len(tlds) > 0 {
			if _, ok := tlds[strings.ToLower(value.Tld)]; !ok {
				continue
			}
		}
		// 根据registrar过滤
		if len(registrars) > 0 {
			if _, ok := registrars[strings.ToLower(value.Registrar)]; !ok {
				continue
			}
		}
		ps = append(ps, value)
	}
	sort.Sort(ps)
	if count > 0 && ps.Len() > count {
		return ps[:count]
	}
	return ps
}

func ParseData(bytes []byte) *Price {
	price := &Price{}
	json.Unmarshal(bytes, &price)
	return price
}
