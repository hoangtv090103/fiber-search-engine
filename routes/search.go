package routes

import (
	"fiber-search-engine/db"

	"github.com/gofiber/fiber/v2"
)

type searchInput struct {
	Term string `json:"term"`
}

// HandleSearch is a Fiber handler function that processes the search request.
// It parses the request body into a searchInput struct and performs a full-text search on the SearchIndex table in the database.
// If there is an error parsing the request body, the search term is empty, or there is an error performing the search, it responds with a 500 status code and an error message.
// If the search is successful, it responds with a 200 status code and the search results.
//
// Parameters:
// c *fiber.Ctx: The context of the request.
//
// Returns:
// error: An error object that describes an error that occurred during the function's execution.
func HandleSearch(c *fiber.Ctx) error {
	input := searchInput{}
	if err := c.BodyParser(&input); err != nil {
		c.Status(500)
		c.Append("content-type", "application/json")
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Invalid input",
			"data":    nil,
		})
	}
	if input.Term == "" {
		c.Status(500)
		c.Append("content-type", "application/json")
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Invalid input",
			"data":    nil,
		})
	}
	idx := &db.SearchIndex{}
	data, err := idx.FullTextSearch(input.Term)
	if err != nil {
		c.Status(500)
		c.Append("content-type", "application/json")
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Invalid input",
			"data":    nil,
		})
	}
	c.Status(200)
	c.Append("content-type", "application/json")
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Search results",
		"data":    data,
	})
}
