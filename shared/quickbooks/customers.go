package quickbooks

import "github.com/HarshMohanSason/AHSChemicalsGCShared/shared"

type Phone struct {
	FreeFormNumber string `json:"FreeFormNumber,omitempty"`
}

type Email struct {
	Address string `json:"Address,omitempty"`
}

type Address struct {
	Id                     string `json:"Id,omitempty"`
	Line1                  string `json:"Line1,omitempty"`
	Line2                  string `json:"Line2,omitempty"`
	City                   string `json:"City,omitempty"`
	CountrySubDivisionCode string `json:"CountrySubDivisionCode,omitempty"`
	PostalCode             string `json:"PostalCode,omitempty"`
	Lat                    string `json:"Lat,omitempty"`
	Long                   string `json:"Long,omitempty"`
}
type Customer struct {
	Id                      string            `json:"Id"`
	SyncToken               string            `json:"SyncToken,omitempty"`
	DisplayName             string            `json:"DisplayName"`
	GivenName               string            `json:"GivenName,omitempty"`
	MiddleName              string            `json:"MiddleName,omitempty"`
	FamilyName              string            `json:"FamilyName,omitempty"`
	CompanyName             string            `json:"CompanyName,omitempty"`
	FullyQualifiedName      string            `json:"FullyQualifiedName,omitempty"`
	PrintOnCheckName        string            `json:"PrintOnCheckName,omitempty"`
	PrimaryPhone            *Phone            `json:"PrimaryPhone,omitempty"`
	PrimaryEmailAddr        *Email            `json:"PrimaryEmailAddr,omitempty"`
	BillAddr                *Address          `json:"BillAddr,omitempty"`
	ShipAddr                *Address          `json:"ShipAddr,omitempty"`
	Taxable                 bool              `json:"Taxable,omitempty"`
	Job                     bool              `json:"Job,omitempty"`
	Balance                 float64           `json:"Balance,omitempty"`
	BalanceWithJobs         float64           `json:"BalanceWithJobs,omitempty"`
	Active                  bool              `json:"Active,omitempty"`
	MetaData                *shared.MetaData  `json:"MetaData,omitempty"`
	PreferredDeliveryMethod string            `json:"PreferredDeliveryMethod,omitempty"`
	CurrencyRef             *shared.Reference `json:"CurrencyRef,omitempty"`
	V4IDPseudonym           string            `json:"V4IDPseudonym,omitempty"`
	IsProject               bool              `json:"IsProject,omitempty"`
	ClientEntityId          string            `json:"ClientEntityId,omitempty"`
	Notes                   string            `json:"Notes,omitempty"`
}
