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

func Fetch_gameHome() (helpers.Responsegame, error) {
	var obj entities.Model_game
	var arraobj []entities.Model_game
	var res helpers.Responsegame
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idgame , idcategame, idprovidergame, nmgame, imgcovergame, imgthumbgame, endpointgameurl, statusgame,  
			creategame, to_char(COALESCE(createdategame,now()), 'YYYY-MM-DD HH24:MI:SS') as createdategame, 
			updategame, to_char(COALESCE(updatedategame,now()), 'YYYY-MM-DD HH24:MI:SS') as updatedategame 
			FROM ` + configs.DB_tbl_mst_game + ` 
			ORDER BY creategame DESC   
		`

	row, err := con.QueryContext(ctx, sql_select)

	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idgame_db                                                                                                        int
			idcategame_db, idprovidergame_db, nmgame_db, imgcovergame_db, imgthumbgame_db, endpointgameurl_db, statusgame_db string
			creategame_db, createdategame_db, updategame_db, updatedategame_db                                               string
		)

		err = row.Scan(
			&idgame_db, &idcategame_db, &idprovidergame_db, &nmgame_db, &imgcovergame_db, &imgthumbgame_db, &endpointgameurl_db, &statusgame_db,
			&creategame_db, &createdategame_db, &updategame_db, &updatedategame_db)

		helpers.ErrorCheck(err)

		create := creategame_db + ", " + createdategame_db
		update := ""
		if updategame_db != "" {
			update = updategame_db + ", " + updatedategame_db
		}
		obj.Game_id = idgame_db
		obj.Game_idcategame = idcategame_db
		obj.Game_idprovidergame = idprovidergame_db
		obj.Game_name = nmgame_db
		obj.Game_imgcover = imgcovergame_db
		obj.Game_imgthumb = imgthumbgame_db
		obj.Game_endpointurl = endpointgameurl_db
		obj.Game_status = statusgame_db
		obj.Game_create = create
		obj.Game_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	var objCategame entities.Model_gamecate
	var arraobjCategame []entities.Model_gamecate
	sql_listcategame := `SELECT 
		idcategame, nmcategame 	
		FROM ` + configs.DB_tbl_mst_categame + ` 
		WHERE statuscategame = 'Y' 
	`
	row_listcategame, err_listcategame := con.QueryContext(ctx, sql_listcategame)
	helpers.ErrorCheck(err_listcategame)
	for row_listcategame.Next() {
		var (
			idcategame_db, nmcategame_db string
		)

		err = row_listcategame.Scan(&idcategame_db, &nmcategame_db)

		helpers.ErrorCheck(err)

		objCategame.Categame_id = idcategame_db
		objCategame.Categame_name = nmcategame_db
		arraobjCategame = append(arraobjCategame, objCategame)
		msg = "Success"
	}

	var objgameprovider entities.Model_gameprovider
	var arraobjgameprovider []entities.Model_gameprovider
	sql_listgameprovider := `SELECT 
		idprovidergame, nmprovidergame 	
		FROM ` + configs.DB_tbl_mst_providergame + ` 
		WHERE statusprovidergame = 'Y' 
	`
	row_listgameprovider, err_listgameprovider := con.QueryContext(ctx, sql_listgameprovider)
	helpers.ErrorCheck(err_listgameprovider)
	for row_listgameprovider.Next() {
		var (
			idprovidergame_db, nmprovidergame_db string
		)

		err = row_listgameprovider.Scan(&idprovidergame_db, &nmprovidergame_db)

		helpers.ErrorCheck(err)

		objgameprovider.Providergame_id = idprovidergame_db
		objgameprovider.Providergame_name = nmprovidergame_db
		arraobjgameprovider = append(arraobjgameprovider, objgameprovider)
		msg = "Success"
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Listcategame = arraobjCategame
	res.Listprovidergame = arraobjgameprovider
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_gameHome(
	admin, idcategame, idprovidergame, nmgame, imgcovergame, imgthumbgame, endpointgameurl,
	status, sData string, idgame int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	if sData == "New" {
		sql_insert := `
			insert into
			` + configs.DB_tbl_mst_game + ` (
				idgame , idcategame, idprovidergame, nmgame, imgcovergame, imgthumbgame, endpointgameurl, statusgame,
				creategame, createdategame 
			) values (
				$1, $2, $3, $4, $5, $6, $7, $8, 
				$9, $10 
			)
		`
		yearcounter := tglnow.Format("YYYY")
		monthcounter := tglnow.Format("MM")
		idcounter := Get_counter(configs.DB_tbl_mst_game + yearcounter + monthcounter)
		idcounter_final := yearcounter + monthcounter + strconv.Itoa(idcounter)
		flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_game, "INSERT",
			idcounter_final, idcategame, idprovidergame, nmgame, imgcovergame, imgthumbgame, endpointgameurl, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			msg = "Succes"
			log.Println(msg_insert)

			notelog := ""
			notelog += "NEW GAME <br>"
			notelog += "GAME CATEGORY : " + idcategame + "<br>"
			notelog += "GAME PROVIDER : " + idprovidergame + "<br>"
			notelog += "GAME NAME : " + nmgame + "<br>"
			notelog += "IMAGE COVER : " + imgcovergame + "<br>"
			notelog += "IMAGE THUMB : " + imgthumbgame + "<br>"
			notelog += "STATUS : " + status
			Insert_log("SUPERADMIN", "", admin, "GAME", "INSERT", notelog)
		} else {
			log.Println(msg_insert)
		}

	} else {
		sql_update2 := `
				UPDATE 
				` + configs.DB_tbl_mst_game + `   
				SET idcategame=$1, idprovidergame=$2, nmgame=$3, imgcovergame=$4, imgthumbgame=$5, endpointgameurl=$6, statusgame=$7,
				updategame=$8, updatedategame=$9  
				WHERE idgame=$10   
			`
		flag_update, msg_update := Exec_SQL(sql_update2, configs.DB_tbl_mst_game, "UPDATE",
			idcategame, idprovidergame, nmgame, imgcovergame, imgthumbgame, endpointgameurl, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idgame)

		if flag_update {
			msg = "Succes"
			log.Println(msg_update)

			notelog := ""
			notelog += "UPDATE GAME <br>"
			notelog += "GAME CATEGORY : " + idcategame + "<br>"
			notelog += "GAME PROVIDER : " + idprovidergame + "<br>"
			notelog += "GAME NAME : " + nmgame + "<br>"
			notelog += "IMAGE COVER : " + imgcovergame + "<br>"
			notelog += "IMAGE THUMB : " + imgthumbgame + "<br>"
			notelog += "STATUS : " + status
			Insert_log("SUPERADMIN", "", admin, "GAME", "UPDATE", notelog)
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
