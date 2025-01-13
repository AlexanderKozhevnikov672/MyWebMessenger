package internal

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Friendship struct {
	UserId1 int `gorm:"column:user_id1" json:"user_id1"`
	UserId2 int `gorm:"column:user_id2" json:"user_id2"`
}

func (Friendship) TableName() string {
	return "msg.friendship"
}

func (f *Friendship) Normalize() {
	if f.UserId1 > 0 && f.UserId2 > 0 && f.UserId1 > f.UserId2 {
		f.UserId1, f.UserId2 = f.UserId2, f.UserId1
	}
}

type FriendshipHandlers struct {
	db *gorm.DB
}

func (fh *FriendshipHandlers) AdminPage(c *fiber.Ctx) error {
	var friendships []Friendship
	err := fh.db.Model(&Friendship{}).Find(&friendships).Error
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	return c.Render("admin-friendship", fiber.Map{"Friendships": friendships})
}

func (fh *FriendshipHandlers) AdminCreate(c *fiber.Ctx) error {
	var friendship Friendship
	if err := c.BodyParser(&friendship); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	if friendship.UserId1 == friendship.UserId2 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: "ids should be different"})
	}

	friendship.Normalize()
	err := fh.db.Create(&friendship).Error
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (fh *FriendshipHandlers) AdminDelete(c *fiber.Ctx) error {
	ids := c.Params("ids")
	if ids == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: "ids param is missing"})
	}

	parts := strings.Split(ids, "-")
	if len(parts) != 2 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: "ids param should match the format 'id1-id2'"})
	}

	var id1, id2 int
	var err error
	id1, err = strconv.Atoi(parts[0])
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}
	id2, err = strconv.Atoi(parts[1])
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	friendship := Friendship{UserId1: id1, UserId2: id2}
	friendship.Normalize()
	err = fh.db.Where(&friendship).Delete(&Friendship{}).Error
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
