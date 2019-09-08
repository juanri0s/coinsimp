package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"
)

// Data represents the response that we expect from an api single coin call
type Data struct {
	// Coin represents the specific coin data
	Coin struct {
		ID                string `json:"id"`
		Rank              string `json:"rank"`
		Symbol            string `json:"symbol"`
		Name              string `json:"name"`
		Supply            string `json:"supply"`
		MaxSupply         string `json:"maxSupply"`
		MarketCapUsd      string `json:"marketCapUsd"`
		VolumeUsd24Hr     string `json:"volumeUsd24Hr"`
		PriceUsd          string `json:"priceUsd"`
		ChangePercent24Hr string `json:"changePercent24Hr"`
		Vwap24Hr          string `json:"vwap24Hr"`
	} `json:"data"`
	Timestamp int64 `json:"timestamp"`
}

func main() {
	url := "https://api.coincap.io/v2/assets/"
	method := "GET"

	var c string
	p := "coin"
	flag.StringVar(&c, p, "", "Usage")
	flag.Parse()

	fc := isFlagPassed(p)
	if fc != true {
		fmt.Println("Welcome to simpcoin, please provide a coin that you're interested in using the -coin= flag")
		return
	}

	t := time.Now().Local()
	fmt.Println("simpcoin started -", t)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest(method, url+c, nil)

	if err != nil {
		fmt.Println(err)
	}

	res, err := client.Do(req)

	if err != nil {
		fmt.Printf("The HTTP %s request to %s failed with error %s\n", method, url, err)
	}

	if res.StatusCode <= 200 && res.StatusCode >= 299 {
		fmt.Println("Unsuccessful HTTP request:", res.StatusCode, http.StatusText(res.StatusCode))
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	d := Data{}
	jserr := json.Unmarshal(body, &d)
	if jserr != nil {
		log.Fatal(jserr)
	}

	buildCoinTable(d)
}

// buildTable takes in data for a single coin and outputs it in ASCII table format
func buildCoinTable(d Data) {

	t := tablewriter.NewWriter(os.Stdout)

	buildHeaders(t, d)
	buildRows(t, d)

	t.Render() // Send output
}

// buildHeaders generates the headers slice for the table.
func buildHeaders(table *tablewriter.Table, d Data) {
	v := reflect.ValueOf(d.Coin)
	n := v.NumField()
	typeOfS := v.Type()

	h := make([]string, 0)
	for i := 0; i < n; i++ {
		h = append(h, typeOfS.Field(i).Name)
	}
	table.SetHeader(h)
}

// buildRows generates the rows slice for the table.
func buildRows(table *tablewriter.Table, d Data) {
	v := reflect.ValueOf(d.Coin)
	n := v.NumField()

	r := make([]string, 0)
	for i := 0; i < n; i++ {
		x := v.Field(i)
		s := fmt.Sprintf("%v", x.Interface())
		r = append(r, s)
	}
	table.Append(r)
}

// isFlagPassed checks to see that the required flags were inserted.
func isFlagPassed(p string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == p {
			found = true
		}
	})

	return found
}
