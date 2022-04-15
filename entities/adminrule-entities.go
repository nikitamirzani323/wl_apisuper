package entities

type Model_adminruleall struct {
	Adminrule_idadmin string `json:"adminrule_idadmin"`
	Adminrule_rule    string `json:"adminrule_rule"`
}

type Controller_adminrule struct {
	Master string `json:"master" validate:"required"`
}

type Controller_adminrulesave struct {
	Sdata   string `json:"sdata" validate:"required"`
	Page    string `json:"page" validate:"required"`
	Idadmin string `json:"adminrule_idadmin" validate:"required"`
	Rule    string `json:"adminrule_rule" `
}

type Responseredis_adminruleall struct {
	Adminrule_idadmin string `json:"adminrule_idadmin"`
	Adminrule_rule    string `json:"adminrule_rule"`
}
