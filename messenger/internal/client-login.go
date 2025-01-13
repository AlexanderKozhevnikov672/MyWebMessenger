package internal

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (ca *ClientApi) ClientLoginPage(c *fiber.Ctx) error {
	return c.Render("client-login", fiber.Map{
		"Path": ca.path,
	})
}

func GenerateJWT(userId int, key []byte) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

func ParseToken(tokenString string, key []byte) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return key, nil
	})
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type LoginRequest struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

func (ca *ClientApi) ClientLoginCheck(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&ErrorResponse{Error: "invalid login request data"})
	}

	var u User
	if err := ca.db.Where(&User{Nickname: req.Nickname}).First(&u).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(&ErrorResponse{Error: "invalid nickname or password"})
	}

	if !CheckPasswordHash(req.Password, u.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(&ErrorResponse{Error: "invalid nickname or password"})
	}

	token, err := GenerateJWT(u.Id, ca.key)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		HTTPOnly: true,
		Secure:   false,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
		SameSite: fiber.CookieSameSiteStrictMode,
	})

	return c.SendStatus(fiber.StatusNoContent)
}

func (ca *ClientApi) ClientAuthenticate(c *fiber.Ctx) error {
	tokenString := c.Cookies("token")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(&ErrorResponse{Error: "no token provided"})
	}

	token, err := ParseToken(tokenString, ca.key)
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(&ErrorResponse{Error: "ivalide or expired token provided"})
	}

	userId := int(token.Claims.(jwt.MapClaims)["user_id"].(float64))
	c.Locals("userId", userId)

	return c.Next()
}

func (ca *ClientApi) ClientLogout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,
		Path:     "/",
		Expires:  time.Now(),
		SameSite: fiber.CookieSameSiteStrictMode,
	})

	return c.SendStatus(fiber.StatusOK)
}

func (ca *ClientApi) GetUserId(c *fiber.Ctx) int {
	return c.Locals("userId").(int)
}

func (ca *ClientApi) ClientSignIn(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&ErrorResponse{Error: "invalid sign in data"})
	}

	return ca.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&User{}).Where(&User{Nickname: req.Nickname}).Count(&count).Error; err != nil || count != 0 {
			return c.Status(fiber.StatusBadRequest).JSON(&ErrorResponse{Error: "this nickname is already in use"})
		}

		passwordHash, err := HashPassword(req.Password)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&ErrorResponse{Error: "invalid password format"})
		}

		err = tx.Create(&User{Nickname: req.Nickname, PasswordHash: passwordHash}).Error
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}
