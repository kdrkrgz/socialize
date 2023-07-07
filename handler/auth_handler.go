package handler

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kdrkrgz/socalize/conf"
	"github.com/kdrkrgz/socalize/repository"
	"github.com/kdrkrgz/socalize/users"
	"github.com/kdrkrgz/socalize/utils"
	"golang.org/x/crypto/bcrypt"
)

const (
	cookieName       = "jwt"
	tokenMessage     = "Success"
	errorMessage     = "Something went wrong"
	passwordMatch    = "Password and Password Confirmation must be the same"
	invalidReqBody   = "Invalid request body"
	invalidCreds     = "Invalid credentials"
	userCreateFailed = "Failed to create user"
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
				"message": invalidReqBody,
			})
		}

		user := repo.GetUserByEmail(payload.Email)
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": invalidCreds,
			})
		}
		ttl, _ := time.ParseDuration(conf.Get("TokenExpiredIn"))

		token, err := utils.GenerateToken(ttl, user, conf.Get("TokenSecret"))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": errorMessage,
			})
		}
		cookie := &fiber.Cookie{
			Name:     cookieName,
			Value:    token,
			Expires:  time.Now().Add(ttl),
			SameSite: fiber.CookieSameSiteStrictMode,
			HTTPOnly: true,
			Secure:   true,
			MaxAge:   int(ttl.Seconds()),
		}
		c.Cookie(cookie)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": tokenMessage,
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
			return fiber.NewError(fiber.StatusBadRequest, invalidReqBody)
		}

		if errors := users.ValidateStruct(payload); errors != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(errors)
		}

		if payload.Password != payload.PasswordConfirm {
			return fiber.NewError(fiber.StatusBadRequest, passwordMatch)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, errorMessage)
		}

		newUser := &users.User{
			Email:     strings.ToLower(payload.Email),
			Password:  string(hashedPassword),
			FirstName: payload.FirstName,
			LastName:  payload.LastName,
			Username:  payload.Username,
		}

		if err := repo.CreateUser(newUser); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, userCreateFailed)
		}

		ttl, _ := time.ParseDuration(conf.Get("TokenExpiredIn"))
		token, err := utils.GenerateToken(ttl, newUser, conf.Get("TokenSecret"))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, errorMessage)
		}

		cookie := fiber.Cookie{
			Name:     cookieName,
			Value:    token,
			Expires:  time.Now().Add(ttl),
			SameSite: fiber.CookieSameSiteStrictMode,
			HTTPOnly: true,
			Secure:   true,
			MaxAge:   int(ttl.Seconds()),
		}
		ctx.Cookie(&cookie)

		return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": tokenMessage,
			"data": fiber.Map{
				"user": users.FilterUserRecord(newUser),
			},
		})
	}
}
