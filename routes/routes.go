package routes

import (
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

// render is a helper function that renders a component and sends the result to the client.
// It takes a Fiber context, a templ.Component, and an optional list of options.
// It returns an error if the rendering fails.
// The options are passed to the component handler.
// The component handler is created from the component using the templ.Handler function.
// The component handler is then passed to the HTTPHandler function from the adaptor package to create a Fiber handler.
// The Fiber handler is then called with the Fiber context to render the component.
// If the rendering is successful, the Fiber handler returns nil.
// If the rendering fails, the Fiber handler returns an error.
//
// Parameters:
// c *fiber.Ctx: The context of the request.
// component templ.Component: The component to render.
// options ...func(*templ.ComponentHandler): An optional list of options to pass to the component handler.
//
// Returns:
// error: An error object that describes an error that occurred during the rendering process.
func render(c *fiber.Ctx, component templ.Component, options ...func(*templ.ComponentHandler)) error {
	componentHandler := templ.Handler(component)
	for _, o := range options {
		o(componentHandler)
	}
	return adaptor.HTTPHandler(componentHandler)(c)
}

// SetRoutes is a function that sets up the routes for the Fiber application.
// It takes a pointer to a Fiber App as a parameter and does not return any values.
// The routes it sets up are:
// - GET /login: The login view
// - POST /login: The form submission from the login view
// - POST /logout: The logout action
// - POST /search: The search action
// - GET /: The dashboard view (requires authentication)
// - POST /: The form submission from the dashboard view (requires authentication)
//
// It also sets up a cache middleware for the /search route that caches responses for 30 minutes, unless the "noCache" query parameter is set to "true".
func SetRoutes(app *fiber.App) {
	app.Get("/login", LoginHandler)
	app.Post("/login", LoginPostHandler)
	app.Post("/logout", LogoutHandler)

	app.Post("/search", HandleSearch)
	app.Use("/search", cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("noCache") == "true"
		},
		Expiration:   30 * time.Minute,
		CacheControl: true,
	}))

	// app.Get("/create", func(c *fiber.Ctx) error {
	// 	u := &db.User{}
	// 	u.CreateAdmin()
	// 	return c.SendString("Admin Created")
	// })

	app.Get("/", AuthMiddleware, DashboardHandler)
	app.Post("/", AuthMiddleware, DashboardPostHandler)
}
