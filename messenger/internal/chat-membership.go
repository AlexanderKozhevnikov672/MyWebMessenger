package internal

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ChatMembership struct {
	ChatId int `gorm:"column:chat_id" json:"chat_id"`
	UserId int `gorm:"column:user_id" json:"user_id"`
}

func (ChatMembership) TableName() string {
	return "msg.chat_membership"
}

type ChatMembershipHandlers struct {
	db *gorm.DB
}

func (cmh *ChatMembershipHandlers) AdminPage(c *fiber.Ctx) error {
	var chatMemberships []ChatMembership
	err := cmh.db.Model(&ChatMembership{}).Find(&chatMemberships).Error
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	return c.Render("admin-chat-membership", fiber.Map{"ChatMemberships": chatMemberships})
}

func (cmh *ChatMembershipHandlers) AdminCreate(c *fiber.Ctx) error {
	var chatMembership ChatMembership
	if err := c.BodyParser(&chatMembership); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	err := cmh.db.Create(&chatMembership).Error
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (cmh *ChatMembershipHandlers) AdminDelete(c *fiber.Ctx) error {
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

	err = cmh.db.Where(&ChatMembership{ChatId: chatId, UserId: userId}).Delete(&ChatMembership{}).Error
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
