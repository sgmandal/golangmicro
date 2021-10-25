package data

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/go-hclog"
)

type ExchangeRates struct {
	l     hclog.Logger // number of different loggers available in this package
	rates map[string]float64
}

func NewRates(log hclog.Logger) (*ExchangeRates, error) { // idiomatic go says use addresses incase of structs
	er := &ExchangeRates{l: log, rates: map[string]float64{}}
	_ = er.getRates()
	return er, nil
}

func (e *ExchangeRates) getRates() error {
	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")

	// rates are based on euros
	if err != nil {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected error code 200 got %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	md := &Cubes{}
	xml.NewDecoder(resp.Body).Decode(&md)

	for _, c := range md.CubeData {
		r, err := strconv.ParseFloat(c.Rate, 64) // conversion from string to float64
		if err != nil {
			return err
		}

		e.rates[c.Currency] = r
	}

	e.rates["EUR"] = 1

	return nil
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"` // specifying we need attribute, try not to keep same variable and xml names
	Rate     string `xml:"rate,attr"`
}

func (e *ExchangeRates) GetRae(base, dest string) (float64, error) {
	br, ok := e.rates[base] // or returns if the key-value pair is present or not - boolean
	if !ok {
		return 0, fmt.Errorf("rate not found for currency %s", base)
	}

	dr, ok := e.rates[dest]
	if !ok {
		return 0, fmt.Errorf("rate not found for currency %s", base)
	}

	return dr / br, nil
}
