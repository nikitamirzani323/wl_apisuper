package models

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/wl_apisuper/configs"
	"github.com/nikitamirzani323/wl_apisuper/db"
	"github.com/nikitamirzani323/wl_apisuper/entities"
	"github.com/nikitamirzani323/wl_apisuper/helpers"
	"github.com/nleeper/goment"
)

func Fetch_banktypeHome() (helpers.ResponseBanktype, error) {
	var obj entities.Model_banktype
	var arraobj []entities.Model_banktype
	var res helpers.ResponseBanktype
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idbanktype, idcatebank , nmbanktype, imgbanktype, statusbanktype, 
			createbanktype, to_char(COALESCE(createdatebanktype,now()), 'YYYY-MM-DD HH24:MI:SS') as createdatebanktype, 
			updatebanktype, to_char(COALESCE(updatedatebanktype,now()), 'YYYY-MM-DD HH24:MI:SS') as updatedatebanktype 
			FROM ` + configs.DB_tbl_mst_banktype + ` 
			ORDER BY createdatebanktype DESC 
		`

	row, err := con.QueryContext(ctx, sql_select)

	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idcatebank_db                                                                      int
			idbanktype_db, nmbanktype_db, imgbanktype_db, statusbanktype_db                    string
			createbanktype_db, createdatebanktype_db, updatebanktype_db, updatedatebanktype_db string
		)

		err = row.Scan(
			&idbanktype_db, &idcatebank_db, &nmbanktype_db, &imgbanktype_db, &statusbanktype_db,
			&createbanktype_db, &createdatebanktype_db, &updatebanktype_db, &updatedatebanktype_db)

		helpers.ErrorCheck(err)

		create := createbanktype_db + ", " + createdatebanktype_db
		update := ""
		if updatebanktype_db != "" {
			update = updatebanktype_db + ", " + updatedatebanktype_db
		}
		obj.Banktype_id = idbanktype_db
		obj.Banktype_idcatebank = idcatebank_db
		obj.Banktype_nmcatebank = _categorybank("nmcatebank", idcatebank_db)
		obj.Banktype_name = nmbanktype_db
		obj.Banktype_img = imgbanktype_db
		obj.Banktype_status = statusbanktype_db
		obj.Banktype_create = create
		obj.Banktype_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	var objCatebank entities.Model_banktypecatebank
	var arraobjCatebank []entities.Model_banktypecatebank
	sql_listcatebank := `SELECT 
		idcatebank, nmcatebank  	
		FROM ` + configs.DB_tbl_mst_catebank + ` 
		WHERE statuscatebank = 'Y' 
	`
	row_listcatebank, err_listcatebank := con.QueryContext(ctx, sql_listcatebank)
	helpers.ErrorCheck(err_listcatebank)
	for row_listcatebank.Next() {
		var (
			idcatebank_db int
			nmcatebank_db string
		)

		err = row_listcatebank.Scan(&idcatebank_db, &nmcatebank_db)

		helpers.ErrorCheck(err)

		objCatebank.Catebank_id = idcatebank_db
		objCatebank.Catebank_name = nmcatebank_db
		arraobjCatebank = append(arraobjCatebank, objCatebank)
		msg = "Success"
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Listcatebank = arraobjCatebank
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_banktypeHome(
	admin, idbanktype, nmbanktype, imgbanktype,
	status, sData string, idcatebank int) (helpers.Response, error) {
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
		sql_insert := `
				insert into
				` + configs.DB_tbl_mst_banktype + ` (
					idbanktype, idcatebank, nmbanktype, imgbanktype, statusbanktype, 
					createbanktype, createdatebanktype
				) values (
					$1, $2, $3, $4, $5,
					$6, $7  
					
				)
			`
		flag = CheckDB(configs.DB_tbl_mst_banktype, "idbanktype", idbanktype)
		if !flag {
		} else {

		}
		flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_banktype, "INSERT",
			idbanktype, idcatebank, nmbanktype, imgbanktype, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			msg = "Succes"
			log.Println(msg_insert)

			notelog := ""
			notelog += "NEW BANK TYPE <br>"
			notelog += "NAME : " + nmbanktype + "<br>"
			notelog += "STATUS : " + status
			Insert_log("SUPERADMIN", "", admin, "BANK TYPE", "INSERT", notelog)
		} else {
			log.Println(msg_insert)
		}
	} else {
		sql_update2 := `
				UPDATE 
				` + configs.DB_tbl_mst_banktype + `   
				SET idcatebank=$1, nmbanktype=$2, imgbanktype=$3, statusbanktype=$4, 
				updatebanktype=$5, updatedatebanktype=$6   
				WHERE idbanktype =$7  
			`
		flag_update, msg_update := Exec_SQL(sql_update2, configs.DB_tbl_mst_banktype, "UPDATE",
			idcatebank, nmbanktype, imgbanktype, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idbanktype)

		if flag_update {
			msg = "Succes"
			log.Println(msg_update)

			notelog := ""
			notelog += "UPDATE BANK TYPE <br>"
			notelog += "NAME : " + nmbanktype + "<br>"
			notelog += "STATUS : " + status
			Insert_log("SUPERADMIN", "", admin, "BANK TYPE", "UPDATE", notelog)
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
func _categorybank(tipe string, idcatebank int) string {
	con := db.CreateCon()
	ctx := context.Background()
	flag := false
	result := ""
	nmcatebank := ""

	sql_select := `SELECT
		nmcatebank   
		FROM ` + configs.DB_tbl_mst_catebank + `  
		WHERE idcatebank = $1 
	`
	row := con.QueryRowContext(ctx, sql_select, idcatebank)
	switch e := row.Scan(&nmcatebank); e {
	case sql.ErrNoRows:
		flag = false
	case nil:
		flag = true

	default:
		panic(e)
	}
	if flag {
		switch tipe {
		case "nmcatebank":
			result = nmcatebank
		}

	}
	return result
}
