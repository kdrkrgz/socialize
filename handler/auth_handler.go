package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kdrkrgz/socalize/conf"
	"github.com/kdrkrgz/socalize/repository"
	"github.com/kdrkrgz/socalize/users"
	"github.com/kdrkrgz/socalize/utils"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

// SignIn godoc
//
//	@Summary	Login
//	@Tags		Auth
//	@Produce	json
//
// @Param data body users.SignInInput true "Login"
// @Success	200		{object}	string
// @Router		/auth/signin [POST]
func SignIn(repo *repository.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var payload *users.SignInInput

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid request body",
			})
		}

		user := repo.GetUserByEmail(payload.Email)
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid credentials",
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
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Success",
			"token":   token,
		})
	}
}

// SignUp godoc
//
//	@Summary	Register
//	@Tags		Auth
//	@Produce	json
//
// @Param data body users.SignUpInput true "Register"
//
//	@Success	201		{object}	string
//	@Router		/auth/signup [POST]
func SignUp(repo *repository.Repository) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var payload *users.SignUpInput
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid request body",
			})
		}
		errors := users.ValidateStruct(payload)
		if errors != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(errors)
		}
		if payload.Password != payload.PasswordConfirm {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Password and Password Confirmation must be the same",
			})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Something went wrong",
			})
		}
		newUser := &users.User{
			Email:     strings.ToLower(payload.Email),
			Password:  string(hashedPassword),
			FirstName: payload.FirstName,
			LastName:  payload.LastName,
			Username:  payload.Username,
		}
		result := repo.CreateUser(newUser)
		if result.Error != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Something went wrong",
			})
		}
		return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Success",
			"data": fiber.Map{
				"user": users.FilterUserRecord(newUser),
			},
		})
	}
}
