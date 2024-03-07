package interfaces

type SimplifiedStock struct {
	Total             float64 `json:"total"`
	DivPerShare       float64 `json:"divPerShare"`
	DivInPercent      float64 `json:"divInPercent"`
	ExpectedDivReturn float64 `json:"expectedDivReturn"`
	PercentageInPort  float64 `json:"percentageInPort"`
	DivPercentPort    float64 `json:"divPercentPort"`
	StockSymbol       string  `json:"stockSymbol"`
	Volume            uint    `json:"volume"`
	AveragePrice      float64 `json:"averagePrice"`
	StockType         string  `json:"stockType"`
}

type GetPortResponse struct {
	PortName string            `json:"portName"`
	Stock    []SimplifiedStock `json:"stock"`
}
