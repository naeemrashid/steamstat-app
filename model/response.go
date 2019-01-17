package model

import (
	"github.com/naeemkhan12/golang-moving-average"
	"time"
)

type TransformedGameDetail struct {
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
type Trending struct {
	GameID int64 `json:"id"`
	GameTitle string `json:"game_title"`
	Change24hr float64 `json:"change_24_hr"`
	Change48hr []movingaverage.Values `json:"change_48_hr"`
	CurrentPlayers int64 `json:"current_players"`
}
type TopGameByCP struct {
	GameID int64 `json:"game_id"`
	GameTitle string `json:"game_title"`
	CurrentPlayers int64 `json:"current_players"`
	Last30Days []TimeSeries `json:"last_30_days"`
}
type TimeSeries struct {
	PeakPlayer int64 `json:"peak_player"`
	Time string `json:"time"`
}
type TopRecords struct {
	GameID int64 `json:"game_id"`
	GameTitle string `json:"game_title"`
	PeakPlayers int64 `json:"peak_players"`
	Date time.Time `json:"date"`
	Change48hr []movingaverage.Values `json:"change_48_hr"`
}
