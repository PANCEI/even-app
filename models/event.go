package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	UserId      int       `json:"userId"`
	User        User      `gorm:"foreignKey:UserId"`
	Datetime    time.Time `json:"datetime" binding:"required"`
}

// // buat variabel untuk menampung data event
// var events []Event = []Event{}

// // fungsi untuk menyimpan data event
// func (e Event) Save() {
// 	events = append(events, e)
// }

// //	func Save(e Event) {
// //		events = append(events, e)
// //	}
// func GetAllEvents() []Event {
// 	return events
// }
