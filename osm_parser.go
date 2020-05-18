package main

import (
	"encoding/xml"
	"io"
	"os"
	"strings"
)

// DecodeFile Décode un fichier OpenStreetMap
func DecodeFile(fileName string) (*Map, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Decode(file)
}

// DecodeString Décode un string
func DecodeString(data string) (*Map, error) {
	return Decode(strings.NewReader(data))
}

// Decode Décode un lecteur
func Decode(reader io.Reader) (*Map, error) {
	var (
		o   = new(Map)
		err error
	)

	decoder := xml.NewDecoder(reader)
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}

		switch typedToken := token.(type) {
		case xml.StartElement:
			switch typedToken.Name.Local {
			case "bounds":
				var b Bounds
				err = decoder.DecodeElement(&b, &typedToken)
				if err != nil {
					return nil, err
				}
				o.Bounds = b

			case "node":
				var n Node
				err = decoder.DecodeElement(&n, &typedToken)
				if err != nil {
					return nil, err
				}
				o.Nodes = append(o.Nodes, n)

			case "way":
				var w Way
				err = decoder.DecodeElement(&w, &typedToken)
				if err != nil {
					return nil, err
				}
				o.Ways = append(o.Ways, w)

			}
		}
	}
	return o, nil
}
