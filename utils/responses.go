package utils

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Errors  any    `json:"errors,omitempty"`
	Results any    `json:"results,omitempty"`
	Page any `json:"page,omitempty"`
	PerPage any `json:"perPage,omitempty"`
}
