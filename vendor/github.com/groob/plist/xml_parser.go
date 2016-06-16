package plist

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

// xmlParser uses xml.Decoder to parse an xml plist into the corresponding plistValues
type xmlParser struct {
	*xml.Decoder
}

// newXMLParser returns a new xmlParser
func newXMLParser(r io.Reader) *xmlParser {
	return &xmlParser{xml.NewDecoder(r)}
}

func (p *xmlParser) parseDocument(start *xml.StartElement) (*plistValue, error) {
	if start == nil {
		for {
			tok, err := p.Token()
			if err != nil {
				return nil, err
			}
			if t, ok := tok.(xml.StartElement); ok {
				start = &t
				break
			}
		}
	}
	return p.parseXMLElement(start)
}

func (p *xmlParser) parseXMLElement(element *xml.StartElement) (*plistValue, error) {
	switch element.Name.Local {
	case "plist":
		return p.parsePlist(element)
	case "dict":
		return p.parseDict(element)
	case "string":
		return p.parseString(element)
	case "true", "false":
		return p.parseBoolean(element)
	case "array":
		return p.parseArray(element)
	case "real":
		return p.parseReal(element)
	case "integer":
		return p.parseInteger(element)
	case "data":
		return p.parseData(element)
	case "date":
		return p.parseDate(element)
	default:
		return nil, fmt.Errorf("plist: Unknown plist element %s", element.Name.Local)
	}
}

func (p *xmlParser) parsePlist(element *xml.StartElement) (*plistValue, error) {
	for {
		token, err := p.Token()
		if err != nil {
			return nil, err
		}
		if el, ok := token.(xml.EndElement); ok && el.Name.Local == "plist" {
			break
		}
		if el, ok := token.(xml.StartElement); ok {
			return p.parseXMLElement(&el)
		}
	}
	return nil, errors.New("plist: Invalid plist")
}

func (p *xmlParser) parseDict(element *xml.StartElement) (*plistValue, error) {
	var key *string
	var subvalues = make(map[string]*plistValue)
	for {
		token, err := p.Token()
		if err != nil {
			return nil, err
		}
		if el, ok := token.(xml.EndElement); ok && el.Name.Local == "dict" {
			break
		}
		if el, ok := token.(xml.StartElement); ok {
			if el.Name.Local == "key" {
				var k string
				if err := p.DecodeElement(&k, &el); err != nil {
					return nil, err
				}
				key = &k
				continue
			}
			if key == nil {
				return nil, errors.New("plist: missing key in dict")
			}
			subvalues[*key], err = p.parseXMLElement(&el)
			if err != nil {
				return nil, err
			}
			key = nil
		}
	}
	return &plistValue{Dictionary, &dictionary{m: subvalues}}, nil
}

func (p *xmlParser) parseString(element *xml.StartElement) (*plistValue, error) {
	var value string
	if err := p.DecodeElement(&value, element); err != nil {
		return nil, err
	}
	return &plistValue{String, value}, nil
}

func (p *xmlParser) parseBoolean(element *xml.StartElement) (*plistValue, error) {
	if err := p.Skip(); err != nil {
		return nil, err
	}
	plistBoolean := element.Name.Local == "true"
	return &plistValue{Boolean, plistBoolean}, nil
}

func (p *xmlParser) parseArray(element *xml.StartElement) (*plistValue, error) {
	var subvalues []*plistValue
	for {
		token, err := p.Token()
		if err != nil {
			return nil, err
		}
		if el, ok := token.(xml.EndElement); ok && el.Name.Local == "array" {
			break
		}
		if el, ok := token.(xml.StartElement); ok {
			subv, err := p.parseXMLElement(&el)
			if err != nil {
				return nil, err
			}
			subvalues = append(subvalues, subv)
		}
	}
	return &plistValue{Array, subvalues}, nil
}

func (p *xmlParser) parseReal(element *xml.StartElement) (*plistValue, error) {
	var n float64
	if err := p.DecodeElement(&n, element); err != nil {
		return nil, err
	}
	return &plistValue{Real, sizedFloat{n, 64}}, nil
}

func (p *xmlParser) parseInteger(element *xml.StartElement) (*plistValue, error) {
	var n uint64
	if err := p.DecodeElement(&n, element); err != nil {
		return nil, err
	}
	return &plistValue{Integer, signedInt{n, false}}, nil
}

func (p *xmlParser) parseData(element *xml.StartElement) (*plistValue, error) {
	replacer := strings.NewReplacer("\t", "", "\n", "", " ", "", "\r", "")
	var data []byte
	if err := p.DecodeElement(&data, element); err != nil {
		return nil, err
	}
	str := replacer.Replace(string(data))
	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	data = []byte(decoded)
	return &plistValue{Data, data}, nil
}

func (p *xmlParser) parseDate(element *xml.StartElement) (*plistValue, error) {
	var date time.Time
	if err := p.DecodeElement(&date, element); err != nil {
		return nil, err
	}
	return &plistValue{Date, date}, nil
}
