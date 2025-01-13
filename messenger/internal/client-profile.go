package internal

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (ca *ClientApi) ClientProfilePage(c *fiber.Ctx) error {
	userId := ca.GetUserId(c)

	return ca.db.Transaction(func(tx *gorm.DB) error {
		var u User
		err := tx.Where(&User{Id: userId}).First(&u).Error
		if err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		var friendCount1, friendCount2 int64
		err = tx.Model(&Friendship{}).Where(&Friendship{UserId1: userId}).Count(&friendCount1).Error
		if err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}
		err = tx.Model(&Friendship{}).Where(&Friendship{UserId2: userId}).Count(&friendCount2).Error
		if err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		var chatCount int64
		err = tx.Model(&ChatMembership{}).Where(&ChatMembership{UserId: userId}).Count(&chatCount).Error
		if err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		return c.Render("client-profile", fiber.Map{
			"Path":        ca.path,
			"Nickname":    u.Nickname,
			"FriendCount": friendCount1 + friendCount2,
			"ChatCount":   chatCount,
		})
	})
}

type PasswordInfo struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (ca *ClientApi) ClientProfileChangePassword(c *fiber.Ctx) error {
	userId := ca.GetUserId(c)

	var pi PasswordInfo
	err := c.BodyParser(&pi)
	if err != nil || pi.NewPassword == "" || pi.OldPassword == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: "invalid password info"})
	}

	return ca.db.Transaction(func(tx *gorm.DB) error {
		var u User
		err = tx.Where(&User{Id: userId}).First(&u).Error
		if err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		if !CheckPasswordHash(pi.OldPassword, u.PasswordHash) {
			return c.SendStatus(fiber.StatusForbidden)
		}

		var newHash string
		newHash, err = HashPassword(pi.NewPassword)
		if err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		err = tx.Model(&User{}).Where(&u).Updates(map[string]interface{}{"password_hash": newHash}).Error
		if err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}
