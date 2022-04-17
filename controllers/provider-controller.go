package controllers

import (
	"log"
	"time"

	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nikitamirzani323/wl_apisuper/entities"
	"github.com/nikitamirzani323/wl_apisuper/helpers"
	"github.com/nikitamirzani323/wl_apisuper/models"
)

const Fieldprovidergame_home_redis = "LISTPROVIDERGAME_SUPER_WL"

func Providergamehome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_providergame)
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

	var obj entities.Model_providergame
	var arraobj []entities.Model_providergame
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldprovidergame_home_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		providergame_id, _ := jsonparser.GetString(value, "providergame_id")
		providergame_name, _ := jsonparser.GetString(value, "providergame_name")
		providergame_phone, _ := jsonparser.GetString(value, "providergame_phone")
		providergame_email, _ := jsonparser.GetString(value, "providergame_email")
		providergame_note, _ := jsonparser.GetString(value, "providergame_note")
		providergame_status, _ := jsonparser.GetString(value, "providergame_status")
		providergame_create, _ := jsonparser.GetString(value, "providergame_create")
		providergame_update, _ := jsonparser.GetString(value, "providergame_update")

		obj.Providergame_id = providergame_id
		obj.Providergame_name = providergame_name
		obj.Providergame_phone = providergame_phone
		obj.Providergame_email = providergame_email
		obj.Providergame_note = providergame_note
		obj.Providergame_status = providergame_status
		obj.Providergame_create = providergame_create
		obj.Providergame_update = providergame_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_providergameHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldprovidergame_home_redis, result, 30*time.Hour)
		log.Println("PROVIDERGAME MYSQL")
		return c.JSON(result)
	} else {
		log.Println("PROVIDERGAME CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func ProvidergameSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_providergamesave)
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
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	client_admin, _ := helpers.Parsing_Decry(temp_decp, "==")

	result, err := models.Save_providergameHome(
		client_admin,
		client.Providergame_id, client.Providergame_name, client.Providergame_phone, client.Providergame_email,
		client.Providergame_note, client.Providergame_status, client.Sdata)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_providergame()
	return c.JSON(result)
}
func _deleteredis_providergame() {
	val_super := helpers.DeleteRedis(Fieldprovidergame_home_redis)
	log.Printf("REDIS DELETE SUPER PROVIDER GAME : %d", val_super)

	val_superlog := helpers.DeleteRedis(Fieldlog_home_redis)
	log.Printf("REDIS DELETE SUPER LOG : %d", val_superlog)
}
