package main

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

type TickerError string
func (t TickerError) Error() string {
	return string(t);
}

func (t tickerRange) ToYahooUrl() (*string, error) {
	if t.startmonth == 0 || t.endmonth == 0 { return nil, TickerError("months are expressed 1-12.") }
	url := fmt.Sprintf("http://ichart.yahoo.com/table.csv?s=%v&a=%v&b=%v&c=%v&d=%v&e=%v&f=%v&g=d&ignore=.csv",
		t.ticker.symbol, t.startmonth-1, t.startday, t.startyear, t.endmonth-1, t.endday, t.endyear)

	return &url, nil;
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

	url,err := tickerRange.ToYahooUrl()
	if err != nil { fmt.Println(err); return }

	resp, err := http.Get(*url)
	if err != nil { fmt.Println(err); return }
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	} else {
		//fmt.Println("p ==", body)
		fmt.Println(string(body))
	}

	// write whole the body
	err = ioutil.WriteFile("data."+ticker.symbol+"", body, 0644)
	if err != nil {
		panic(err)
	}

	// ...
}
