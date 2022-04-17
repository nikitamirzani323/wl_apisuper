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

const Fieldgame_home_redis = "LISTGAME_SUPER_WL"

func Gamehome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_game)
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

	var obj entities.Model_game
	var arraobj []entities.Model_game
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldgame_home_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		game_id, _ := jsonparser.GetInt(value, "game_id")
		game_idcategame, _ := jsonparser.GetString(value, "game_idcategame")
		game_idprovidergame, _ := jsonparser.GetString(value, "game_idprovidergame")
		game_name, _ := jsonparser.GetString(value, "game_name")
		game_imgcover, _ := jsonparser.GetString(value, "game_imgcover")
		game_imgthumb, _ := jsonparser.GetString(value, "game_imgthumb")
		game_endpointurl, _ := jsonparser.GetString(value, "game_endpointurl")
		game_status, _ := jsonparser.GetString(value, "game_status")
		game_create, _ := jsonparser.GetString(value, "game_create")
		game_update, _ := jsonparser.GetString(value, "game_update")

		obj.Game_id = int(game_id)
		obj.Game_idcategame = game_idcategame
		obj.Game_idprovidergame = game_idprovidergame
		obj.Game_name = game_name
		obj.Game_imgcover = game_imgcover
		obj.Game_imgthumb = game_imgthumb
		obj.Game_endpointurl = game_endpointurl
		obj.Game_status = game_status
		obj.Game_create = game_create
		obj.Game_update = game_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_gameHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldgame_home_redis, result, 30*time.Hour)
		log.Println("GAME MYSQL")
		return c.JSON(result)
	} else {
		log.Println("GAME CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func GameSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_gamesave)
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

	result, err := models.Save_gameHome(
		client_admin,
		client.Game_idcategame, client.Game_idcategame, client.Game_name, client.Game_imgcover,
		client.Game_imgthumb, client.Game_endpointurl, client.Game_status, client.Sdata, client.Game_id)
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
func _deleteredis_game() {
	val_super := helpers.DeleteRedis(Fieldgame_home_redis)
	log.Printf("REDIS DELETE SUPER GAME : %d", val_super)

	val_superlog := helpers.DeleteRedis(Fieldlog_home_redis)
	log.Printf("REDIS DELETE SUPER LOG : %d", val_superlog)
}
