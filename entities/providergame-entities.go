package entities

type Model_providergame struct {
	Providergame_id     string `json:"providergame_id"`
	Providergame_name   string `json:"providergame_name"`
	Providergame_phone  string `json:"providergame_phone"`
	Providergame_email  string `json:"providergame_email"`
	Providergame_note   string `json:"providergame_note"`
	Providergame_status string `json:"providergame_status"`
	Providergame_create string `json:"providergame_create"`
	Providergame_update string `json:"providergame_update"`
}
type Controller_providergame struct {
	Master string `json:"master" validate:"required"`
}
type Controller_providergamesave struct {
	Sdata               string `json:"sdata" validate:"required"`
	Page                string `json:"page" validate:"required"`
	Providergame_id     string `json:"providergame_id" validate:"required"`
	Providergame_name   string `json:"providergame_name" validate:"required"`
	Providergame_phone  string `json:"providergame_phone" validate:"required"`
	Providergame_email  string `json:"providergame_email" validate:"required"`
	Providergame_note   string `json:"providergame_note"`
	Providergame_status string `json:"providergame_status" validate:"required"`
}
