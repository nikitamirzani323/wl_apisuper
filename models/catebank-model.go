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

func Fetch_catebankHome() (helpers.Response, error) {
	var obj entities.Model_catebank
	var arraobj []entities.Model_catebank
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idcatebank , nmcatebank, statuscatebank, 
			createcatebank, to_char(COALESCE(createdatecatebank,now()), 'YYYY-MM-DD HH24:MI:SS') as createdatecatebank, 
			updatecatebank, to_char(COALESCE(updatedatecatebank,now()), 'YYYY-MM-DD HH24:MI:SS') as updatedatecatebank 
			FROM ` + configs.DB_tbl_mst_catebank + ` 
			ORDER BY createdatecatebank DESC 
		`

	row, err := con.QueryContext(ctx, sql_select)

	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idcatebank_db                                                                      int
			nmcatebank_db, statuscatebank_db                                                   string
			createcatebank_db, createdatecatebank_db, updatecatebank_db, updatedatecatebank_db string
		)

		err = row.Scan(
			&idcatebank_db, &nmcatebank_db, &statuscatebank_db,
			&createcatebank_db, &createdatecatebank_db, &updatecatebank_db, &updatedatecatebank_db)

		helpers.ErrorCheck(err)

		create := createcatebank_db + ", " + createdatecatebank_db
		update := ""
		if updatecatebank_db != "" {
			update = updatecatebank_db + ", " + updatedatecatebank_db
		}
		obj.Catebank_id = idcatebank_db
		obj.Catebank_name = nmcatebank_db
		obj.Catebank_status = statuscatebank_db
		obj.Catebank_create = create
		obj.Catebank_update = update
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
func Save_catebankHome(
	admin, nmcatebank,
	status, sData string, idcatebank int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	if sData == "New" {
		sql_insert := `
				insert into
				` + configs.DB_tbl_mst_catebank + ` (
					idcatebank , nmcatebank, statuscatebank, 
					createcatebank, createdatecatebank
				) values (
					$1, $2, $3,   
					$4, $5 
				)
			`
		yearcounter := tglnow.Format("YYYY")
		monthcounter := tglnow.Format("MM")
		idcounter := Get_counter(configs.DB_tbl_mst_catebank + yearcounter + monthcounter)
		idcounter_final := yearcounter + monthcounter + strconv.Itoa(idcounter)
		flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_catebank, "INSERT",
			idcounter_final, nmcatebank, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			msg = "Succes"
			log.Println(msg_insert)

			notelog := ""
			notelog += "NEW CATEGORY BANK <br>"
			notelog += "NAME : " + nmcatebank + "<br>"
			notelog += "STATUS : " + status
			Insert_log("SUPERADMIN", "", admin, "CATEGORY BANK", "INSERT", notelog)
		} else {
			log.Println(msg_insert)
		}
	} else {
		sql_update2 := `
				UPDATE 
				` + configs.DB_tbl_mst_catebank + `   
				SET nmcatebank=$1, statuscatebank=$2, 
				updatecatebank=$3, updatedatecatebank=$4  
				WHERE idcatebank =$5 
			`
		flag_update, msg_update := Exec_SQL(sql_update2, configs.DB_tbl_mst_catebank, "UPDATE",
			nmcatebank, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idcatebank)

		if flag_update {
			msg = "Succes"
			log.Println(msg_update)

			notelog := ""
			notelog += "UPDATE CATEGORY BANK <br>"
			notelog += "NAME : " + nmcatebank + "<br>"
			notelog += "STATUS : " + status
			Insert_log("SUPERADMIN", "", admin, "CATEGORY BANK", "UPDATE", notelog)
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
