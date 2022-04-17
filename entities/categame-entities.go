package entities

type Model_categame struct {
	Categame_id      string `json:"categame_id"`
	Categame_name    string `json:"categame_name"`
	Categame_display int    `json:"categame_display"`
	Categame_status  string `json:"categame_status"`
	Categame_create  string `json:"categame_create"`
	Categame_update  string `json:"categame_update"`
}
type Controller_categame struct {
	Master string `json:"master" validate:"required"`
}
type Controller_categamesave struct {
	Sdata            string `json:"sdata" validate:"required"`
	Page             string `json:"page" validate:"required"`
	Categame_id      string `json:"categame_id"`
	Categame_name    string `json:"categame_name" validate:"required"`
	Categame_display int    `json:"categame_display" validate:"required"`
	Categame_status  string `json:"categame_status" validate:"required"`
}
