package models

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/wl_apisuper/configs"
	"github.com/nikitamirzani323/wl_apisuper/db"
	"github.com/nikitamirzani323/wl_apisuper/entities"
	"github.com/nikitamirzani323/wl_apisuper/helpers"
	"github.com/nleeper/goment"
)

func Fetch_companyHome() (helpers.ResponseCompany, error) {
	var obj entities.Model_company
	var arraobj []entities.Model_company
	var res helpers.ResponseCompany
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idcompany , to_char(COALESCE(startjoincompany,now()), 'YYYY-MM-DD HH24:MI:SS'), to_char(COALESCE(endjoincompany,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			idcurr , nmcompany, nmowner, phoneowner, emailowner, companyurl, statuscompany, 
			createcompany, to_char(COALESCE(createdatecompany,now()), 'YYYY-MM-DD HH24:MI:SS') as createdatecompany, 
			updatecompany, to_char(COALESCE(updatedatecompany,now()), 'YYYY-MM-DD HH24:MI:SS') as updatedatecompany 
			FROM ` + configs.DB_tbl_mst_company + ` 
			ORDER BY createcompany DESC 
		`

	row, err := con.QueryContext(ctx, sql_select)

	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idcompany_db, startjoincompany_db, endjoincompany_db                                               string
			idcurr_db, nmcompany_db, nmowner_db, phoneowner_db, emailowner_db, companyurl_db, statuscompany_db string
			createcompany_db, createdatecompany_db, updatecompany_db, updatedatecompany_db                     string
		)

		err = row.Scan(
			&idcompany_db, &startjoincompany_db, &endjoincompany_db,
			&idcurr_db, &nmcompany_db, &nmowner_db, &phoneowner_db, &emailowner_db, &companyurl_db, &statuscompany_db,
			&createcompany_db, &createdatecompany_db, &updatecompany_db, &updatedatecompany_db)

		helpers.ErrorCheck(err)
		if startjoincompany_db == "0000-00-00 00:00:00" {
			startjoincompany_db = ""
		}
		if endjoincompany_db == "0000-00-00 00:00:00" {
			endjoincompany_db = ""
		}
		create := createcompany_db + ", " + createdatecompany_db
		update := ""
		if updatecompany_db != "" {
			update = updatecompany_db + ", " + updatedatecompany_db
		}
		obj.Company_idcomp = idcompany_db
		obj.Company_startjoin = startjoincompany_db
		obj.Company_endjoin = endjoincompany_db
		obj.Company_idcurr = idcurr_db
		obj.Company_nmcompany = nmcompany_db
		obj.Company_nmowner = nmowner_db
		obj.Company_phoneowner = phoneowner_db
		obj.Company_emailowner = emailowner_db
		obj.Company_urlendpoint = companyurl_db
		obj.Company_status = statuscompany_db
		obj.Company_create = create
		obj.Company_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	var objCurr entities.Model_compcurr
	var arraobjCurr []entities.Model_compcurr
	sql_listcurr := `SELECT 
		idcurr 	
		FROM ` + configs.DB_tbl_mst_currency + ` 
	`
	row_listcurr, err_listcurr := con.QueryContext(ctx, sql_listcurr)
	helpers.ErrorCheck(err_listcurr)
	for row_listcurr.Next() {
		var (
			idcurr_db string
		)

		err = row_listcurr.Scan(&idcurr_db)

		helpers.ErrorCheck(err)

		objCurr.Curr_idcurr = idcurr_db
		arraobjCurr = append(arraobjCurr, objCurr)
		msg = "Success"
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Listcurr = arraobjCurr
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_companyHome(
	admin, idcompany, idcurr, nmcompany,
	nmowner, phoneowner, emailowner, companyurl,
	status, sData string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		flag = CheckDB(configs.DB_tbl_mst_company, "idcompany", idcompany)
		if !flag {
			sql_insert := `
				insert into
				` + configs.DB_tbl_mst_company + ` (
					idcompany , startjoincompany, idcurr, nmcompany, nmowner, 
					phoneowner, emailowner, companyurl, statuscompany
					createcompany, createdatecompany
				) values (
					$1, $2, $3, $4, $5,  
					$6, $7, $8, $9, 
					$10, $11 
				)
			`
			flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_company, "INSERT",
				idcompany, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idcurr, nmcompany, nmowner, phoneowner, emailowner, companyurl, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

			if flag_insert {
				flag = true
				msg = "Succes"
				log.Println(msg_insert)

				notelog := ""
				notelog += "NEW COMPANY <br>"
				notelog += "IDCOMPANY : " + idcompany + "<br>"
				notelog += "CURRENCY : " + idcurr + "<br>"
				notelog += "COMPANY : " + nmcompany + "<br>"
				notelog += "OWNER : " + nmowner + "<br>"
				notelog += "PHONE : " + phoneowner + "<br>"
				notelog += "EMAIL : " + emailowner + "<br>"
				notelog += "STATUS : " + status
				Insert_log("SUPERADMIN", "", admin, "COMPANY", "INSERT", notelog)
			} else {
				log.Println(msg_insert)
			}
		} else {
			msg = "Duplicate Entry"
		}
	} else {
		sql_update2 := `
				UPDATE 
				` + configs.DB_tbl_mst_company + `   
				SET nmcompany=$1, nmowner=$2, phoneowner=$3, emailowner=$4,  
				companyurl =$5, statuscompany=$6,  
				updatecompany=$7, updatedatecompany=$8 
				WHERE idcompany =$9 
			`
		flag_update, msg_update := Exec_SQL(sql_update2, configs.DB_tbl_mst_company, "UPDATE",
			nmcompany, nmowner, phoneowner, emailowner, companyurl, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idcompany)

		if flag_update {
			flag = true
			msg = "Succes"
			log.Println(msg_update)

			notelog := ""
			notelog += "UPDATE COMPANY <br>"
			notelog += "COMPANY : " + nmcompany + "<br>"
			notelog += "OWNER : " + nmowner + "<br>"
			notelog += "PHONE : " + phoneowner + "<br>"
			notelog += "EMAIL : " + emailowner + "<br>"
			notelog += "STATUS : " + status
			Insert_log("SUPERADMIN", "", admin, "COMPANY", "UPDATE", notelog)
		} else {
			log.Println(msg_update)
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
func Fetch_companyListAdmin(idcompany string) (helpers.Response, error) {
	var obj entities.Model_companyadmin
	var arraobj []entities.Model_companyadmin
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			username_comp , typeadmin, nama_comp, email_comp, phone_comp, 
			status_comp , to_char(COALESCE(lastlogin_comp,now()), 'YYYY-MM-DD HH24:MI:SS'), lastipaddres_comp, 
			createcomp_admin, to_char(COALESCE(createdatecomp_admin,now()), 'YYYY-MM-DD HH24:MI:SS') as createdatecomp_admin, 
			updatecomp_admin, to_char(COALESCE(updatedatecomp_admin,now()), 'YYYY-MM-DD HH24:MI:SS') as updatedatecomp_admin 
			FROM ` + configs.DB_tbl_mst_company_admin + ` 
			WHERE idcompany = $1 
			ORDER BY lastlogin_comp DESC 
		`

	row, err := con.QueryContext(ctx, sql_select, idcompany)

	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			username_comp_db, typeadmin_db, nama_comp_db, email_comp_db, phone_comp_db                 string
			status_comp_db, lastlogin_comp_db, lastipaddres_comp_db                                    string
			createcomp_admin_db, createdatecomp_admin_db, updatecomp_admin_db, updatedatecomp_admin_db string
		)

		err = row.Scan(
			&username_comp_db, &typeadmin_db, &nama_comp_db, &email_comp_db, &phone_comp_db,
			&status_comp_db, &lastlogin_comp_db, &lastipaddres_comp_db,
			&createcomp_admin_db, &createdatecomp_admin_db, &updatecomp_admin_db, &updatedatecomp_admin_db)

		helpers.ErrorCheck(err)
		if lastlogin_comp_db == "0000-00-00 00:00:00" {
			lastlogin_comp_db = ""
		}
		create := createcomp_admin_db + ", " + createdatecomp_admin_db
		update := ""
		if updatecomp_admin_db != "" {
			update = updatecomp_admin_db + ", " + updatedatecomp_admin_db
		}
		obj.Companyadmin_username = username_comp_db
		obj.Companyadmin_type = typeadmin_db
		obj.Companyadmin_name = nama_comp_db
		obj.Companyadmin_email = email_comp_db
		obj.Companyadmin_phone = phone_comp_db
		obj.Companyadmin_status = status_comp_db
		obj.Companyadmin_lastlogin = lastlogin_comp_db
		obj.Companyadmin_lastipaddress = lastipaddres_comp_db
		obj.Companyadmin_create = create
		obj.Companyadmin_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_companylistadmin(
	admin, idcompany, username, password,
	name, email, phone, status,
	sData string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		flag = CheckDB(configs.DB_tbl_mst_company_admin, "username_comp", username)
		if !flag {
			sql_insert := `
				insert into
				` + configs.DB_tbl_mst_company_admin + ` (
					username_comp , password_comp, idcompany, typeadmin, idruleadmin,  
					nama_comp, email_comp, phone_comp, status_comp,   
					createcomp_admin, createdatecomp_admin 
				) values (
					$1, $2, $3, $4, $5,  
					$6, $7, $8, $9, 
					$10, $11  
				)
			`
			flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_company_admin, "INSERT",
				username, password, idcompany, "MASTER", "0",
				name, email, phone, status, admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

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
		if password != "" {
			haspwd := helpers.HashPasswordMD5(password)
			sql_update2 := `
				UPDATE 
				` + configs.DB_tbl_mst_company_admin + `   
				SET password_comp=$1, nama_comp=$2, email_comp=$3, phone_comp=$4,  
				status_comp=$5, updatecomp_admin=$6, updatedatecomp_admin=$7 
				WHERE idcompany =$8 AND username_comp=$9   
			`
			flag_update, msg_update := Exec_SQL(sql_update2, configs.DB_tbl_mst_company_admin, "UPDATE",
				haspwd, name, email, phone, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idcompany, username)

			if flag_update {
				flag = true
				msg = "Succes"
				log.Println(msg_update)
			} else {
				log.Println(msg_update)
			}
		} else {
			sql_update2 := `
				UPDATE 
				` + configs.DB_tbl_mst_company_admin + `   
				SET nama_comp=$1, email_comp=$2, phone_comp=$3,  
				status_comp=$4, updatecomp_admin=$5, updatedatecomp_admin=$6 
				WHERE idcompany=$7 AND username_comp=$8  
			`
			flag_update, msg_update := Exec_SQL(sql_update2, configs.DB_tbl_mst_company_admin, "UPDATE",
				name, email, phone, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idcompany, username)

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
