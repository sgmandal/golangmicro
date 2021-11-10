package data

import (
	"encoding/json"
	"io"
)

// this function serializes the given interface into a string based JSON format
// i interface is the input data
// interface is used such that any data can be expected here

// which makes these two kinda a universal implementation for Json conversion
func ToJSON(i interface{}, e io.Writer) error {
	w := json.NewEncoder(e) // creating an encoder instance wrt io.writer
	return w.Encode(i)      // passing the data
}

func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r) // creating a decoder instance wrt io.reader
	return d.Decode(i)      // passing the data
}
