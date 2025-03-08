package exceptions

import "github.com/niflheimdevs/backend/internal/enums"

type Exception struct {
	Tag    enums.GeneralError
	Errors []enums.SpecificError
}
