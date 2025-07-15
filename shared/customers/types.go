package customers

import "github.com/HarshMohanSason/AHSChemicalsGCShared/shared"

type Email struct {
	Address string `json:"Address"`
}

type Phone struct {
	FreeFormNumber string `json:"FreeFormNumber"`
}

type Address struct {
	Id                     string `json:"Id,omitempty"`
	Line1                  string `json:"Line1,omitempty"`
	Line2                  string `json:"Line2,omitempty"`
	City                   string `json:"City,omitempty"`
	CountrySubDivisionCode string `json:"CountrySubDivisionCode,omitempty"` // State code like "CA"
	PostalCode             string `json:"PostalCode,omitempty"`
	Country                string `json:"Country,omitempty"`
	Lat                    string `json:"Lat,omitempty"`
	Long                   string `json:"Long,omitempty"`
}

type Ref struct {
	Value string `json:"value"`
	Name  string `json:"name,omitempty"`
}

type Customer struct {
	ID                   string           `json:"Id"`
	SyncToken            string           `json:"SyncToken"`
	Title                string           `json:"Title,omitempty"`
	GivenName            string           `json:"GivenName,omitempty"`
	MiddleName           string           `json:"MiddleName,omitempty"`
	FamilyName           string           `json:"FamilyName,omitempty"`
	DisplayName          string           `json:"DisplayName"`
	FullyQualifiedName   string           `json:"FullyQualifiedName,omitempty"`
	CompanyName          string           `json:"CompanyName,omitempty"`
	PrimaryEmailAddr     *Email           `json:"PrimaryEmailAddr,omitempty"`
	PrimaryPhone         *Phone           `json:"PrimaryPhone,omitempty"`
	Mobile               *Phone           `json:"Mobile,omitempty"`
	Fax                  *Phone           `json:"Fax,omitempty"`
	AlternatePhone       *Phone           `json:"AlternatePhone,omitempty"`
	BillAddr             *Address         `json:"BillAddr,omitempty"`
	ShipAddr             *Address         `json:"ShipAddr,omitempty"`
	Notes                string           `json:"Notes,omitempty"`
	Active               bool             `json:"Active"`
	Balance              float64          `json:"Balance,omitempty"`
	BalanceWithJobs      float64          `json:"BalanceWithJobs,omitempty"`
	OpenBalanceDate      string           `json:"OpenBalanceDate,omitempty"` // Date string like "2023-09-01"
	CustomerTypeRef      *Ref             `json:"CustomerTypeRef,omitempty"`
	PaymentMethodRef     *Ref             `json:"PaymentMethodRef,omitempty"`
	TermsRef             *Ref             `json:"SalesTermRef,omitempty"`
	Taxable              bool             `json:"Taxable"`
	TaxExemptionReasonId string           `json:"TaxExemptionReasonId,omitempty"`
	TaxExemptionReason   string           `json:"ExemptionReasonCode,omitempty"`
	ResaleNum            string           `json:"ResaleNum,omitempty"`
	Level                string           `json:"Level,omitempty"`
	PreferredDelivery    string           `json:"PreferredDeliveryMethod,omitempty"` // "Email", "Print", "None"
	MetaData             *shared.MetaData `json:"MetaData,omitempty"`
}
