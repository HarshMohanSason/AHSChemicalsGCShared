package utils

func CalculateShippingTableCellWidths(cellWidths []float64) float64 {
	var totalWidth float64 = 0.0
	for i := range cellWidths {
		totalWidth += cellWidths[i]
	}
	return totalWidth
}