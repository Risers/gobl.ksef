package ksef

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
)

type Address struct {
	CountryCode string `xml:"KodKraju"`
	AddressL1   string `xml:"AdresL1"`
	AddressL2   string `xml:"AdresL2"`
}

type Seller struct {
	NIP     string          `xml:"DaneIdentyfikacyjne>NIP"`
	Name    string          `xml:"DaneIdentyfikacyjne>Nazwa"`
	Address *Address        `xml:"Adres"`
	Contact *ContactDetails `xml:"DaneKontaktowe,omitempty"`
}

type ContactDetails struct {
	Phone string `xml:"Telefon,omitempty"`
	Email string `xml:"Email,omitempty"`
}

type Buyer struct {
	NIP string `xml:"DaneIdentyfikacyjne>NIP,omitempty"`
	// or
	UECode      string `xml:"DaneIdentyfikacyjne>KodUE,omitempty"` //TODO
	UEVatNumber string `xml:"DaneIdentyfikacyjne>NrVatUE,omitempty"`
	// or
	CountryCode string `xml:"DaneIdentyfikacyjne>KodKraju,omitempty"`
	IdNumber    string `xml:"DaneIdentyfikacyjne>NrId,omitempty"`
	// or
	NoId int `xml:"DaneIdentyfikacyjne>BrakID,omitempty"`

	Name    string          `xml:"DaneIdentyfikacyjne>Nazwa,omitempty"`
	Address *Address        `xml:"Adres,omitempty"`
	Contact *ContactDetails `xml:"DaneKontaktowe,omitempty"`
}

func NewAddress(address *org.Address) *Address {
	adres := &Address{
		CountryCode: string(address.Country),
		AddressL1:   addressLine1(address),
		AddressL2:   addressLine2(address),
	}

	return adres
}

func NameToString(name org.Name) string {
	return name.Prefix + nameMaybe(name.Given) +
		nameMaybe(name.Middle) + nameMaybe(name.Surname) +
		nameMaybe(name.Surname2) + nameMaybe(name.Suffix)
}

func NewSeller(supplier *org.Party) *Seller {
	var name string
	if supplier.Name != "" {
		name = supplier.Name
	} else {
		name = NameToString(supplier.People[0].Name)
	}
	seller := &Seller{
		Address: NewAddress(supplier.Addresses[0]),
		NIP:     string(supplier.TaxID.Code),
		Name:    name,
	}
	if len(supplier.Telephones) > 0 {
		seller.Contact = &ContactDetails{
			Phone: supplier.Telephones[0].Number,
		}
	}
	if len(supplier.Emails) > 0 {
		if seller.Contact == nil {
			seller.Contact = &ContactDetails{}
		}
		seller.Contact.Email = supplier.Emails[0].Address
	}

	return seller
}

func NewBuyer(customer *org.Party) *Buyer {

	buyer := &Buyer{
		Name: customer.Name,
		NIP:  string(customer.TaxID.Code),
	}

	if customer.TaxID.Country == l10n.PL {
		buyer.NIP = string(customer.TaxID.Code)
	} else {
		if len(customer.TaxID.Code) > 0 {
			buyer.IdNumber = string(customer.TaxID.Code)
			buyer.CountryCode = string(customer.TaxID.Country)
		} else {
			buyer.NoId = 1
		}
	}
	// TODO NrVatUE and UECode if applicable

	if len(customer.Addresses) > 0 {
		buyer.Address = NewAddress(customer.Addresses[0])
	}

	if len(customer.Telephones) > 0 {
		buyer.Contact = &ContactDetails{
			Phone: customer.Telephones[0].Number,
		}
	}
	if len(customer.Emails) > 0 {
		if buyer.Contact == nil {
			buyer.Contact = &ContactDetails{}
		}
		buyer.Contact.Email = customer.Emails[0].Address
	}

	return buyer
}

func addressLine1(address *org.Address) string {
	if address.PostOfficeBox != "" {
		return address.PostOfficeBox
	}

	return address.Street +
		", " + address.Number +
		addressMaybe(address.Block) +
		addressMaybe(address.Floor) +
		addressMaybe(address.Door)
}

func addressLine2(address *org.Address) string {
	return address.Code + ", " + address.Locality
}

func addressMaybe(element string) string {
	if element != "" {
		return ", " + element
	}
	return ""
}

func nameMaybe(element string) string {
	if element != "" {
		return " " + element
	}
	return ""
}
