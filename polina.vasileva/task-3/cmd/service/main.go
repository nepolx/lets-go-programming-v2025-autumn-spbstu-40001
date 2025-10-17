package main

import (
	"flag"
	"slices"

	"polina.vasileva/task-3/internal/config"
	"polina.vasileva/task-3/internal/currency"
	"polina.vasileva/task-3/internal/json"
	"polina.vasileva/task-3/internal/xml"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config")
	flag.Parse()

	config, err := config.ParseYaml(*configPath)
	if err != nil {
		panic(err)
	}

	currencyList := currency.Rates{Data: []currency.Currency{}}

	err = xml.ParseXML(config.InputFilePath, &currencyList)
	if err != nil {
		panic(err)
	}

	slices.SortStableFunc(currencyList.Data, currency.DescendingComparatorCurrency)

	err = json.ParseJSON(config.OutputFilePath, currencyList.Data)
	if err != nil {
		panic(err)
	}
}
