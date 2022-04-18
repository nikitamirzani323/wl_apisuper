package models

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/wl_apisuper/configs"
	"github.com/nikitamirzani323/wl_apisuper/db"
	"github.com/nikitamirzani323/wl_apisuper/entities"
	"github.com/nikitamirzani323/wl_apisuper/helpers"
	"github.com/nleeper/goment"
)

func Fetch_categameHome() (helpers.Response, error) {
	var obj entities.Model_categame
	var arraobj []entities.Model_categame
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idcategame , nmcategame, displaycategame, statuscategame, 
			createcategame, to_char(COALESCE(createdatecategame,now()), 'YYYY-MM-DD HH24:MI:SS') as createdatecategame, 
			updatecategame, to_char(COALESCE(updatedatecategame,now()), 'YYYY-MM-DD HH24:MI:SS') as updatedatecategame 
			FROM ` + configs.DB_tbl_mst_categame + ` 
			ORDER BY displaycategame ASC  
		`

	row, err := con.QueryContext(ctx, sql_select)

	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			displaycategame_db                                                                 int
			idcategame_db, nmcategame_db, statuscategame_db                                    string
			createcategame_db, createdatecategame_db, updatecategame_db, updatedatecategame_db string
		)

		err = row.Scan(
			&idcategame_db, &nmcategame_db, &displaycategame_db, &statuscategame_db,
			&createcategame_db, &createdatecategame_db, &updatecategame_db, &updatedatecategame_db)

		helpers.ErrorCheck(err)

		create := createcategame_db + ", " + createdatecategame_db
		update := ""
		if updatecategame_db != "" {
			update = updatecategame_db + ", " + updatedatecategame_db
		}
		obj.Categame_id = idcategame_db
		obj.Categame_name = nmcategame_db
		obj.Categame_display = displaycategame_db
		obj.Categame_status = statuscategame_db
		obj.Categame_create = create
		obj.Categame_update = update
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
func Save_categameHome(
	admin, idcategame, nmcategame,
	status, sData string, display int) (helpers.Response, error) {
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
		flag = CheckDB(configs.DB_tbl_mst_categame, "idcategame", idcategame)
		if !flag {
			sql_insert := `
				insert into
				` + configs.DB_tbl_mst_categame + ` (
					idcategame , nmcategame, displaycategame, statuscategame, 
					createcategame, createdatecategame
				) values (
					$1, $2, $3, $4,   
					$5, $6 
				)
			`

			flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_categame, "INSERT",
				idcategame, nmcategame, display, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

			if flag_insert {
				msg = "Succes"
				log.Println(msg_insert)

				notelog := ""
				notelog += "NEW CATEGORY GAME <br>"
				notelog += "NAME : " + nmcategame + "<br>"
				notelog += "DISPLAY : " + strconv.Itoa(display) + "<br>"
				notelog += "STATUS : " + status
				Insert_log("SUPERADMIN", "", admin, "CATEGORY GAME", "INSERT", notelog)
			} else {
				log.Println(msg_insert)
			}
		} else {
			msg = "Duplicate Entry"
		}

	} else {
		sql_update2 := `
				UPDATE 
				` + configs.DB_tbl_mst_categame + `   
				SET nmcategame=$1, displaycategame=$2, statuscategame=$3, 
				updatecategame=$4, updatedatecategame=$5   
				WHERE idcategame=$6 
			`
		flag_update, msg_update := Exec_SQL(sql_update2, configs.DB_tbl_mst_categame, "UPDATE",
			nmcategame, display, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idcategame)

		if flag_update {
			msg = "Succes"
			log.Println(msg_update)

			notelog := ""
			notelog += "UPDATE CATEGORY GAME <br>"
			notelog += "NAME : " + nmcategame + "<br>"
			notelog += "DISPLAY : " + strconv.Itoa(display) + "<br>"
			notelog += "STATUS : " + status
			Insert_log("SUPERADMIN", "", admin, "CATEGORY GAME", "UPDATE", notelog)
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
