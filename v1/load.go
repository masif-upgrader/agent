package v1

type Load struct {
	Query     [3]float64 `json:"query"`
	Install   [3]float64 `json:"install"`
	Update    [3]float64 `json:"update"`
	Configure [3]float64 `json:"configure"`
	Remove    [3]float64 `json:"remove"`
	Purge     [3]float64 `json:"purge"`
	Error     [3]float64 `json:"error"`
}
