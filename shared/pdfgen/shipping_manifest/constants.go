package shipping_manifest

//Shipping Manifest constant vars for pdf generation
var (
	ShippingManifestHeaders        = []string{"UNITS", "HM", "TYPE CONTAINER", "DESCRIPTION", "CLASS", "SKU", "NET WEIGHT", "GROSS WEIGHT NHM", "GROSS WEIGHT HM"}
	ShippingManifestTableColWidths = []float64{13, 10, 20, 40, 15, 30, 20.6, 20.6, 20.6}
)