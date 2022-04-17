package helpers

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Record  interface{} `json:"record"`
	Time    string      `json:"time"`
}
type ResponseBanktype struct {
	Status       int         `json:"status"`
	Message      string      `json:"message"`
	Record       interface{} `json:"record"`
	Listcatebank interface{} `json:"listcatebank"`
	Time         string      `json:"time"`
}
type ResponseCompany struct {
	Status   int         `json:"status"`
	Message  string      `json:"message"`
	Record   interface{} `json:"record"`
	Listcurr interface{} `json:"listcurr"`
	Time     string      `json:"time"`
}
type ResponseAdmin struct {
	Status   int         `json:"status"`
	Message  string      `json:"message"`
	Record   interface{} `json:"record"`
	Listrule interface{} `json:"listruleadmin"`
	Time     string      `json:"time"`
}
type ErrorResponse struct {
	Field string
	Tag   string
}

func ErrorCheck(err error) {
	if err != nil {
		panic(err.Error())
	}
}
