package main

import "io/ioutil"
import "fmt"
import "net/http"
import "encoding/csv"
import "os"
import "sort"
import "strconv"
import "errors"


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

	tickerRange := tickerRange{ticker, 2000, 01, 01, 2005, 12, 31}

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

type Logger interface {
	Info(format string, a ...interface{}) error
}

type TradeLogger struct {
	prefix string;
}

func (t TradeLogger) Info(format string, a ...interface{}) error {
	_, e := fmt.Printf("%s %s\n", t.prefix, fmt.Sprintf(format, a...));
	return e;
}

type Holding struct {
	// A holding is a single investment
	shares int
	cost currency // per-share purchase price
	ticker ticker
}

type Account struct {
	// An account holds cash and shares
	cash currency
	shares []Holding
}

/** 
 * Adds the Holding if there is sufficient funds and adjusts cash on hand.
 * cost is per-share cost.
 */
func (a *Account) BuyHolding(shares int, cost currency, ticker ticker) error {
	buycost := cost * currency(shares)
	if a.cash < buycost {
		return errors.New("Insufficient funds.");
	}

	h := Holding{shares, cost, ticker}
	a.shares = append(a.shares, h)
	a.cash -= buycost
	return nil
}

func (a *Account) SellHolding(shares int, saleprice currency, ticker ticker) (remaining_shares int, e error) {
	for _, h := range(a.shares) {
		thisround := shares
		if h.ticker != ticker { continue }
		if h.shares < thisround {
			thisround = h.shares
		}
		if h.shares >= thisround {
			h.shares -= thisround
			a.cash += currency(thisround) * saleprice
			shares -= thisround
		}
		if shares == 0 { break; }
	}

	return shares, nil
}

func simple_buy_if_down_sell_if_up(in chan TickerData, doneflag chan currency, log Logger) {
	cash := currency(1000.0)
	shares := 0

	// how to track gains?
	buycost := currency(0)
	finalsaleprice := currency(0.0)

	// If the stock ends lower than it opened, then buy at market.
	// If the stock ends high and we've made 10%, then sell at market.
	// This is sure to lose money, but I'll try the algorithm anyway.
	for v := range in {
		if v.closev < v.openv && cash > v.adjclose {
			// Buy
			// How many shares can we buy?
			// Assume we can buy at market open
			buyfor := v.adjclose
			sharesbuy := int(cash / buyfor)

			shares += sharesbuy
			buycost = currency(sharesbuy) * buyfor
			cash -= buycost

			log.Info("Bought %v shares at %v, cost %v and have %v cash remaining.", sharesbuy, buyfor, buycost, cash)
		}

		if v.closev > v.openv && (v.adjclose * currency(shares)) > (buycost * 1.10) {
			// Sell
			// Assume we sell all our shares
			// Assume we can sell at market open
			sellfor := v.adjclose;
			saleprice := currency(shares) * sellfor
			cash += saleprice
			log.Info("Selling %v shares at %v: Cash is %v", shares, sellfor, cash)
			shares -= shares // strange way to say shares = 0, but I don't want to change it.
		}

		finalsaleprice = v.adjclose
	}

	finalsell:=currency(shares)*finalsaleprice
	log.Info("Final sale: %v shares at %v yield %v", shares, finalsaleprice, finalsell)
	cash+=finalsell
	log.Info("Final cash: %v", cash)

	doneflag <- cash
}

func analyze(t []TickerData) error {
	fmt.Printf("\nBeginning trading\n");

	out := make(chan TickerData)
	doneflag := make(chan currency)
	log := TradeLogger{"simple"}

	go simple_buy_if_down_sell_if_up(out, doneflag, log)

	for _, v := range t {
		out <- v
	}

	close(out)

	<- doneflag

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




