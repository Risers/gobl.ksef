// Package ksef implements the conversion from GOBL to FA_VAT XML
package ksef

import (
	"encoding/xml"
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
)

const (
	XSINamespace    = "http://www.w3.org/2001/XMLSchema-instance"
	XSDNamespace    = "http://www.w3.org/2001/XMLSchema"
	XMLNamespace    = "http://crd.gov.pl/wzor/2023/06/29/12648/"
	RootElementName = "Faktura"
)

// Invoice is a pseudo-model for containing the XML document being created
type Invoice struct {
	XMLName      xml.Name
	XSINamespace string  `xml:"xmlns:xsi,attr"`
	XSDNamespace string  `xml:"xmlns:xsd,attr"`
	XMLNamespace string  `xml:"xmlns,attr"`
	Header       *Header `xml:"Naglowek"`
	Seller       *Seller `xml:"Podmiot1"`
	Buyer        *Buyer  `xml:"Podmiot2"`
	ThirdParty   *Buyer  `xml:"Podmiot3,omitempty"` // third party
	Inv          *Inv    `xml:"Fa"`
	Footer       *Footer `xml:"Stopka,omitempty"`
}

// NewDocument converts a GOBL envelope into a FA_VAT document
func NewDocument(env *gobl.Envelope) (*Invoice, error) {
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return nil, fmt.Errorf("invalid type %T", env.Document)
	}

	invoice := &Invoice{
		XMLName:      xml.Name{Local: RootElementName},
		XSINamespace: XSINamespace,
		XSDNamespace: XSDNamespace,
		XMLNamespace: XMLNamespace,

		Header: NewHeader(inv),
		Seller: NewSeller(inv.Supplier),
		Buyer:  NewBuyer(inv.Customer),
		Inv:    NewInv(inv),
		Footer: NewFooter(inv),
	}

	return invoice, nil
}

// Bytes returns the XML representation of the document in bytes
func (d *Invoice) Bytes() ([]byte, error) {
	bytes, err := xml.MarshalIndent(d, "", "  ")
	if err != nil {
		return nil, err
	}

	return append([]byte(xml.Header), bytes...), nil
}
