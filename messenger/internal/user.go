package internal

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id           int    `gorm:"column:id"`
	Nickname     string `gorm:"column:nickname"`
	PasswordHash string `gorm:"column:password_hash"`
}

func (User) TableName() string {
	return "msg.user"
}

type UserHandlers struct {
	db *gorm.DB
}

func (uh *UserHandlers) AdminPage(c *fiber.Ctx) error {
	var users []User
	err := uh.db.Model(&User{}).Find(&users).Error
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	return c.Render("admin-user", fiber.Map{"Users": users})
}

type UserCreateInfo struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (uh *UserHandlers) AdminCreate(c *fiber.Ctx) error {
	var userCreateInfo UserCreateInfo
	if err := c.BodyParser(&userCreateInfo); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	if userCreateInfo.Nickname == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: "field nickname should not be empty"})
	}
	if userCreateInfo.Password == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: "field password should not be empty"})
	}

	passwordHash, err := HashPassword(userCreateInfo.Password)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	err = uh.db.Create(&User{Nickname: userCreateInfo.Nickname, PasswordHash: passwordHash}).Error
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (uh *UserHandlers) AdminDelete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	err = uh.db.Where(&User{Id: id}).Delete(&User{}).Error
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

type UserChangeInfo struct {
	Id          int    `json:"id"`
	NewNickname string `json:"new_nickname"`
	NewPassword string `json:"new_password"`
}

func (uh *UserHandlers) AdminChange(c *fiber.Ctx) error {
	var userChangeInfo UserChangeInfo
	if err := c.BodyParser(&userChangeInfo); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
	}

	if userChangeInfo.NewNickname == "" && userChangeInfo.NewPassword == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: "new field should not be empty"})
	}

	if userChangeInfo.NewNickname != "" {
		err := uh.db.Model(&User{}).Where(&User{Id: userChangeInfo.Id}).Updates(map[string]interface{}{"nickname": userChangeInfo.NewNickname}).Error
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
		}
	}

	if userChangeInfo.NewPassword != "" {
		newPasswordHash, err := HashPassword(userChangeInfo.NewPassword)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
		}

		err = uh.db.Model(&User{}).Where(&User{Id: userChangeInfo.Id}).Updates(map[string]interface{}{"password_hash": newPasswordHash}).Error
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: err.Error()})
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}
