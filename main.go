package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	url := "https://api.coincap.io/v2/assets"
	method := "GET"

	var c string
	p := "coin"
	flag.StringVar(&c, p, "", "Usage")
	flag.Parse()

	fc := isFlagPassed(p)
	if fc != true {
		fmt.Println("Welcome to coinsimp, please provide a coin that you're interested in using the -coin flag")
		return
	}

	t := time.Now().Local()
	fmt.Println("coinsimp started -", t)

	fmt.Println("requested coin:", c)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("The HTTP %s request to %s failed with error %s\n", method, url, err)
	}

	if resp.StatusCode <= 200 && resp.StatusCode >= 299 {
		fmt.Println("Something happened:", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Printf("%s\n", body)
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
