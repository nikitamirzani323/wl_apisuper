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

const Fieldbanktype_home_redis = "LISTBANKTYPE_SUPER_WL"

func Banktypehome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_banktype)
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

	var obj entities.Model_banktype
	var arraobj []entities.Model_banktype
	var objcatebank entities.Model_banktypecatebank
	var arraobjcatebank []entities.Model_banktypecatebank
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldbanktype_home_redis)
	jsonredis := []byte(resultredis)
	listcatebank_RD, _, _, _ := jsonparser.Get(jsonredis, "listcatebank")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		banktype_id, _ := jsonparser.GetString(value, "banktype_id")
		banktype_idcatebank, _ := jsonparser.GetInt(value, "banktype_idcatebank")
		banktype_nmcatebank, _ := jsonparser.GetString(value, "banktype_nmcatebank")
		banktype_name, _ := jsonparser.GetString(value, "banktype_name")
		banktype_img, _ := jsonparser.GetString(value, "banktype_img")
		banktype_status, _ := jsonparser.GetString(value, "banktype_status")
		banktype_create, _ := jsonparser.GetString(value, "banktype_create")
		banktype_update, _ := jsonparser.GetString(value, "banktype_update")

		obj.Banktype_id = banktype_id
		obj.Banktype_idcatebank = int(banktype_idcatebank)
		obj.Banktype_nmcatebank = banktype_nmcatebank
		obj.Banktype_name = banktype_name
		obj.Banktype_img = banktype_img
		obj.Banktype_status = banktype_status
		obj.Banktype_create = banktype_create
		obj.Banktype_update = banktype_update
		arraobj = append(arraobj, obj)
	})
	jsonparser.ArrayEach(listcatebank_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		catebank_id, _ := jsonparser.GetInt(value, "catebank_id")
		catebank_name, _ := jsonparser.GetString(value, "catebank_name")

		objcatebank.Catebank_id = int(catebank_id)
		objcatebank.Catebank_name = catebank_name
		arraobjcatebank = append(arraobjcatebank, objcatebank)
	})
	if !flag {
		result, err := models.Fetch_banktypeHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldbanktype_home_redis, result, 30*time.Hour)
		log.Println("CATEBANK MYSQL")
		return c.JSON(result)
	} else {
		log.Println("CATEBANK CACHE")
		return c.JSON(fiber.Map{
			"status":       fiber.StatusOK,
			"message":      "Success",
			"record":       arraobj,
			"listcatebank": arraobjcatebank,
			"time":         time.Since(render_page).String(),
		})
	}
}
func BanktypeSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_banktypesave)
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

	result, err := models.Save_banktypeHome(
		client_admin,
		client.Banktype_id, client.Banktype_name, client.Banktype_img, client.Banktype_status, client.Sdata, client.Banktype_idcatebank)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_banktype()
	return c.JSON(result)
}
func _deleteredis_banktype() {
	val_super := helpers.DeleteRedis(Fieldbanktype_home_redis)
	log.Printf("REDIS DELETE SUPER BANKTYPE : %d", val_super)

	val_superlog := helpers.DeleteRedis(Fieldlog_home_redis)
	log.Printf("REDIS DELETE SUPER LOG : %d", val_superlog)
}
