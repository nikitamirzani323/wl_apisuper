package models

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/wl_apisuper/configs"
	"github.com/nikitamirzani323/wl_apisuper/db"
	"github.com/nikitamirzani323/wl_apisuper/entities"
	"github.com/nikitamirzani323/wl_apisuper/helpers"
	"github.com/nleeper/goment"
)

func Login_Model(username, password string) (helpers.Response, bool, string, error) {
	var obj entities.Model_login
	var arraobj []entities.Model_login
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	render := time.Now()
	flag := false
	flag_login := false
	temp_password := ""
	ruleadmin := ""

	sql_select := `SELECT 
			username , password, idadmin 
			FROM ` + configs.DB_tbl_admin + ` 
			WHERE statuslogin = 'Y'   
		`

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			username_db, password_db, idadmin_db string
		)
		err = row.Scan(&username_db, &password_db, &idadmin_db)
		temp_username := strings.Trim(username, "")
		temp_username = strings.ToLower(temp_username)
		if username_db == temp_username {
			flag_login = true
			temp_password = password_db
			ruleadmin = idadmin_db
		}

		obj.Login_username = username_db
		obj.Login_password = password_db
		obj.Login_idadmin = idadmin_db
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	if flag_login {
		hashpass := helpers.HashPasswordMD5(password)
		log.Println("Pwd : ", password)
		log.Println("Pwd Hash :", hashpass)
		log.Println("DB Hash :", temp_password)
		if hashpass == temp_password {
			flag = true
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Time = time.Since(render).String()

	return res, flag, ruleadmin, nil
}
func Update_login(username, ipaddress, timezone string) {
	tglnow, _ := goment.New()
	sql_update := `
			UPDATE ` + configs.DB_tbl_admin + ` 
			SET lastlogin=$1, ipaddress=$2 , timezone=$3, 
			updateadmin=$4,  updatedateadmin=$5  
			WHERE username  = $6 
			AND statuslogin = 'Y' 
		`
	lastlogin := tglnow.Format("YYYY-MM-DD HH:mm:ss")
	flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_admin, "UPDATE",
		lastlogin, ipaddress, timezone, username,
		tglnow.Format("YYYY-MM-DD HH:mm:ss"), username)

	if flag_update {
		log.Println(msg_update)
	}

	notelog := ""
	notelog += "DATE : " + lastlogin + " <br>"
	notelog += "IP : " + ipaddress
	Insert_log("SUPERADMIN", "", username, "LOGIN", "UPDATE", notelog)
}
