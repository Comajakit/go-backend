package types

type StockSummary struct {
	Theme   string   `json:"theme"`
	Symbol  []string `json:"symbol"`
	Target  string   `json:"target"`
	Current string   `json:"current"`
}

type SummaryResponse struct {
	PortName     string                    `json:"portName"`
	StrategyName string                    `json:"strategyName"`
	Summary      map[string][]StockSummary `json:"summary"`
}
