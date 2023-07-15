package interfaces

type CreateStrategyRequest struct {
	PortName     string              `json:"portName"`
	StrategyName string              `json:"strategyName"`
	Theme        []map[string]string `json:"theme"`
}
