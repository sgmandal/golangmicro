package data

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-hclog"
)

func TestNewrates(t *testing.T) {
	tr, err := NewRates(hclog.Default())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%#v", tr.rates)
}
