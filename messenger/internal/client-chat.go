package internal

import (
	"container/list"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"gorm.io/gorm"
)

func (ca *ClientApi) ClientChatPage(c *fiber.Ctx) error {
	userId := ca.GetUserId(c)

	var chats []Chat
	if err := ca.db.Table(ChatMembership{}.TableName()+" AS cm").
		Select("c.id, c.name").
		Joins("LEFT JOIN "+Chat{}.TableName()+" AS c ON c.id = cm.chat_id").
		Where("cm.user_id = $1", userId).
		Scan(&chats).Error; err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	return c.Render("client-chat", fiber.Map{
		"Path":  ca.path,
		"Chats": chats,
	})
}

type Elem struct {
	C      *websocket.Conn
	UserId int
}

type MessageInfo struct {
	SenderNickname string `json:"sender_nickname"`
	MsgText        string `json:"msg_text"`
	Type           bool   `json:"type"`
}

func (ca *ClientApi) ClientChatSocket(c *fiber.Ctx) error {
	chatId, err := c.ParamsInt("chat_id")
	if err != nil {
		return err
	}

	userId := ca.GetUserId(c)

	var messagesInfo []MessageInfo
	if err = ca.db.Table(Message{}.TableName()+" AS m").
		Select("u.nickname AS sender_nickname, m.msg_text, CASE WHEN m.user_id = $1 THEN true ELSE false END AS type", userId).
		Joins("LEFT JOIN "+User{}.TableName()+" AS u ON u.id = m.user_id").
		Where("m.chat_id = $2", chatId).
		Order("m.id ASC").
		Scan(&messagesInfo).Error; err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	return websocket.New(func(c *websocket.Conn) {
		for _, mi := range messagesInfo {
			if err := c.WriteJSON(mi); err != nil {
				log.Println(err)
			}
		}

		ca.mu.Lock()
		l, ok := ca.conns[chatId]
		if !ok {
			l = list.New()
			ca.conns[chatId] = l
		}
		elem := l.PushBack(&Elem{C: c, UserId: userId})
		ca.mu.Unlock()

		defer func() {
			ca.mu.Lock()
			ca.conns[chatId].Remove(elem)
			ca.mu.Unlock()
		}()

		for {
			messageType, p, err := c.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}
			if messageType != websocket.TextMessage {
				log.Println("Not supported type of message:", messageType)
				break
			}

			msgText := string(p)
			var u User
			err = ca.db.Transaction(func(tx *gorm.DB) error {
				err = tx.Create(&Message{UserId: userId, ChatId: chatId, MsgText: msgText}).Error
				if err != nil {
					return err
				}

				err = tx.Where(&User{Id: userId}).First(&u).Error
				if err != nil {
					return err
				}

				return nil
			})
			if err != nil {
				log.Println(err)
				continue
			}

			var disabledConns []*list.Element
			ca.mu.Lock()
			for e := ca.conns[chatId].Front(); e != nil; e = e.Next() {
				nowConn := e.Value.(*Elem).C
				err = nowConn.WriteJSON(&MessageInfo{SenderNickname: u.Nickname, MsgText: msgText, Type: u.Id == e.Value.(*Elem).UserId})
				if err != nil {
					disabledConns = append(disabledConns, e)
				}
			}
			for _, e := range disabledConns {
				ca.conns[chatId].Remove(e)
			}
			ca.mu.Unlock()
		}
	})(c)
}

type InviteInfo struct {
	Nickname string `json:"nickname"`
}

func (ca *ClientApi) ClientChatInvite(c *fiber.Ctx) error {
	chatId, err := c.ParamsInt("chat_id")
	if err != nil {
		return err
	}

	userId := ca.GetUserId(c)

	var i InviteInfo
	err = c.BodyParser(&i)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: "invalid nickname"})
	}

	return ca.db.Transaction(func(tx *gorm.DB) error {
		var friend User
		err = tx.Where(&User{Nickname: i.Nickname}).First(&friend).Error
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: "invalid nickname"})
		}

		f := Friendship{UserId1: userId, UserId2: friend.Id}
		f.Normalize()
		var count int64
		err = tx.Model(&f).Where(&f).Count(&count).Error
		if err != nil || count == 0 {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: "not your friend"})
		}

		cm := ChatMembership{UserId: friend.Id, ChatId: chatId}
		err = tx.Model(&cm).Where(&cm).Count(&count).Error
		if err != nil || count != 0 {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: "already in this chat"})
		}

		err = tx.Create(&ChatInvitation{UserId: friend.Id, ChatId: chatId}).Error
		if err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}

func (ca *ClientApi) ClientChatLeave(c *fiber.Ctx) error {
	chatId, err := c.ParamsInt("chat_id")
	if err != nil {
		return err
	}

	userId := ca.GetUserId(c)

	return ca.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		err = tx.Model(&ChatMembership{}).Where(&ChatMembership{UserId: userId, ChatId: chatId}).Count(&count).Error
		if err != nil || count == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(&ErrorResponse{Error: "already left"})
		}

		err = tx.Where(&ChatMembership{UserId: userId, ChatId: chatId}).Delete(&ChatMembership{}).Error
		if err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}

type ChatInfo struct {
	Name string `json:"name"`
}

func (ca *ClientApi) ClientChatCreate(c *fiber.Ctx) error {
	userId := ca.GetUserId(c)

	var i ChatInfo
	err := c.BodyParser(&i)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorResponse{Error: "invalid name"})
	}

	return ca.db.Transaction(func(tx *gorm.DB) error {
		chat := Chat{Name: i.Name}
		err = tx.Create(&chat).Error
		if err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		err = tx.Create(&ChatMembership{UserId: userId, ChatId: chat.Id}).Error
		if err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}
