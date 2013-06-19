package main

import "errors"
import "io/ioutil"
import "fmt"
import "net/http"

type ticker struct {
	symbol string
}
type tickerRange struct {
	ticker     ticker
	startyear  int
	startmonth int
	startday   int
	endyear    int
	endmonth   int
	endday     int
}

func (t *tickerRange) ToYahooUrl() (string, error) {
	if t.startmonth == 0 || t.endmonth == 0 { return "", errors.New("months are expressed 1-12.") }
	url := fmt.Sprintf("http://ichart.yahoo.com/table.csv?s=%v&a=00&b=01&c=%v&d=11&e=31&f=%v&g=d&ignore=.csv",
		t.ticker, t.startmonth-1, t.startday, t.startyear, t.endmonth-1, t.endday, t.endyear)

	return url, nil;
}

func main() {
	fmt.Println("Starting")

	//	client := &http.Client{
	//		CheckRedirect: redirectPolicyFunc,
	//	}
	//
	//	resp, err := client.Get("http://example.com")
	//	// ...
	//
	//	req, err := http.NewRequest("GET", "http://example.com", nil)
	//	// ...
	//	req.Header.Add("If-None-Match", `W/"wyzzy"`)
	//	resp, err := client.Do(req)
	//	// ...

	ticker := ticker{"GE"}

	tickerRange := tickerRange{ticker, 2000, 01, 01, 2000, 12, 31}
	
	tickerRange.ToYahooUrl();

	url := fmt.Sprintf("http://ichart.yahoo.com/table.csv?s=%v&a=00&b=01&c=%v&d=11&e=31&f=%v&g=d&ignore=.csv", symbol, startyear, endyear)
	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	} else {
		//fmt.Println("p ==", body)
		fmt.Println(string(body))
	}

	// write whole the body
	err = ioutil.WriteFile("data."+symbol+"", body, 0644)
	if err != nil {
		panic(err)
	}

	// ...
}
