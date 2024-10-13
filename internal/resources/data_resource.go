package resources

type Response struct {
	Data any `json:"data"`
}

type MapResponse struct {
	Data map[string]any `json:"data"`
}
