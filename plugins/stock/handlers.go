// Package stock implements a plugin to get stock quotes.
package stock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/transport"
)

const quoteURL = "https://api.iextrading.com/1.0/stock/%s/quote"

type quote struct {
	Symbol                string
	CompanyName           string
	PrimaryExchange       string
	Sector                string
	CalculationPrice      string
	Open                  float64
	OpenTime              int64
	Close                 float64
	CloseTime             int64
	High                  float64
	Low                   float64
	LatestPrice           float64
	LatestSource          string
	LatestTime            string
	LatestUpdate          int64
	LatestVolume          int64
	IexRealtimePrice      float64
	IexRealtimeSize       int64
	iexLastUpdated        int64
	DelayedPrice          float64
	DelayedPriceTime      int64
	ExtendedPrice         float64
	ExtendedChange        float64
	ExtendedChangePercent float64
	ExtendedPriceTime     int64
	PreviousClose         float64
	Change                float64
	ChangePercent         float64
	IexMarketPercent      float64
	IexVolume             int64
	AvgTotalVolume        int64
	IexBidPrice           float64
	IexBidSize            int64
	IexAskPrice           float64
	IexAskSize            int64
	MarketCap             int64
	PeRatio               float64
	Week52High            float64
	Week52Low             float64
	YtdChange             float64
}

func (q *quote) String() string {
	plus := ""
	if q.Change > 0 {
		plus = "+"
	}
	return fmt.Sprintf("%v (%v): %v %v (%v%v, %v%.2f%%), %v market cap, %v volume, %.2f P/E - https://iextrading.com/apps/stocks/%v",
		q.Symbol, q.CompanyName,
		q.LatestSource, q.LatestPrice, plus, q.Change, plus, q.ChangePercent*100,
		humanize(q.MarketCap), humanize(q.LatestVolume), q.PeRatio, q.Symbol)
}

func humanize(i int64) string {
	format := func(i, unit int64) string {
		f := "%.0f"
		if i < 10*unit {
			f = "%.2f"
		} else if i < 100*unit {
			f = "%.1f"
		}
		return fmt.Sprintf(f, float64(i)/float64(unit))
	}
	switch {
	case i >= 1e9: // Billion
		return fmt.Sprintf("%vB$", format(i, 1e9))
	case i >= 1e6: // Million
		return fmt.Sprintf("%vM$", format(i, 1e6))
	case i >= 1e3: // Thousand
		return fmt.Sprintf("%vK$", format(i, 1e3))
	}
	return fmt.Sprintf("%v$", i)
}

func stock(symbol string) (*quote, error) {
	uri := fmt.Sprintf(quoteURL, url.PathEscape(symbol))
	c, err := transport.Client(uri)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v", resp.Status)
	}
	var v quote
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return &v, nil
}

func handleStock(e *bot.Event) {
	symbol := strings.TrimSpace(e.Args)
	if len(symbol) == 0 {
		return
	}
	q, err := stock(symbol)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	e.Bot.Privmsg(e.Target, q.String())
}

// Register registers the plugin with a bot.
func Register(b bot.Bot) {
	b.Commands().Add("stock", bot.Command{
		Help:    "get trading stock information like price",
		Handler: handleStock,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
