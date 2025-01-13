package internal

import (
	"container/list"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ClientApi struct {
	app   *fiber.App
	db    *gorm.DB
	conns map[int]*list.List
	mu    sync.Mutex
	key   []byte
	path  string
}

func NewClientApi(dsn string, key []byte, path string) (*ClientApi, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	engine := html.New("template", ".html")
	app := fiber.New(fiber.Config{Views: engine})

	ca := &ClientApi{
		app:   app,
		db:    db,
		conns: make(map[int]*list.List),
		mu:    sync.Mutex{},
		key:   key,
		path:  path,
	}

	app.Get("/client/login", ca.ClientLoginPage)
	app.Post("/client/login/check", ca.ClientLoginCheck)
	app.Delete("/client/logout", ca.ClientAuthenticate, ca.ClientLogout)
	app.Post("/client/sign-up", ca.ClientSignIn)

	app.Get("/client/chat", ca.ClientAuthenticate, ca.ClientChatPage)
	app.Get("/client/chat/:chat_id", ca.ClientAuthenticate, ca.ClientChatSocket)
	app.Post("/client/chat/:chat_id/invite", ca.ClientAuthenticate, ca.ClientChatInvite)
	app.Delete("/client/chat/:chat_id/leave", ca.ClientAuthenticate, ca.ClientChatLeave)
	app.Post("/client/chat/create", ca.ClientAuthenticate, ca.ClientChatCreate)

	app.Get("/client/profile", ca.ClientAuthenticate, ca.ClientProfilePage)
	app.Post("/client/profile/change-password", ca.ClientAuthenticate, ca.ClientProfileChangePassword)

	app.Get("/client/management", ca.ClientAuthenticate, ca.ClientManagementPage)
	app.Delete("/client/management/friend/remove/:friend_id", ca.ClientAuthenticate, ca.ClientManagementFriendRemove)
	app.Post("/client/management/friend-invitation/approve/:friend_id", ca.ClientAuthenticate, ca.ClientManagementFriendInvitationApprove)
	app.Delete("/client/management/friend-invitation/decline/:friend_id", ca.ClientAuthenticate, ca.ClientManagementFriendInvitationDecline)
	app.Post("/client/management/chat-invitation/accept/:chat_id", ca.ClientAuthenticate, ca.ClientManagementChatInvitationAccept)
	app.Delete("/client/management/chat-invitation/reject/:chat_id", ca.ClientAuthenticate, ca.ClientManagementChatInvitationReject)
	app.Post("/client/management/friend-invitation/create", ca.ClientAuthenticate, ca.ClientManagementFriendInvitationCreate)

	return ca, nil
}

func (ca *ClientApi) Listen(addr string) error {
	return ca.app.Listen(addr)
}
