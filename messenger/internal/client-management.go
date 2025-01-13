package internal

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (ca *ClientApi) ClientManagementPage(c *fiber.Ctx) error {
	userId := ca.GetUserId(c)

	return ca.db.Transaction(func(tx *gorm.DB) error {
		var friends []User
		if err := ca.db.Table(Friendship{}.TableName()+" AS f").
			Select("u.id, u.nickname").
			Joins("LEFT JOIN "+User{}.TableName()+" AS u ON u.id = CASE WHEN f.user_id1 = $1 THEN f.user_id2 ELSE f.user_id1 END").
			Where("f.user_id1 = $1 OR f.user_id2 = $1", userId).
			Scan(&friends).Error; err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		var friendInvitations []User
		if err := ca.db.Table(FriendshipInvitation{}.TableName()+" AS fi").
			Select("u.id, u.nickname").
			Joins("LEFT JOIN "+User{}.TableName()+" AS u ON u.id = fi.from_user_id").
			Where("fi.to_user_id = $1", userId).
			Scan(&friendInvitations).Error; err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		var chatsInvitations []Chat
		if err := ca.db.Table(ChatInvitation{}.TableName()+" AS ci").
			Select("c.id, c.name").
			Joins("LEFT JOIN "+Chat{}.TableName()+" AS c ON c.id = ci.chat_id").
			Where("ci.user_id = $1", userId).
			Scan(&chatsInvitations).Error; err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		return c.Render("client-management", fiber.Map{
			"Path":              ca.path,
			"Friends":           friends,
			"FriendInvitations": friendInvitations,
			"ChatsInvitations":  chatsInvitations,
		})
	})
}

func (ca *ClientApi) ClientManagementFriendRemove(c *fiber.Ctx) error {
	friendId, err := c.ParamsInt("friend_id")
	if err != nil {
		return err
	}

	userId := ca.GetUserId(c)

	f := Friendship{UserId1: userId, UserId2: friendId}
	f.Normalize()

	return ca.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		err = tx.Model(&Friendship{}).Where(&f).Count(&count).Error
		if err != nil || count == 0 {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		err = tx.Where(&f).Delete(&Friendship{}).Error
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}

func (ca *ClientApi) ClientManagementFriendInvitationApprove(c *fiber.Ctx) error {
	friendId, err := c.ParamsInt("friend_id")
	if err != nil {
		return err
	}

	userId := ca.GetUserId(c)

	return ca.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		err = tx.Model(&FriendshipInvitation{}).Where(&FriendshipInvitation{FromUserId: friendId, ToUserId: userId}).Count(&count).Error
		if err != nil || count == 0 {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		f := Friendship{UserId1: userId, UserId2: friendId}
		f.Normalize()
		err = tx.Create(&f).Error
		if err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		err = tx.Where(&FriendshipInvitation{FromUserId: friendId, ToUserId: userId}).Delete(&FriendshipInvitation{}).Error
		if err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}

func (ca *ClientApi) ClientManagementFriendInvitationDecline(c *fiber.Ctx) error {
	friendId, err := c.ParamsInt("friend_id")
	if err != nil {
		return err
	}

	userId := ca.GetUserId(c)

	return ca.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		err = tx.Model(&FriendshipInvitation{}).Where(&FriendshipInvitation{FromUserId: friendId, ToUserId: userId}).Count(&count).Error
		if err != nil || count == 0 {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		err = tx.Where(&FriendshipInvitation{FromUserId: friendId, ToUserId: userId}).Delete(&FriendshipInvitation{}).Error
		if err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}

func (ca *ClientApi) ClientManagementChatInvitationAccept(c *fiber.Ctx) error {
	chatId, err := c.ParamsInt("chat_id")
	if err != nil {
		return err
	}

	userId := ca.GetUserId(c)

	return ca.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		err = tx.Model(&ChatInvitation{}).Where(&ChatInvitation{UserId: userId, ChatId: chatId}).Count(&count).Error
		if err != nil || count == 0 {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		cm := ChatMembership{UserId: userId, ChatId: chatId}
		err = tx.Create(&cm).Error
		if err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		err = tx.Where(&ChatInvitation{UserId: userId, ChatId: chatId}).Delete(&ChatInvitation{}).Error
		if err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}

func (ca *ClientApi) ClientManagementChatInvitationReject(c *fiber.Ctx) error {
	chatId, err := c.ParamsInt("chat_id")
	if err != nil {
		return err
	}

	userId := ca.GetUserId(c)

	return ca.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		err = tx.Model(&ChatInvitation{}).Where(&ChatInvitation{UserId: userId, ChatId: chatId}).Count(&count).Error
		if err != nil || count == 0 {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		err = tx.Where(&ChatInvitation{UserId: userId, ChatId: chatId}).Delete(&ChatInvitation{}).Error
		if err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}

func (ca *ClientApi) ClientManagementFriendInvitationCreate(c *fiber.Ctx) error {
	userId := ca.GetUserId(c)

	var i InviteInfo
	err := c.BodyParser(&i)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&ErrorResponse{Error: "invalid friend invitation data"})
	}

	return ca.db.Transaction(func(tx *gorm.DB) error {
		var u User
		err = tx.Where(&User{Nickname: i.Nickname}).First(&u).Error
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&ErrorResponse{Error: "invalid nickname"})
		}

		f := Friendship{UserId1: userId, UserId2: u.Id}
		f.Normalize()
		var count int64
		err = tx.Model(&Friendship{}).Where(&f).Count(&count).Error
		if err != nil || count > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(&ErrorResponse{Error: "already friends"})
		}

		err = tx.Model(&FriendshipInvitation{}).Where(&FriendshipInvitation{FromUserId: userId, ToUserId: u.Id}).Count(&count).Error
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&ErrorResponse{Error: "invitation already exists"})
		}

		err = tx.Model(&FriendshipInvitation{}).Where(&FriendshipInvitation{FromUserId: u.Id, ToUserId: userId}).Count(&count).Error
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&ErrorResponse{Error: "opposite invitation exists"})
		}

		err = tx.Create(&FriendshipInvitation{FromUserId: userId, ToUserId: u.Id}).Error
		if err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}
