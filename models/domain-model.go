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

func Fetch_domainHome() (helpers.Response, error) {
	var obj entities.Model_domain
	var arraobj []entities.Model_domain
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			iddomain , nmdomain, statusdomain,  
			createdomain, COALESCE(createdatedomain,now()), updatedomain, COALESCE(updatedatedomain,now())
			FROM ` + configs.DB_tbl_mst_domain + `  
			ORDER BY iddomain DESC   
	`

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			iddomain_db                                                                int
			nmdomain_db, statusdomain_db                                               string
			createdomain_db, createdatedomain_db, updatedomain_db, updatedatedomain_db string
		)

		err = row.Scan(&iddomain_db, &nmdomain_db, &statusdomain_db,
			&createdomain_db, &createdatedomain_db, &updatedomain_db, &updatedatedomain_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		if createdomain_db != "" {
			create = createdomain_db + ", " + createdatedomain_db
		}
		if updatedomain_db != "" {
			update = updatedomain_db + ", " + updatedatedomain_db
		}

		obj.Domain_id = iddomain_db
		obj.Domain_name = nmdomain_db
		obj.Domain_status = statusdomain_db
		obj.Domain_create = create
		obj.Domain_update = update
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
func Save_domain(admin, nmdomain, status, sData string, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		flag = CheckDB(configs.DB_tbl_mst_domain, "nmdomain", nmdomain)
		if !flag {
			sql_insert := `
				insert into
				` + configs.DB_tbl_mst_domain + ` (
					iddomain , nmdomain, statusdomain, 
					createdomain, createdatedomain
				) values (
					$1, $2, $3, 
					$4, $5
				)
			`
			field_column := configs.DB_tbl_mst_domain + tglnow.Format("YYYY")
			idrecord_counter := Get_counter(field_column)
			flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_domain, "INSERT",
				tglnow.Format("YY")+strconv.Itoa(idrecord_counter), nmdomain, status,
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
		sql_update := `
				UPDATE 
				` + configs.DB_tbl_mst_domain + `  
				SET nmdomain =$1, statusdomain=$2, 
				updatedomain=$3, updatedatedomain=$4 
				WHERE iddomain =$5 
			`

		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_domain, "UPDATE",
			nmdomain, status,
			admin,
			tglnow.Format("YYYY-MM-DD HH:mm:ss"),
			idrecord)

		if flag_update {
			flag = true
			msg = "Succes"
			log.Println(msg_update)
		} else {
			log.Println(msg_update)
		}
	}

	if flag {
		res.Status = fiber.StatusOK
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	} else {
		res.Status = fiber.StatusBadRequest
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	}

	return res, nil
}
