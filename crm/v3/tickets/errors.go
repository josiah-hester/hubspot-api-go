package tickets

import "strings"

type ObjectError struct {
	Message     string              `json:"message" required:"yes"`
	SubCategory string              `json:"subCategory"`
	Code        string              `json:"code"`
	In          string              `json:"in"`
	Context     map[string][]string `json:"context"`
}

type BatchError struct {
	Context     map[string][]string `json:"context" required:"yes"`
	Links       map[string]string   `json:"links" required:"yes"`
	Category    string              `json:"category" required:"yes"`
	Message     string              `json:"message" required:"yes"`
	Errors      []ObjectError       `json:"errors" required:"yes"`
	Status      string              `json:"status" required:"yes"`
	SubCategory any                 `json:"subCategory"`
	ID          string              `json:"id"`
}

func (e *BatchError) Error() string {
	if e.Message != "" && len(e.Errors) == 0 {
		return e.Message
	}

	var sb strings.Builder
	if e.Message != "" {
		sb.WriteString(e.Message)
	} else {
		sb.WriteString("batch error")
	}

	if len(e.Errors) > 0 {
		sb.WriteString(": ")
		for i, err := range e.Errors {
			if i > 0 {
				sb.WriteString("; ")
			}
			sb.WriteString(err.Message)
		}
	}

	return sb.String()
}
