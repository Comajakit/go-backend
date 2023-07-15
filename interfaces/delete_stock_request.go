package interfaces

type DeleteStockRequest struct {
	PortName string   `json:"portName"`
	Stock    []string `json:"stock"`
}
