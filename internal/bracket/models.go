package bracket

import "time"

type User struct {
	Name   string `json:"name"`
	Points int    `json:"points"`
	Active bool   `json:"active"`
}

type Match struct {
	DateTime   time.Time `json:"date_time"`
	Team1      []*User   `json:"team1"`
	Team2      []*User   `json:"team2"`
	ScoreTeam1 int       `json:"score_team_1"`
	ScoreTeam2 int       `json:"score_team_2"`
}

type PlayForm struct {
	Users []string `form:"users[]"`
}

type SaveMatchForm struct {
	ScoreTeam1 int `form:"scoreTeam1"`
	ScoreTeam2 int `form:"scoreTeam2"`
}
