package dto

type GetTagDto struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
}

type RecieveTagDTO struct {
	ID    int `json:"id" validate:"required,numeric"`
	Level int `json:"level" validate:"required,regexp=^-?\\d+$"`
}

func (rt *RecieveTagDTO) IsEqualToGetTagDTO(gt *GetTagDto) bool {
	return rt.Level == gt.Level &&
		rt.ID == gt.ID
}
