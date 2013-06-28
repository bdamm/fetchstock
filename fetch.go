package main

import "io/ioutil"
import "fmt"
import "net/http"
import "encoding/csv"
import "os"
import "sort"
import "strconv"


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

func fetchSymbol(sym string, filename string) error {
	ticker := ticker{sym}

	tickerRange := tickerRange{ticker, 2000, 01, 01, 2000, 12, 31}

	url,err := tickerRange.ToYahooUrl()
	if err != nil { return err }

	resp, err := http.Get(*url)
	if err != nil { return err }

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	} else {
		fmt.Println(string(body))
	}

	// write whole the body
	err = ioutil.WriteFile(filename, body, 0644)
	if err != nil { return err }

	return nil
}

type currency float32
type TickerData struct {
	date string
	openv currency
	highv currency
	lowv currency
	closev currency
	volume int
	adjclose currency
}

func pc(s string) currency {
	a, e := strconv.ParseFloat(s, 32)
	if e != nil {panic(e)}
	return currency(a)
}
func pi(s string) int {
	a, e := strconv.ParseInt(s, 10, 0)
	if e != nil {panic(e)}
	return int(a)
}

type TickerDataSlice []TickerData

func (t TickerDataSlice) Len() int {
	return len(t)
}
func (t TickerDataSlice) Less(i, j int) bool {
	return t[i].date < t[j].date
}
func (t TickerDataSlice) Swap(i, j int) {
	t[i],t[j] = t[j],t[i]
}


func readTickerData(filename string) ([]TickerData, error) {
	filereader, err := os.Open(filename)
	if err != nil { return nil, err }
	r := csv.NewReader(filereader)
	records, err := r.ReadAll()
	if err != nil { return nil, err }
	records = records[1:] // chomp the first line; it's the header

	result := TickerDataSlice(make([]TickerData, len(records)))
	for i := 0; i < len(records); i++ {
		r := records[i];
		d := TickerData{r[0], pc(r[1]), pc(r[2]), pc(r[3]), pc(r[4]), pi(r[5]), pc(r[6])}
		result[i] = d
	}
	sort.Sort(result)

	return result, nil
}

func analyze(t []TickerData) error {

	// Just do the analysis here, I'll refactor it out later.

	// If the stock ends lower than it opened, then buy at market.
	// If the stock ends high and we've made 10%, then sell at market.
	// This is sure to lose money, but I'll try the algorithm anyway.
	cash := currency(1000.0)
	shares := 0

	// how to track gains?
	buycost := currency(0)

	for _, v := range t {
		if v.closev < v.openv && cash > 0 {
			// Buy
			fmt.Println("Buying.")
			// How many shares can we buy?
			// Assume we can buy at market open
			sharesbuy := int(cash / v.adjclose)

			shares += sharesbuy
			buycost = currency(sharesbuy) * v.adjclose
			cash -= currency(sharesbuy) * v.adjclose

		}

		if v.closev > v.openv && (v.adjclose * currency(shares)) > (buycost * 1.10) {
			// Sell
			fmt.Println("Selling.")
			// Assume we sell all our shares
			// Assume we can sell at market open
			saleprice := currency(shares) * v.adjclose
			cash += saleprice
			shares -= shares // strange way to say shares = 0, but I don't want to change it.

		}
	}

	return nil
}

func main() {
	fmt.Println("Starting")

	sym := "GE"

	filename := "data."+sym

	err := fetchSymbol(sym, filename)
	if err != nil { panic(err) }

	// ...
	// Ok, now read it in.
	data, err := readTickerData(filename);

	for v := range data {
		fmt.Println(data[v])
	}

	analyze(data)
}




