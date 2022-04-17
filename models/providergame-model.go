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

func Fetch_providergameHome() (helpers.Response, error) {
	var obj entities.Model_providergame
	var arraobj []entities.Model_providergame
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idprovidergame , nmprovidergame, phoneprovidergame, emailprovidergame, noteprovidergame, statusprovidergame, 
			createprovidergame, to_char(COALESCE(createdateprovidergame,now()), 'YYYY-MM-DD HH24:MI:SS') as createdateprovidergame, 
			updateprovidergame, to_char(COALESCE(updatedateprovidergame,now()), 'YYYY-MM-DD HH24:MI:SS') as updatedateprovidergame 
			FROM ` + configs.DB_tbl_mst_providergame + ` 
			ORDER BY createprovidergame DESC   
		`

	row, err := con.QueryContext(ctx, sql_select)

	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idprovidergame_db, nmprovidergame_db, phoneprovidergame_db, emailprovidergame_db, noteprovidergame_db, statusprovidergame_db string
			createprovidergame_db, createdateprovidergame_db, updateprovidergame_db, updatedateprovidergame_db                           string
		)

		err = row.Scan(
			&idprovidergame_db, &nmprovidergame_db, &phoneprovidergame_db, &emailprovidergame_db, &noteprovidergame_db, &statusprovidergame_db,
			&createprovidergame_db, &createdateprovidergame_db, &updateprovidergame_db, &updatedateprovidergame_db)

		helpers.ErrorCheck(err)

		create := createprovidergame_db + ", " + createdateprovidergame_db
		update := ""
		if updateprovidergame_db != "" {
			update = updateprovidergame_db + ", " + updatedateprovidergame_db
		}
		obj.Providergame_id = idprovidergame_db
		obj.Providergame_name = nmprovidergame_db
		obj.Providergame_phone = phoneprovidergame_db
		obj.Providergame_email = emailprovidergame_db
		obj.Providergame_note = noteprovidergame_db
		obj.Providergame_status = statusprovidergame_db
		obj.Providergame_status = create
		obj.Providergame_update = update
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
func Save_providergameHome(
	admin, idprovidergame, nmprovidergame, phoneprovidergame, emailprovidergame, noteprovidergame,
	status, sData string) (helpers.Response, error) {
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
		flag = CheckDB(configs.DB_tbl_mst_providergame, "idprovidergame", idprovidergame)
		if !flag {
			sql_insert := `
				insert into
				` + configs.DB_tbl_mst_providergame + ` (
					idprovidergame , nmprovidergame, phoneprovidergame, emailprovidergame, noteprovidergame, statusprovidergame,  
					createprovidergame, createdateprovidergame
				) values (
					$1, $2, $3, $4, $5, $6,  
					$7, $8 
				)
			`

			flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_providergame, "INSERT",
				idprovidergame, nmprovidergame, phoneprovidergame, emailprovidergame, noteprovidergame, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

			if flag_insert {
				msg = "Succes"
				log.Println(msg_insert)

				notelog := ""
				notelog += "NEW PROVIDER GAME <br>"
				notelog += "NAME : " + nmprovidergame + "<br>"
				notelog += "PHONE : " + phoneprovidergame + "<br>"
				notelog += "EMAIL : " + emailprovidergame + "<br>"
				notelog += "NOTE : " + noteprovidergame + "<br>"
				notelog += "STATUS : " + status
				Insert_log("SUPERADMIN", "", admin, "PROVIDER GAME", "INSERT", notelog)
			} else {
				log.Println(msg_insert)
			}
		} else {
			msg = "Duplicate Entry"
		}

	} else {
		sql_update2 := `
				UPDATE 
				` + configs.DB_tbl_mst_providergame + `   
				SET nmprovidergame=$1, phoneprovidergame=$2, emailprovidergame=$3, noteprovidergame=$4, statusprovidergame=$5,
				updateprovidergame=$6, updatedateprovidergame=$7  
				WHERE idprovidergame =$8  
			`
		flag_update, msg_update := Exec_SQL(sql_update2, configs.DB_tbl_mst_providergame, "UPDATE",
			nmprovidergame, phoneprovidergame, emailprovidergame, noteprovidergame, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idprovidergame)

		if flag_update {
			msg = "Succes"
			log.Println(msg_update)

			notelog := ""
			notelog += "UPDATE PROVIDER GAME <br>"
			notelog += "NAME : " + nmprovidergame + "<br>"
			notelog += "PHONE : " + phoneprovidergame + "<br>"
			notelog += "EMAIL : " + emailprovidergame + "<br>"
			notelog += "NOTE : " + noteprovidergame + "<br>"
			notelog += "STATUS : " + status
			Insert_log("SUPERADMIN", "", admin, "PROVIDER GAME", "UPDATE", notelog)
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
