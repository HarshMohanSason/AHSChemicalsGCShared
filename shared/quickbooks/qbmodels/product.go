package qbmodels

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/models"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared/quickbooks"
)

var (
	Brands = map[string]struct{}{
		"microtech": {},
		"problend":  {},
	}
	slugSanitizer = regexp.MustCompile(`[^a-z0-9]+`)
)

type QBItemsResponse struct {
	QueryResponse struct {
		Item []QBItem `json:"Item"`
	} `json:"QueryResponse"`
}

type QBItem struct {
	ID                 string               `json:"Id"`
	Name               string               `json:"Name"`
	SKU                string               `json:"Sku,omitempty"`
	SyncToken          string               `json:"SyncToken,omitempty"`
	Description        string               `json:"Description,omitempty"`
	Active             bool                 `json:"Active"`
	FullyQualifiedName string               `json:"FullyQualifiedName,omitempty"`
	UnitPrice          float64              `json:"UnitPrice,omitempty"`
	Type               string               `json:"Type"` // "Inventory", "NonInventory", "Service", "OtherCharge"
	TrackQtyOnHand     bool                 `json:"TrackQtyOnHand,omitempty"`
	QtyOnHand          float64              `json:"QtyOnHand,omitempty"`
	InvStartDate       string               `json:"InvStartDate,omitempty"`
	PurchaseCost       float64              `json:"PurchaseCost,omitempty"`
	PurchaseDesc       string               `json:"PurchaseDesc,omitempty"`
	IncomeAccountRef   *QBItemRef           `json:"IncomeAccountRef,omitempty"`
	ExpenseAccountRef  *QBItemRef           `json:"ExpenseAccountRef,omitempty"`
	AssetAccountRef    *QBItemRef           `json:"AssetAccountRef,omitempty"`
	TaxCodeRef         *QBItemRef           `json:"TaxCodeRef,omitempty"`
	ParentRef          *QBItemRef           `json:"ParentRef,omitempty"`
	SalesTaxIncluded   bool                 `json:"SalesTaxIncluded,omitempty"`
	Taxable            bool                 `json:"Taxable,omitempty"`
	MetaData           *quickbooks.MetaData `json:"MetaData,omitempty"`
	CustomField        []QBItemCustomField  `json:"CustomField,omitempty"`
}

type QBItemRef struct {
	Value string `json:"value"`
	Name  string `json:"name,omitempty"`
}

type QBItemCustomField struct {
	Type  string `json:"Type"`        // "StringType", "DateType", "NumberType", etc.
	Name  string `json:"Name"`        // Custom field label
	Value string `json:"StringValue"` // Field value
}

// parseName parses product name into firestore product
//
// Quickbooks does not allow storing additional properties like product brand.
// So each product is inputted as one entire string: `BrandName - ProductName` in quickbooks manually. We parse
// it from the delimeter '-' to get the appropriate brand and product name
func (qb *QBItem) parseNameInto(product *models.Product) {
	splitString := strings.SplitN(qb.Name, "-", 2)
	if len(splitString) != 2 {
		product.Name = qb.Name
		return
	}
	parsedBrand := strings.TrimSpace(splitString[0])
	parsedProductName := strings.TrimSpace(splitString[1])

	if _, ok := Brands[strings.ToLower(parsedBrand)]; ok {
		product.Brand = parsedBrand
		product.Name = parsedProductName
	} else {
		product.Name = qb.Name
	}
}

// parseSKUInto parses product SKU into Firestore Product
//
// Quickbooks does not allow storing additional properties like size unit, size or pack of.
// So each product is inputted as one entire string: `SKU - Size - SizeUnit - PackOf` in quickbooks manually. We parse
// it from the delimeter '-' to get the appropriate SKU, Size, SizeUnit and PackOf. If the delimeter is not present or
// there is a malformed string, we set the default values as "QBSKU - 0.0 - GAL - 1"
func (qb *QBItem) parseSKUInto(product *models.Product) {
	splitString := strings.SplitN(qb.SKU, "-", 4)
	if len(splitString) == 4 {
		size, err := strconv.ParseFloat(splitString[1], 64)
		if err != nil {
			size = 0.0
		}
		packOf, err := strconv.Atoi(splitString[3])
		if err != nil {
			packOf = 0
		}

		product.SKU = splitString[0]
		product.Size = size
		product.SizeUnit = splitString[2]
		product.PackOf = packOf
	} else {
		product.SKU = qb.SKU
		product.Size = 0.0
		product.SizeUnit = "GAL"
		product.PackOf = 1
	}
}

// parseSlugAndNameKeyInto parses product name into firestore product
//
// Creates a unique slug with a namekey for each product in quickbooks.
func (qb *QBItem) parseSlugAndNameKeyInto(product *models.Product) {
	slug := strings.ToLower(product.Name)
	slug = slugSanitizer.ReplaceAllString(slug, "-") //Replace everthing with a hyphen
	slug = strings.Trim(slug, "-")                   //Cleaning any leading or trailing spaces if any present

	product.NameKey = slug
	product.Slug = fmt.Sprintf("%s-%s", slug, qb.ID)
}

func (qb *QBItem) parseCategoryInto(product *models.Product) {
	if qb.ParentRef == nil {
		return
	}
	product.Category = qb.ParentRef.Name
}

func (qb *QBItem) MapToProduct() *models.Product {
	product := &models.Product{
		ID:            qb.ID,
		IsActive:      qb.Active,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Hazardous:     false,
		Price:         qb.UnitPrice,
		PurchasePrice: qb.PurchaseCost,
		Desc:          qb.Description,
	}
	qb.parseNameInto(product)
	qb.parseSKUInto(product)
	qb.parseSlugAndNameKeyInto(product)
	qb.parseCategoryInto(product)
	return product
}
