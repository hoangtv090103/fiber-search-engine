package routes

import (
	utils "fiber-search-engine/utills"
	"fiber-search-engine/views"
	"fmt"
	"os"
	"strconv"
	"time"

	"fiber-search-engine/db"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func LoginHandler(c *fiber.Ctx) error {
	return render(c, views.Login())
}

type loginform struct { // Lowercase means doesn't export
	Email    string `form:"email"`
	Password string `form:"password"`
}

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
		return c.SendString("<h2>Error: Unauthorized</h2>")
	}

	signedToken, err := utils.CreateNewAuthToken(user.ID, user.Email, user.IsAdmin)

	if err != nil {
		c.Status(401) // 401: Unauthorized
		return c.SendString("<h2>Error: Something went wrong</h2>")
	}

	cookie := fiber.Cookie{
		Name:     "admin",
		Value:    signedToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	c.Append("HX-Redirect", "/")
	return c.SendStatus(200)
}

func LogoutHandler(c *fiber.Ctx) error {
	c.ClearCookie("admin")
	// c.ClearCookie("admin") - Clear the cookie named "admin"
	c.Set("HX-Redirect", "/login")
	// c.Set("HX-Redirect", "/login") - Set the header "HX-Redirect" to "/login"
	return c.SendStatus(200)
}

type AdminClaims struct {
	User                 string          `json:"user"`
	Id                   string          `json:"id"`
	jwt.RegisteredClaims `json:"claims"` // jwt.RegisteredClaims is a struct that contains the standard claims
}

// AuthMiddleware is a middleware function that checks if the user is authenticated as an admin.
// It retrieves the "admin" cookie from the request and validates it using JWT.
// If the cookie is missing or invalid, the function redirects the user to the login page.
// If the cookie is valid and the claims are of type AdminClaims, the function allows the request to proceed.
// Otherwise, it redirects the user to the login page.
//
// Parameters:
// - c: The fiber.Ctx object representing the current request and response.
//
// Returns:
// - An error if the cookie is missing, invalid, or the claims are not of type AdminClaims.
// - Otherwise, it allows the request to proceed.
func AuthMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("admin")
	// cookie := c.Cookies("admin") - Get the cookie named "admin"

	if cookie == "" {
		return c.Redirect("/login", 302)
		// return c.Redirect("/login", 302) - Redirect to "/login" with status code 302
	}

	token, err := jwt.ParseWithClaims(cookie, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return c.Redirect("/login", 302)
	}

	_, ok := token.Claims.(*AdminClaims) // Check if the claims are of type AdminClaims

	if ok && token.Valid {
		return c.Next()
	}
	return c.Redirect("/login", 302)
}

func DashboardHandler(c *fiber.Ctx) error {
	settings := db.SearchSetting{}
	err := settings.Get()
	if err != nil {
		c.Status(500)
		return c.SendString("<h2>Error: Cannot get settings</h2>")
	}
	amount := strconv.FormatUint(uint64(settings.Amount), 10)
	return render(c, views.Home(amount, settings.SearchOn, settings.AddNew))
}

type settingsform struct {
	Amount   int  `form:"amount"`
	SearchOn bool `form:"searchOn"`
	AddNew   bool `form:"addNew"`
}

func DashboardPostHandler(c *fiber.Ctx) error {
	input := settingsform{}
	if err := c.BodyParser(&input); err != nil {
		c.Status(500)
		return c.SendString("<h2>Error: cannot get settings</h2>")
	}

	addNew := false
	if input.AddNew {
		addNew = true
	}

	searchOn := false
	if input.SearchOn {
		searchOn = true
	}

	settings := &db.SearchSetting{}
	settings.Amount = input.Amount
	settings.SearchOn = searchOn
	settings.AddNew = addNew

	err := settings.Update()
	if err != nil {
		fmt.Println(err)
		c.Status(500)
		return c.SendString("<h2>Error: cannot update settings</h2>")
	}

	c.Append("HX-Refresh", "true")
	return c.SendStatus(200)
}
