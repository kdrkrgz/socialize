package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/kdrkrgz/socalize/conf"
	"github.com/kdrkrgz/socalize/repository"
	resp "github.com/kdrkrgz/socalize/utils"
	"strconv"
	"strings"
)

func DeserializeUser(c *fiber.Ctx) error {
	var tokenString string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("token") != "" {
		tokenString = c.Cookies("token")
	}

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(
			resp.ErrorResponse(resp.ApiErrors("UnAuthorized")))
	}

	tokenByte, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
		}

		return []byte(conf.Get("TokenSecret")), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			resp.ErrorResponse(resp.ApiErrors("InvalidToken")))
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "invalid token claim"})

	}

	repo := repository.New()
	claimsSub := claims["sub"].(map[string]interface{})
	id, err := strconv.Atoi(fmt.Sprintf("%v", claimsSub["id"]))
	user := repo.GetUserById(uint(id))
	if float64(user.Id) != claimsSub["id"] {
		return c.Status(fiber.StatusForbidden).JSON(resp.ErrorResponse(resp.ApiErrors("UserNotExist")))
	}

	//c.Locals("user", models.FilterUserRecord(&user))

	return c.Next()
}
