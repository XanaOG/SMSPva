package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/XanaOG/SmsPVA-PriceCheck/Core/Client"
)

type Response struct {
	Response int    `json:"response"`
	Country  string `json:"country"`
	Service  string `json:"service"`
	Price    string `json:"price"`
}

func main() {
	Config := Client.GetConfig(Client.CountryFile)
	if len(os.Args) < 2 {
		for _, service := range Config.List.Options {
			fmt.Printf("%s - %d\n", service.Name, service.Number)
		}
		return
	}

	lowestPrices := make(map[string]string)
	for _, country := range Config.List.Countries {
		req, _ := http.NewRequest("GET", "https://smspva.com/priemnik.php?metod=get_service_price&country="+country+"&service=opt"+os.Args[1]+"&apikey="+Config.APIKey, nil)

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		var priceResponse Response
		json.Unmarshal(body, &priceResponse)
		if priceResponse.Price != "" {
			if _, ok := lowestPrices[priceResponse.Country]; !ok {
				lowestPrices[priceResponse.Country] = priceResponse.Price
			}
		}
	}

	if len(lowestPrices) == 0 {
		fmt.Println("No prices found.")
		return
	}

	type priceCountry struct {
		Price   float64
		Country string
	}

	var prices []priceCountry

	for country, price := range lowestPrices {
		priceFloat, err := strconv.ParseFloat(price, 64)
		if err != nil {
			continue
		}
		prices = append(prices, priceCountry{Price: priceFloat, Country: country})
	}

	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Price < prices[j].Price
	})

	fmt.Println("Top 5 lowest prices:")
	for i, pc := range prices {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. $%.2f in %s\n", i+1, pc.Price, pc.Country)
	}
}
