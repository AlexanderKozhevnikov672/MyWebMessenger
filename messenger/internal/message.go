package internal

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Message struct {
	Id      int    `gorm:"column:id"`
	UserId  int    `gorm:"column:user_id"`
	ChatId  int    `gorm:"column:chat_id"`
	MsgText string `gorm:"column:msg_text"`
}

func (Message) TableName() string {
	return "msg.message"
}

type MessageHandlers struct {
	db *gorm.DB
}

func (mh *MessageHandlers) AdminPage(c *fiber.Ctx) error {
	var messages []Message
	err := mh.db.Model(&Message{}).Find(&messages).Error
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	return c.Render("admin-message", fiber.Map{"Messages": messages})
}

func (mh *MessageHandlers) AdminDelete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	err = mh.db.Where(&Message{Id: id}).Delete(&Message{}).Error
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
