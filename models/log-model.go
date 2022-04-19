package models

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/wl_apisuper/configs"
	"github.com/nikitamirzani323/wl_apisuper/db"
	"github.com/nikitamirzani323/wl_apisuper/entities"
	"github.com/nikitamirzani323/wl_apisuper/helpers"
)

func Fetch_logHome(typeuser string, page int) (helpers.Response, error) {
	var obj entities.Model_log
	var arraobj []entities.Model_log
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	totalrecord := 0
	offset := page
	sql_selectcount := ""
	sql_selectcount += ""
	sql_selectcount += "SELECT "
	sql_selectcount += "COUNT(idlog) as totalidlog  "
	sql_selectcount += "FROM " + configs.DB_tbl_trx_log + "  "
	sql_selectcount += "WHERE typeuser=$1  "

	row_selectcount := con.QueryRowContext(ctx, sql_selectcount)
	switch e_selectcount := row_selectcount.Scan(&totalrecord); e_selectcount {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e_selectcount)
	}
	sql_select := ""
	sql_select += ""
	sql_select += "SELECT "
	sql_select += "idlog, to_char(COALESCE(datetimelog,now()), 'YYYY-MM-DD HH24:MI:SS'), userlog, pagelog,  "
	sql_select += "tipelog, notelog    "
	sql_select += "FROM " + configs.DB_tbl_trx_log + "  "
	sql_select += "WHERE typeuser=$1   "
	sql_select += "ORDER BY datetimelog DESC  OFFSET " + strconv.Itoa(offset) + " LIMIT " + strconv.Itoa(configs.PERPAGE)

	row, err := con.QueryContext(ctx, sql_select, typeuser)

	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idlog_db                                                       int
			datetimelog_db, userlog_db, pagelog_db, tipelog_db, notelog_db string
		)

		err = row.Scan(&idlog_db, &datetimelog_db, &userlog_db, &pagelog_db, &tipelog_db, &notelog_db)

		helpers.ErrorCheck(err)

		obj.Log_id = idlog_db
		obj.Log_datetime = datetimelog_db
		obj.Log_user = userlog_db
		obj.Log_page = pagelog_db
		obj.Log_tipe = tipelog_db
		obj.Log_note = notelog_db
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
