package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kdrkrgz/socalize/conf"
	"github.com/kdrkrgz/socalize/repository"
	"github.com/kdrkrgz/socalize/users"
	"github.com/kdrkrgz/socalize/utils"
	"time"
)

// SignUpUser godoc
// @Summary	Sign Up User
// @Tags		Auth
// @Accept		json
// @Produce	json
// @Success	200		{object}	string
func SignUp(repo *repository.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var payload *users.SignInInput

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid request body",
			})
		}
		user := repo.GetUserByEmail(payload.Email)
		if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid password",
			})
		}
		ttl, _ := time.ParseDuration(conf.Get("TokenExpiredIn"))

		token, err := utils.GenerateToken(ttl, user, conf.Get("TokenSecret"))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Something went wrong",
			})
		}
		cookie := new(fiber.Cookie)
		cookie.Name = "jwt"
		cookie.Value = token
		cookie.Expires = time.Now().Add(ttl)
		c.Cookie(cookie)
		return c.JSON(fiber.Map{
			"message": "Success",
			"token":   token,
		})
	}
}
