/*this file only fetches and organizes the rates and returns it, nothing more, nothing less*/

package data

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-hclog"
)

type ExchangeRates struct {
	log hclog.Logger

	rates map[string]float64 //this is returned as rates of each currency which is present or updated
}

func NewRates(log hclog.Logger) (*ExchangeRates, error) {
	er := &ExchangeRates{log: log, rates: map[string]float64{}} // no data is passed in linear fashion hence the curly parentheses are empty

	_ = er.okRates() // returns an error, we can throw that or check error
	return er, nil
}

// helper method to fetch rates
func (e *ExchangeRates) okRates() error {
	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml") // getting the live data to process further in xml format
	if err != nil {
		return nil
	}

	if resp.StatusCode != http.StatusOK { // http package return in resp and checking if everything is ok or not
		return fmt.Errorf("expected error code 200 %d", resp.StatusCode)
	}
	defer resp.Body.Close() // not leaving any open(ghost) connections

	md := &Cubes{} // creating an empty container to hold data
	xml.NewDecoder(resp.Body).Decode(&md)

	for _, c := range md.cubedata {
		r, err := strconv.ParseFloat(c.Rate, 64) // Rate conversion to float64
		if err != nil {
			return err
		}

		e.rates[c.cx] = r // here we'll have all the data in exchange rates struct
	}

	e.rates["EUR"] = 1 // adding the base rate

	return nil

}

type Cube struct {
	cx   string `xml:"currency,attr"`
	Rate string `xml:"rate,attr"`
}

type Cubes struct {
	cubedata []Cube `xml:"Cube>Cube>Cube"` // slice of cube struct, and this is how the xml format is, so we extract the deep down data, telling the compiler what we need specifiaclly
}

// another set of methods for getting the required rate is written here, can be writter in a different file for readability

// as we're going to call this function outside of this body, it is preferable to capitalize

func (er *ExchangeRates) GetRae(base, dest string) (float64, error) {
	b, ok := er.rates[base] // returns the rate for a specified base string, eg er.rates["EUR"] fetches or stores 1 in b, ok has a job of returning a boolean value, if the value is present then it holds "true"
	if !ok {
		return 0, fmt.Errorf("rate not found for currency %s", base)
	}

	d, ok := er.rates[dest]
	if !ok {
		return 0, fmt.Errorf("rate not found for currency %s", dest)
	}

	return d / b, nil // dest rate / base rate for calculation

}

// monitor rates function, take time.Duration type as input and a channel of type struct as output
func (er *ExchangeRates) MonitorRates(xd time.Duration) chan struct{} {
	r := make(chan struct{}) // creating a channel using make inbuilt function

	// concurrent execution
	go func() {
		ticker := time.NewTicker(xd) // creating a new ticker

		// ticker fetches data at equal interval of ticks

		// neverending for loop
		for {

			// whichever channel comes first, execute that statement and exit
			select {
			case <-ticker.C:

				// just add a random difference to rate and return it
				// this simulates the fluctiuations in original currency rates daily
				for x, v := range er.rates { // x string(key), v(value(rate), float)

					// for now lets change 10% of the original value
					cg := (rand.Float64() / 10)

					// positive or negative change
					dir := rand.Intn(1)

					if dir == 0 {
						// new value will be min 90% of old
						cg = 1 - cg
					} else {
						cg = 1 + cg
					}

					// modify the rate
					er.rates[x] = v * cg
				}

				// notify updates, this will block unless there is a listener on the other end
				r <- struct{}{} // returning empty struct depicting the execution is complete
			}
		}
	}()

	return r
}
