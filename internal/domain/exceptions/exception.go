package exceptions

type Exception struct {
	Tag    GeneralError    `json:"tag"`
	Errors []SpecificError `json:"errors"`
}

func (e *Exception) AddError(se SpecificError) {
	e.Errors = append(e.Errors, se)
}
