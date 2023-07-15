package interfaces

type UpdateStrategyRequest struct {
	PortName     string              `json:"portName"`
	StrategyName string              `json:"strategyName"`
	Theme        []map[string]string `json:"theme"`
}
