package interfaces

type UpdateStockRequest struct {
	PortName string      `json:"portName"`
	Stock    []StockInfo `json:"stock"`
}
