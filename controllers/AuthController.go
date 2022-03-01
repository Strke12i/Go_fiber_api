package controllers

import (
	"freecode/database"
	"freecode/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
	"time"
)

const SecretKey = "secret"

func Register(c *fiber.Ctx) {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		log.Panicln(err)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	if err != nil {
		log.Panicln(err)
	}

	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}

	database.DB.Create(&user)
	c.JSON(user)
}

func Login(c *fiber.Ctx) {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		log.Panicln(err)
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(404)
		c.JSON(fiber.Map{
			"message": "user not found",
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		c.JSON(fiber.Map{
			"message": "Incorrect password!",
		})
		return
	}

	clains := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
	})

	token, err := clains.SignedString([]byte(SecretKey))
	if err != nil {
		c.Status(404)
		c.JSON(fiber.Map{
			"message": "could not login!",
		})
		return
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 2),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	c.JSON(fiber.Map{
		"message": "success",
	})

}

func User(c *fiber.Ctx) {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
		return
	}

	claims := token.Claims.(*jwt.StandardClaims)
	var user models.User

	database.DB.Where("id = ?", claims.Issuer).First(&user)

	c.JSON(user)
}

func Logout(c *fiber.Ctx) {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	c.JSON(fiber.Map{
		"message": "Logout success",
	})
}
