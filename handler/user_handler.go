package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kdrkrgz/socalize/repository"
)

// GetUsers godoc
//
//	@Summary	Get Users
//	@Tags		Users
//	@Produce	json
//	@Success	200		{object}	[]users.User
//	@Security	Bearer
//	@Router		/users [GET]
func GetUsers(repo *repository.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := repo.GetUsers()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Something went wrong",
			})
		}

		return c.JSON(users)
	}
}
