package entity

type Owner struct {
	Username string `json:"username" validate:"required"`
	Folder   string `json:"folder" validate:"required"`
	Client   string `json:"client" validate:"required"`
}
