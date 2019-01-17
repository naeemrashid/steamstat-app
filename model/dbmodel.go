package model

import (
	"time"
)

type Game struct {
	ID int64 `json:"primary_key:yes;auto_increment:false;column:id"`
	Name string
}
type GameDetail struct {
	GameID int64 `gorm:"primary_key:yes;auto_increment:false"`
	Game Game `gorm:"foreign_key:GameID"`
	Title string `json:"title"`
	Type                string          `json:"type"`
	IsFree              bool            `json:"is_free"`
	DetailedDescription string          `json:"detailed_description"`
	AboutTheGame        string          `json:"about_the_game"`
	ShortDescription    string          `json:"short_description"`
	SupportedLanguages  string          `json:"supported_languages"`
	Reviews             string          `json:"reviews"`
	HeaderImage         string          `json:"header_image"`
	Website             string          `json:"website"`
	Background          string          `json:"background"`

}
type Price struct {
	GameID int64 `gorm:"primary_key:yes;auto_increment:false"`
	Game Game `gorm:"foreign_key:GameID"`
	Currency        string `json:"currency"`
	Initial         int    `json:"initial"`
	Final           int    `json:"final"`
	DiscountPercent int    `json:"discount_percent"`
}
type PeakPlayer struct {
	GameID int64
	Game Game `gorm:"foreign_key:GameID"`
	PeakPlayers int64
	Date time.Time
}

