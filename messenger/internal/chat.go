package internal

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Chat struct {
	Id   int    `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

func (Chat) TableName() string {
	return "msg.chat"
}

type ChatHandlers struct {
	db *gorm.DB
}

func (ch *ChatHandlers) AdminPage(c *fiber.Ctx) error {
	var chats []Chat
	err := ch.db.Model(&Chat{}).Find(&chats).Error
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	return c.Render("admin-chat", fiber.Map{"Chats": chats})
}

type ChatCreateInfo struct {
	Name string `json:"name"`
}

func (ch *ChatHandlers) AdminCreate(c *fiber.Ctx) error {
	var chatCreateInfo ChatCreateInfo
	if err := c.BodyParser(&chatCreateInfo); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	if chatCreateInfo.Name == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: "field name should not be empty"})
	}

	err := ch.db.Create(&Chat{Name: chatCreateInfo.Name}).Error
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (ch *ChatHandlers) AdminDelete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	err = ch.db.Where(&Chat{Id: id}).Delete(&Chat{}).Error
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

type ChatChangeInfo struct {
	Id      int    `json:"id"`
	NewName string `json:"new_name"`
}

func (ch *ChatHandlers) AdminChange(c *fiber.Ctx) error {
	var chatChangeInfo ChatChangeInfo
	if err := c.BodyParser(&chatChangeInfo); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	if chatChangeInfo.NewName == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: "field new password should not be empty"})
	}

	err := ch.db.Model(&Chat{}).Where(&Chat{Id: chatChangeInfo.Id}).Updates(map[string]interface{}{"name": chatChangeInfo.NewName}).Error
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
