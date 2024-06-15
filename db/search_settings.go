package db

import (
	"time"
)

type SearchSettings struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	SearchOn  bool      `json:"searchOn"`
	AddNew    bool      `json:"addNew"`
	Amount    uint      `json:"amount"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Get is a method on the SearchSettings struct that retrieves the search settings from the database.
// It fetches the search settings with the ID of 1 and populates the SearchSettings struct with the retrieved values.
//
// This method does not take any parameters.
//
// Returns:
// error: An error object that describes an error that occurred during the method's execution.
func (s *SearchSettings) Get() error {
	err := DBConn.Where("id = 1").First(s).Error
	return err
}

// Update is a method on the SearchSettings struct that updates the search settings in the database.
// It updates the search_on, add_new, amount, and updated_at fields in the database with the values from the SearchSettings struct.
//
// This method does not take any parameters.
//
// Returns:
// error: An error object that describes an error that occurred during the method's execution.
func (s *SearchSettings) Update() error {
	tx := DBConn.Select("search_on", "add_new", "amount", "updated_at").Where("id = 1").Updates(&s)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
