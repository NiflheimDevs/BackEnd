package exceptions

import "github.com/niflheimdevs/backend/internal/enums"

type Exception struct {
	Tag    enums.GeneralError    `json:"tag"`
	Errors []enums.SpecificError `json:"errors"`
}

func (e *Exception) AddError(se enums.SpecificError) {
	e.Errors = append(e.Errors, se)
}
