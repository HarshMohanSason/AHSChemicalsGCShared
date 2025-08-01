package qbmodels

import (
	"time"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/quickbooks"
)

type QBCustomersResponse struct {
	QueryResponse struct {
		Customer []QBCustomer `json:"Customer"`
	} `json:"QueryResponse"`
}

type QBCustomer struct {
	ID                      string               `json:"Id"`
	SyncToken               string               `json:"SyncToken"`
	DisplayName             string               `json:"DisplayName"`
	GivenName               string               `json:"GivenName,omitempty"`
	MiddleName              string               `json:"MiddleName,omitempty"`
	FamilyName              string               `json:"FamilyName,omitempty"`
	CompanyName             string               `json:"CompanyName,omitempty"`
	PrimaryEmailAddr        *QBCustomerEmail     `json:"PrimaryEmailAddr,omitempty"`
	PrimaryPhone            *QBCustomerPhone     `json:"PrimaryPhone,omitempty"`
	Mobile                  *QBCustomerPhone     `json:"Mobile,omitempty"`
	AlternatePhone          *QBCustomerPhone     `json:"AlternatePhone,omitempty"`
	Fax                     *QBCustomerPhone     `json:"Fax,omitempty"`
	BillAddr                *QBCustomerAddress   `json:"BillAddr,omitempty"`
	ShipAddr                *QBCustomerAddress   `json:"ShipAddr,omitempty"`
	Notes                   string               `json:"Notes,omitempty"`
	Balance                 float64              `json:"Balance,omitempty"`
	BalanceWithJobs         float64              `json:"BalanceWithJobs,omitempty"`
	Active                  bool                 `json:"Active"`
	Job                     bool                 `json:"Job,omitempty"`
	OpenBalanceDate         string               `json:"OpenBalanceDate,omitempty"`
	CustomerTypeRef         *QBCustomerRef       `json:"CustomerTypeRef,omitempty"`
	ParentRef               *QBCustomerRef       `json:"ParentRef,omitempty"`
	Taxable                 bool                 `json:"Taxable,omitempty"`
	TaxExemptionReasonId    string               `json:"TaxExemptionReasonId,omitempty"`
	ResaleNum               string               `json:"ResaleNum,omitempty"`
	TaxCodeRef              *QBCustomerRef       `json:"TaxCodeRef,omitempty"`
	PreferredDeliveryMethod string               `json:"PreferredDeliveryMethod,omitempty"`
	WebAddr                 *Web                 `json:"WebAddr,omitempty"`
	MetaData                *quickbooks.MetaData `json:"MetaData,omitempty"`
	PrimaryTaxIdentifier    string               `json:"PrimaryTaxIdentifier,omitempty"`
	Level                   int                  `json:"Level,omitempty"`
	CustomerBalance         float64              `json:"CustomerBalance,omitempty"`
	CustomerBalanceWithJobs float64              `json:"CustomerBalanceWithJobs,omitempty"`
}

type QBCustomerRef struct {
	Value string `json:"value"`
	Name  string `json:"name,omitempty"`
}

type QBCustomerEmail struct {
	Address string `json:"Address"`
}

type QBCustomerPhone struct {
	FreeFormNumber string `json:"FreeFormNumber"`
}

type QBCustomerAddress struct {
	Id                     string `json:"Id,omitempty"`
	Line1                  string `json:"Line1,omitempty"`
	Line2                  string `json:"Line2,omitempty"`
	Line3                  string `json:"Line3,omitempty"`
	City                   string `json:"City,omitempty"`
	Country                string `json:"Country,omitempty"`
	CountrySubDivisionCode string `json:"CountrySubDivisionCode,omitempty"` // e.g. "CA"
	PostalCode             string `json:"PostalCode,omitempty"`
	Lat                    string `json:"Lat,omitempty"`
	Long                   string `json:"Long,omitempty"`
}

type Web struct {
	URI string `json:"URI"`
}

func (qb *QBCustomer) SetEmailInCustomer(c *models.Customer) {
	if qb.PrimaryEmailAddr == nil{
		return
	}
	c.Email = qb.PrimaryEmailAddr.Address
}

func (qb *QBCustomer) SetPhoneInCustomer(c *models.Customer) {
	if qb.PrimaryPhone == nil{
		return
	}
	c.Phone = qb.PrimaryPhone.FreeFormNumber
}

func (qb *QBCustomer) SetAddressInCustomer(c *models.Customer) {
	if qb.BillAddr == nil{
		return
	}
	c.Address1 = qb.BillAddr.Line1
	c.City = qb.BillAddr.City
	c.State = qb.BillAddr.CountrySubDivisionCode
	c.Zip = qb.BillAddr.PostalCode
	c.Country = qb.BillAddr.Country
}

func (qb *QBCustomer) MapToCustomer() *models.Customer{

	customer := &models.Customer{
		ID: qb.ID,
		IsActive: qb.Active,
		Name: qb.DisplayName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	qb.SetEmailInCustomer(customer)
	qb.SetPhoneInCustomer(customer)
	qb.SetAddressInCustomer(customer)
	return customer
}