package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AdminApi struct {
	app *fiber.App
	db  *gorm.DB
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewAdminApi(dsn string) (*AdminApi, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	engine := html.New("template", ".html")
	app := fiber.New(fiber.Config{Views: engine})

	func(uh *UserHandlers) {
		app.Get("/admin/user", uh.AdminPage)
		app.Post("/admin/user/create", uh.AdminCreate)
		app.Delete("/admin/user/delete/:id", uh.AdminDelete)
		app.Post("/admin/user/change", uh.AdminChange)
	}(&UserHandlers{db: db})

	func(mh *MessageHandlers) {
		app.Get("/admin/message", mh.AdminPage)
		app.Delete("/admin/message/delete/:id", mh.AdminDelete)
	}(&MessageHandlers{db: db})

	func(fh *FriendshipHandlers) {
		app.Get("/admin/friendship", fh.AdminPage)
		app.Post("/admin/friendship/create", fh.AdminCreate)
		app.Delete("/admin/friendship/delete/:ids", fh.AdminDelete)
	}(&FriendshipHandlers{db: db})

	func(fih *FriendshipInvitationHandlers) {
		app.Get("/admin/friendship-invitation", fih.AdminPage)
		app.Post("/admin/friendship-invitation/create", fih.AdminCreate)
		app.Delete("/admin/friendship-invitation/delete/:ids", fih.AdminDelete)
	}(&FriendshipInvitationHandlers{db: db})

	func(ch *ChatHandlers) {
		app.Get("/admin/chat", ch.AdminPage)
		app.Post("/admin/chat/create", ch.AdminCreate)
		app.Delete("/admin/chat/delete/:id", ch.AdminDelete)
		app.Post("/admin/chat/change", ch.AdminChange)
	}(&ChatHandlers{db: db})

	func(cmh *ChatMembershipHandlers) {
		app.Get("/admin/chat-membership", cmh.AdminPage)
		app.Post("/admin/chat-membership/create", cmh.AdminCreate)
		app.Delete("/admin/chat-membership/delete/:ids", cmh.AdminDelete)
	}(&ChatMembershipHandlers{db: db})

	func(cih *ChatInvitationHandlers) {
		app.Get("/admin/chat-invitation", cih.AdminPage)
		app.Post("/admin/chat-invitation/create", cih.AdminCreate)
		app.Delete("/admin/chat-invitation/delete/:ids", cih.AdminDelete)
	}(&ChatInvitationHandlers{db: db})

	return &AdminApi{app: app, db: db}, nil
}

func (aa *AdminApi) Listen(addr string) error {
	return aa.app.Listen(addr)
}
