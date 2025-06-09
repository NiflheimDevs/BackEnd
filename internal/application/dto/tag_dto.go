package dto

type GetTagDto struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
}

type RecieveTagDTO struct {
	ID    int `json:"id" validate:"numeric"`
	Level int `json:"level" validate:"required,gt=-2,lt=6"`
}

func (rt *RecieveTagDTO) IsEqualToGetTagDTO(gt *GetTagDto) bool {
	return rt.Level == gt.Level &&
		rt.ID == gt.ID
}
