package ksef

import "github.com/invopop/gobl/bill"

type Line struct {
	LineNumber              int    `xml:"NrWierszaFa"`
	Name                    string `xml:"P_7"`
	Measure                 string `xml:"P_8A"`
	Quantity                string `xml:"P_8B"`
	NetUnitPrice            string `xml:"P_9A"`
	UnitDiscount            string `xml:"P_10"`
	NetPriceTotal           string `xml:"P_11"`
	TaxRate                 string `xml:"P_12"`
	ExciseDuty              string `xml:"KwotaAkcyzy,omitempty"`
	SpecialGoodsCode        string `xml:"GTU,omitempty"` // values GTU_1 to GTU_13
	OSSTaxRate              string `xml:"P_12_XII,omitempty"`
	Attachment15GoodsMarker string `xml:"P_12_Zal_15,omitempty"`
	Procedure               string `xml:"Procedura,omitempty"`
	BeforeCorrectionMarker  string `xml:"StanPrzed,omitempty"`
}

func NewLine(line *bill.Line) *Line {
	Line := &Line{
		LineNumber:    line.Index,
		Name:          line.Item.Name,
		Measure:       string(line.Item.Unit.UNECE()),
		NetUnitPrice:  line.Item.Price.String(),
		Quantity:      line.Quantity.String(),
		UnitDiscount:  line.Total.Subtract(line.Sum).Divide(line.Quantity).String(), // not sure if there should be some rescale or matchPrecision
		NetPriceTotal: line.Total.String(),
		TaxRate:       line.Taxes[0].Percent.Rescale(2).StringWithoutSymbol(),
	}

	return Line
}

func NewLines(lines []*bill.Line) []*Line {
	var Lines []*Line

	for _, line := range lines {
		Lines = append(Lines, NewLine(line))
	}

	return Lines
}
