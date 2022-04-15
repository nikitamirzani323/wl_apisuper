package entities

type Model_login struct {
	Login_username string `json:"login_username"`
	Login_password string `json:"login_password"`
	Login_idadmin  string `json:"login_idadmin"`
}
type Login struct {
	Client_hostname string `json:"client_hostname" validate:"required"`
	Username        string `json:"username" validate:"required,min=4,max=30,lowercase"`
	Password        string `json:"password" validate:"required"`
	Ipaddress       string `json:"ipaddress" validate:"required"`
	Timezone        string `json:"timezone" validate:"required"`
}
type Home struct {
	Page string `json:"page" validate:"required"`
}
