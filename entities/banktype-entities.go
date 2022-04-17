package entities

type Model_banktype struct {
	Banktype_id         int    `json:"banktype_id"`
	Banktype_idcatebank int    `json:"banktype_idcatebank"`
	Banktype_name       string `json:"banktype_name"`
	Banktype_img        string `json:"banktype_img"`
	Banktype_status     string `json:"banktype_status"`
	Banktype_create     string `json:"banktype_create"`
	Banktype_update     string `json:"banktype_update"`
}
type Controller_banktype struct {
	Master string `json:"master" validate:"required"`
}
type Controller_banktypesave struct {
	Sdata               string `json:"sdata" validate:"required"`
	Page                string `json:"page" validate:"required"`
	Banktype_id         int    `json:"banktype_id" validate:"required"`
	Banktype_idcatebank int    `json:"banktype_idcatebank" validate:"required"`
	Banktype_name       string `json:"banktype_name" validate:"required"`
	Banktype_img        string `json:"banktype_img"`
	Banktype_status     string `json:"banktype_status" validate:"required"`
}
