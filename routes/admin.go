package routes

import (
	"fiber-search-engine/db"
	"fiber-search-engine/utils"
	"fiber-search-engine/views"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AdminClaims struct {
	User                 string `json:"user"`
	Id                   string `json:"id"`
	jwt.RegisteredClaims `json:"claims"`
}

// DashboardHandler is a Fiber handler function that renders the dashboard view.
// It fetches the current search settings from the database and passes them to the view.
// If there is an error fetching the settings, it responds with a 500 status code and an error message.
//
// Parameters:
// c *fiber.Ctx: The context of the request.
//
// Returns:
// error: An error object that describes an error that occurred during the function's execution.
func DashboardHandler(c *fiber.Ctx) error {
	settings := &db.SearchSettings{}
	err := settings.Get()
	if err != nil {
		c.Status(500)
		return c.SendString("<h2>Error: Something went wrong</h2>")
	}
	amount := strconv.FormatUint(uint64(settings.Amount), 10)
	return render(c, views.Home(amount, settings.SearchOn, settings.AddNew))
}

type settingsform struct {
	Amount   uint   `form:"amount"`
	SearchOn string `form:"searchOn"`
	AddNew   string `form:"addNew"`
}

// DashboardPostHandler is a Fiber handler function that processes the form submission from the dashboard view.
// It parses the form data into a settingsform struct and updates the search settings in the database.
// If there is an error parsing the form data or updating the settings, it responds with a 500 status code and an error message.
// If the settings are updated successfully, it responds with a 200 status code and triggers a refresh of the dashboard view.
//
// Parameters:
// c *fiber.Ctx: The context of the request.
//
// Returns:
// error: An error object that describes an error that occurred during the function's execution.
func DashboardPostHandler(c *fiber.Ctx) error {
	input := settingsform{}
	if err := c.BodyParser(&input); err != nil {
		c.Status(500)
		return c.SendString("<h2>Error: Something went wrong</h2>")
	}
	// Convert checkbox 'on' values to boolean
	addNew := false
	if input.AddNew == "on" {
		addNew = true
	}
	searchOn := false
	if input.SearchOn == "on" {
		searchOn = true
	}
	settings := &db.SearchSettings{}
	settings.Amount = input.Amount
	settings.SearchOn = searchOn
	settings.AddNew = addNew
	err := settings.Update()
	if err != nil {
		fmt.Println(err)
		return c.SendString("<h2>Error: Something went wrong</h2>")
	}
	c.Append("HX-Refresh", "true")
	return c.SendStatus(200)
}

func LoginHandler(c *fiber.Ctx) error {
	return render(c, views.Login())
}

type loginform struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

// LoginPostHandler is a Fiber handler function that processes the form submission from the login view.
// It parses the form data into a loginform struct and attempts to log in the user as an admin.
// If there is an error parsing the form data or the login credentials are incorrect, it responds with an appropriate status code and an error message.
// If the login is successful, it creates a new auth token, sets a cookie with the token, and redirects the user to the home page.
//
// Parameters:
// c *fiber.Ctx: The context of the request.
//
// Returns:
// error: An error object that describes an error that occurred during the function's execution.
func LoginPostHandler(c *fiber.Ctx) error {
	input := loginform{}
	if err := c.BodyParser(&input); err != nil {
		c.Status(500)
		return c.SendString("<h2>Error: Something went wrong</h2>")
	}
	user := &db.User{}
	user, err := user.LoginAsAdmin(input.Email, input.Password)
	if err != nil {
		c.Status(401)
		c.Append("content-type", "text/html")
		return c.SendString("<h2>Error: Unauthorised</h2>")
	}

	signedToken, err := utils.CreateNewAuthToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		c.Status(500)
		return c.SendString("<h2>Error:Something went wrong logging in, please try again.</h2>")
	}

	// Create and set the cookie
	cookie := fiber.Cookie{
		Name:     "admin",
		Value:    signedToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true, // Meant only for the server
	}
	c.Cookie(&cookie)
	c.Append("HX-Redirect", "/")
	return c.SendStatus(200)
}

// LogoutHandler is a Fiber handler function that logs out the current user.
// It clears the "admin" cookie, sets the "HX-Redirect" header to "/login", and sends a 200 status code.
// After this function is called, the user will no longer be authenticated and will be redirected to the login page.
//
// Parameters:
// c *fiber.Ctx: The context of the request.
//
// Returns:
// error: An error object that describes an error that occurred during the function's execution.
func LogoutHandler(c *fiber.Ctx) error {
	c.ClearCookie("admin")
	c.Set("HX-Redirect", "/login")
	return c.SendStatus(200)
}

// AuthMiddleware is a Fiber middleware function that checks if the user is authenticated.
// It retrieves the "admin" cookie, parses the JWT token from the cookie, and checks if the token is valid.
// If the cookie does not exist, the token is invalid, or the token's claims are incorrect, it redirects the user to the login page.
// If the token is valid and the claims are correct, it allows the request to proceed to the next handler.
//
// Parameters:
// c *fiber.Ctx: The context of the request.
//
// Returns:
// error: An error object that describes an error that occurred during the function's execution.
func AuthMiddleware(c *fiber.Ctx) error {
	// Get the cookie by name
	cookie := c.Cookies("admin")
	if cookie == "" {
		return c.Redirect("/login", 302)
	}
	// Parse the cookie & check for errors
	token, err := jwt.ParseWithClaims(cookie, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return c.Redirect("/login", 302)
	}
	// Parse the custom claims & check jwt is valid
	_, ok := token.Claims.(*AdminClaims)
	if ok && token.Valid {
		return c.Next()
	}
	return c.Redirect("/login", 302)
}
