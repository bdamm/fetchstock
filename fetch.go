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
	ticker := ticker{"GE"}

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

type TickerData struct {
	date string
	openv float32
	highv float32
	lowv float32
	closev float32
	volume int
	adjclose float32
}

func pc(s string) (float32, error) {
	a, e := strconv.ParseFloat(s, 32)
	return float32(a), e
}

func (t []TickerData) Len() int {
	return len(t)
}

func readTickerData(filename string) ([]TickerData, error) {
	filereader, err := os.Open(filename)
	if err != nil { return nil, err }
	r := csv.NewReader(filereader)
	records, err := r.ReadAll()
	if err != nil { return nil, err }
	records = records[1:] // chomp the first line; it's the header

	result := make([]TickerData, len(records))
	for r := range(records) {
		d := TickerData(r[0], pc(r[1]), pc(r[2]), pc(r[3]), pc(r[4]), pc(r[5]), pc(r[6]))
	}
	sort.Sort(records)

	return nil, nil
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
		fmt.Println(v)
	}
}
