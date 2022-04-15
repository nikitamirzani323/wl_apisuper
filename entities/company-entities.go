package entities

type Model_company struct {
	Company_idcomp      string `json:"company_idcompany"`
	Company_startjoin   string `json:"company_startjoin"`
	Company_endjoin     string `json:"company_endjoin"`
	Company_idcurr      string `json:"company_idcurr"`
	Company_nmcompany   string `json:"company_nmcompany"`
	Company_nmowner     string `json:"company_nmowner"`
	Company_phoneowner  string `json:"company_phoneowner"`
	Company_emailowner  string `json:"company_emailowner"`
	Company_urlendpoint string `json:"company_urlendpoint"`
	Company_status      string `json:"company_status"`
	Company_create      string `json:"company_create"`
	Company_update      string `json:"company_update"`
}
type Model_companyadmin struct {
	Companyadmin_username      string `json:"companyadmin_username"`
	Companyadmin_type          string `json:"companyadmin_type"`
	Companyadmin_name          string `json:"companyadmin_name"`
	Companyadmin_email         string `json:"companyadmin_email"`
	Companyadmin_phone         string `json:"companyadmin_phone"`
	Companyadmin_status        string `json:"companyadmin_status"`
	Companyadmin_lastlogin     string `json:"companyadmin_lastlogin"`
	Companyadmin_lastipaddress string `json:"companyadmin_lastipaddress"`
	Companyadmin_create        string `json:"companyadmin_create"`
	Companyadmin_update        string `json:"companyadmin_update"`
}
type Model_compcurr struct {
	Curr_idcurr string `json:"curr_idcurr"`
}
type Controller_company struct {
	Master string `json:"master" validate:"required"`
}
type Controller_companyadmin struct {
	Idcompany string `json:"idcompany" validate:"required"`
}
type Controller_companysave struct {
	Sdata               string `json:"sdata" validate:"required"`
	Page                string `json:"page" validate:"required"`
	Company_idcomp      string `json:"company_idcompany" validate:"required"`
	Company_idcurr      string `json:"company_idcurr" validate:"required"`
	Company_nmcompany   string `json:"company_nmcompany" validate:"required"`
	Company_nmowner     string `json:"company_nmowner" validate:"required"`
	Company_phoneowner  string `json:"company_phoneowner" validate:"required"`
	Company_emailowner  string `json:"company_emailowner" `
	Company_urlendpoint string `json:"company_urlendpoint" validate:"required"`
	Company_status      string `json:"company_status" validate:"required"`
}

type Controller_companysavelistadmin struct {
	Sdata                  string `json:"sdata" validate:"required"`
	Page                   string `json:"page" validate:"required"`
	Companyadmin_idcompany string `json:"companyadmin_idcompany" validate:"required"`
	Companyadmin_username  string `json:"companyadmin_username" validate:"required"`
	Companyadmin_password  string `json:"companyadmin_password" `
	Companyadmin_name      string `json:"companyadmin_name" validate:"required"`
	Companyadmin_email     string `json:"companyadmin_email" `
	Companyadmin_phone     string `json:"companyadmin_phone" validate:"required"`
	Companyadmin_status    string `json:"companyadmin_status" validate:"required"`
}
