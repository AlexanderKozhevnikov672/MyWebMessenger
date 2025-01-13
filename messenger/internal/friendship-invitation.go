package internal

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type FriendshipInvitation struct {
	FromUserId int `gorm:"column:from_user_id" json:"from_user_id"`
	ToUserId   int `gorm:"column:to_user_id" json:"to_user_id"`
}

func (FriendshipInvitation) TableName() string {
	return "msg.friendship_invitation"
}

type FriendshipInvitationHandlers struct {
	db *gorm.DB
}

func (fih *FriendshipInvitationHandlers) AdminPage(c *fiber.Ctx) error {
	var friendshipInvitations []FriendshipInvitation
	err := fih.db.Model(&FriendshipInvitation{}).Find(&friendshipInvitations).Error
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	return c.Render("admin-friendship-invitation", fiber.Map{"FriendshipInvitations": friendshipInvitations})
}

func (fih *FriendshipInvitationHandlers) AdminCreate(c *fiber.Ctx) error {
	var friendshipInvitation FriendshipInvitation
	if err := c.BodyParser(&friendshipInvitation); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	if friendshipInvitation.FromUserId == friendshipInvitation.ToUserId {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: "ids should be different"})
	}

	err := fih.db.Create(&friendshipInvitation).Error
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (fih *FriendshipInvitationHandlers) AdminDelete(c *fiber.Ctx) error {
	ids := c.Params("ids")
	if ids == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: "ids param is missing"})
	}

	parts := strings.Split(ids, "-")
	if len(parts) != 2 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: "ids param should match the format 'from_user_id-to_user_id'"})
	}

	var fromUserId, toUserId int
	var err error
	fromUserId, err = strconv.Atoi(parts[0])
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}
	toUserId, err = strconv.Atoi(parts[1])
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	err = fih.db.Where(&FriendshipInvitation{FromUserId: fromUserId, ToUserId: toUserId}).Delete(&FriendshipInvitation{}).Error
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
