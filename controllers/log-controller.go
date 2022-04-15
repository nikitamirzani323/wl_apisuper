package controllers

import (
	"log"
	"time"

	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/wl_apisuper/entities"
	"github.com/nikitamirzani323/wl_apisuper/helpers"
	"github.com/nikitamirzani323/wl_apisuper/models"
)

const Fieldlog_home_redis = "LISTLOG_SUPER_WL"

func Loghome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_log)
	validate := validator.New()
	if err := c.BodyParser(client); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	err := validate.Struct(client)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element helpers.ErrorResponse
			element.Field = err.StructField()
			element.Tag = err.Tag()
			errors = append(errors, &element)
		}
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "validation",
			"record":  errors,
		})
	}

	var obj entities.Model_log
	var arraobj []entities.Model_log
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldlog_home_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		log_id, _ := jsonparser.GetInt(value, "log_id")
		log_datetime, _ := jsonparser.GetString(value, "log_datetime")
		log_user, _ := jsonparser.GetString(value, "log_user")
		log_page, _ := jsonparser.GetString(value, "log_page")
		log_tipe, _ := jsonparser.GetString(value, "log_tipe")
		log_note, _ := jsonparser.GetString(value, "log_note")

		obj.Log_id = int(log_id)
		obj.Log_datetime = log_datetime
		obj.Log_user = log_user
		obj.Log_page = log_page
		obj.Log_tipe = log_tipe
		obj.Log_note = log_note
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_logHome(client.Typeuser)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldlog_home_redis, result, 30*time.Minute)
		log.Println("LOG MYSQL")
		return c.JSON(result)
	} else {
		log.Println("LOG CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
