package main

import "fmt"
import "net/http"
import "io/ioutil"

func main() {
	fmt.Println("Starting");

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



	symbol := "GE"
	startyear := 2000
	endyear := 2000
	url := fmt.Sprintf("http://ichart.yahoo.com/table.csv?s=%v&a=00&b=01&c=%v&d=00&e=30&f=%v&g=d&ignore=.csv", symbol, startyear, endyear)
	resp, err := http.Get(url);
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	} else {
		//fmt.Println("p ==", body)
		fmt.Println(string(body));
	}



	// ...

}
