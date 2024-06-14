package db

import "time"

type SearchSetting struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Amount   int       `json:"amount"`
	SearchOn bool      `json:"searchOn"`
	AddNew   bool      `json:"addNew"`
	UpdateAt time.Time `json:"updateAt"`
}

func (s *SearchSetting) Get() error {
	err := DBConn.Where("id = 1").First(s).Error
	return err
}

func (s *SearchSetting) Update() error {
	tx := DBConn.Select("search_on", "amount", "add_new").Where("id = 1").Updates(s)
	// No update id
	return tx.Error
}
