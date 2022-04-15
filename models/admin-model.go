package models

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/wl_api_master/configs"
	"github.com/nikitamirzani323/wl_api_master/db"
	"github.com/nikitamirzani323/wl_api_master/entities"
	"github.com/nikitamirzani323/wl_api_master/helpers"
	"github.com/nleeper/goment"
)

func Fetch_adminHome() (helpers.ResponseAdmin, error) {
	var obj entities.Model_admin
	var arraobj []entities.Model_admin
	var res helpers.ResponseAdmin
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			username , name, idadmin,
			statuslogin, to_char(COALESCE(lastlogin,now()), 'YYYY-MM-DD HH24:MI:SS') as lastlogin, to_char(COALESCE(joindate,now()), 'YYYY-MM-DD') as joindate, 
			ipaddress, timezone,   
			createadmin, to_char(COALESCE(createdateadmin,now()), 'YYYY-MM-DD HH24:MI:SS') as createdateadmin, 
			updateadmin, to_char(COALESCE(updatedateadmin,now()), 'YYYY-MM-DD HH24:MI:SS') as updatedateadmin 
			FROM ` + configs.DB_tbl_admin + ` 
			ORDER BY lastlogin DESC 
		`

	row, err := con.QueryContext(ctx, sql_select)

	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			username_db, name_db, idadminlevel_db                                  string
			statuslogin_db, lastlogin_db, joindate_db, ipaddress_db, timezone_db   string
			createadmin_db, createdateadmin_db, updateadmin_db, updatedateadmin_db string
		)

		err = row.Scan(
			&username_db, &name_db, &idadminlevel_db,
			&statuslogin_db, &lastlogin_db, &joindate_db,
			&ipaddress_db, &timezone_db,
			&createadmin_db, &createdateadmin_db, &updateadmin_db, &updatedateadmin_db)

		helpers.ErrorCheck(err)
		if statuslogin_db == "Y" {
			statuslogin_db = "ACTIVE"
		}
		if lastlogin_db == "0000-00-00 00:00:00" {
			lastlogin_db = ""
		}
		create := createadmin_db + ", " + createdateadmin_db
		update := ""
		if updateadmin_db != "" {
			update = updateadmin_db + ", " + updatedateadmin_db
		}
		obj.Admin_username = username_db
		obj.Admin_nama = name_db
		obj.Admin_rule = idadminlevel_db
		obj.Admin_joindate = joindate_db
		obj.Admin_timezone = timezone_db
		obj.Admin_lastlogin = lastlogin_db
		obj.Admin_lastIpaddress = ipaddress_db
		obj.Admin_status = statuslogin_db
		obj.Admin_create = create
		obj.Admin_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	var objRule entities.Model_adminrule
	var arraobjRule []entities.Model_adminrule
	sql_listrule := `SELECT 
		idadmin 	
		FROM ` + configs.DB_tbl_admingroup + ` 
	`
	row_listrule, err_listrule := con.QueryContext(ctx, sql_listrule)

	helpers.ErrorCheck(err_listrule)
	for row_listrule.Next() {
		var (
			idruleadmin_db string
		)

		err = row_listrule.Scan(&idruleadmin_db)

		helpers.ErrorCheck(err)

		objRule.Admin_idrule = idruleadmin_db
		arraobjRule = append(arraobjRule, objRule)
		msg = "Success"
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Listrule = arraobjRule
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_adminHome(admin, username, password, nama, rule, status, sData string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false
	if status == "ACTIVE" {
		status = "Y"
	} else {
		status = "N"
	}
	if sData == "New" {
		flag = CheckDB(configs.DB_tbl_admin, "username", username)
		if !flag {
			sql_insert := `
				insert into
				` + configs.DB_tbl_admin + ` (
					username , password, idadmin, name, statuslogin, joindate, 
					createadmin, createdateadmin
				) values (
					$1, $2, $3, $4, $5, $6, 
					$7, $8
				)
			`
			hashpass := helpers.HashPasswordMD5(password)
			flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_admin, "INSERT",
				username, hashpass,
				rule, nama, status,
				tglnow.Format("YYYY-MM-DD HH:mm:ss"),
				admin,
				tglnow.Format("YYYY-MM-DD HH:mm:ss"))

			if flag_insert {
				flag = true
				msg = "Succes"
				log.Println(msg_insert)
			} else {
				log.Println(msg_insert)
			}
		} else {
			msg = "Duplicate Entry"
		}
	} else {
		if password == "" {
			sql_update := `
				UPDATE 
				` + configs.DB_tbl_admin + `  
				SET name =$1, idadmin=$2, statuslogin=$3,  
				updateadmin=$4, updatedateadmin=$5 
				WHERE username =$6 
			`

			flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_admin, "UPDATE",
				nama, rule, status, admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"),
				username)

			if flag_update {
				flag = true
				msg = "Succes"
				log.Println(msg_update)
			} else {
				log.Println(msg_update)
			}
		} else {
			hashpass := helpers.HashPasswordMD5(password)
			sql_update2 := `
				UPDATE 
				` + configs.DB_tbl_admin + `   
				SET name =$1, password=$2, idadmin=$3, statuslogin=$4,  
				updateadmin=$5, updatedateadmin=$6 
				WHERE username =$7 
			`
			flag_update, msg_update := Exec_SQL(sql_update2, configs.DB_tbl_admin, "UPDATE",
				nama,
				hashpass,
				rule,
				status,
				admin,
				tglnow.Format("YYYY-MM-DD HH:mm:ss"),
				username)

			if flag_update {
				flag = true
				msg = "Succes"
				log.Println(msg_update)
			} else {
				log.Println(msg_update)
			}
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
