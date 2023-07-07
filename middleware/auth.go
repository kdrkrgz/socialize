package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/kdrkrgz/socalize/conf"
	"github.com/kdrkrgz/socalize/repository"
	resp "github.com/kdrkrgz/socalize/utils"
)

func DeserializeUser(c *fiber.Ctx) error {
	tokenString := ""
	if authorization := c.Get("Authorization"); strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if token := c.Cookies("token"); token != "" {
		tokenString = token
	}

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(resp.ErrorResponse(resp.ApiErrors("UnAuthorized")))
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(conf.Get("TokenSecret")), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(resp.ErrorResponse(resp.ApiErrors("InvalidToken")))
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "invalid token claim"})
	}

	id, ok := claims["sub"].(map[string]interface{})["id"].(float64)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(resp.ErrorResponse(resp.ApiErrors("UserNotExist")))
	}

	user := repository.New().GetUserById(uint(id))
	if user.Id != uint(id) {
		return c.Status(fiber.StatusForbidden).JSON(resp.ErrorResponse(resp.ApiErrors("UserNotExist")))
	}

	return c.Next()
}
