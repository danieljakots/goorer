package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v2"
)

const dataPath = "testdata/"

type moneyExchange struct {
	Amount float64
	Date time.Time
	With string
}

func readCategoriesFile(categoriesFilePath string) (map[string]string, error) {
	yamlFile, err := ioutil.ReadFile(categoriesFilePath)
	if err != nil {
		return nil, err
	}
	categories := make(map[string]string)
	err = yaml.Unmarshal(yamlFile, categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func readMonthlyFile(monthlyFilePath string) (map[string][]moneyExchange, error){
	m := make(map[string][]moneyExchange)

	yamlFile, err := ioutil.ReadFile(monthlyFilePath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &m)
	if err != nil {
		return nil, err
	}
	return m, nil

}

func main() {
	_, err := readCategoriesFile(dataPath + "categories.yml")
	if err != nil {
		log.Fatal("Couldn't parse categories file: ", err)
	}
	// fmt.Println(categories)

	entries, err := readMonthlyFile(dataPath + "december-20.yml")
	if err != nil {
		log.Fatal("Couldn't parse records file: ", err)
	}
	fmt.Println(entries)

	var sum float64
	for _, spending := range(entries["spendings"]) {
		sum = sum + spending.Amount
	}
	fmt.Println(sum)
}
