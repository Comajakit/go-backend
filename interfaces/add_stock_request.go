package interfaces

type AddStockRequest struct {
	PortName string      `json:"portName"`
	Stock    []StockInfo `json:"stock"`
}

type StockInfo struct {
	Symbol       string `json:"symbol"`
	Volume       int    `json:"volume"`
	AveragePrice string `json:"averagePrice"`
	DivPerShare  string `json:"divPerShare"`
	Type         string `json:"type"`
}