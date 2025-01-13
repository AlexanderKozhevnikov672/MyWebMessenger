package internal

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ChatInvitation struct {
	ChatId int `gorm:"column:chat_id" json:"chat_id"`
	UserId int `gorm:"column:user_id" json:"user_id"`
}

func (ChatInvitation) TableName() string {
	return "msg.chat_invitation"
}

type ChatInvitationHandlers struct {
	db *gorm.DB
}

func (cih *ChatInvitationHandlers) AdminPage(c *fiber.Ctx) error {
	var chatInvitations []ChatInvitation
	err := cih.db.Model(&ChatInvitation{}).Find(&chatInvitations).Error
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	return c.Render("admin-chat-invitation", fiber.Map{"ChatInvitations": chatInvitations})
}

func (cih *ChatInvitationHandlers) AdminCreate(c *fiber.Ctx) error {
	var chatInvitation ChatInvitation
	if err := c.BodyParser(&chatInvitation); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	err := cih.db.Create(&chatInvitation).Error
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (cih *ChatInvitationHandlers) AdminDelete(c *fiber.Ctx) error {
	ids := c.Params("ids")
	if ids == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: "ids param is missing"})
	}

	parts := strings.Split(ids, "-")
	if len(parts) != 2 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: "ids param should match the format 'chat_id-user_id'"})
	}

	var chatId, userId int
	var err error
	chatId, err = strconv.Atoi(parts[0])
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}
	userId, err = strconv.Atoi(parts[1])
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	err = cih.db.Where(&ChatInvitation{ChatId: chatId, UserId: userId}).Delete(&ChatInvitation{}).Error
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
