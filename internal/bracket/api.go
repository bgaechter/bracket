package bracket

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetIndex(c *gin.Context) {
	session := sessions.Default(c)
	currentMatch := session.Get("currentMatch")
	var m Match
	if currentMatch != nil {

		b := []byte(fmt.Sprint(currentMatch))
		err := json.Unmarshal(b, &m)
		if err != nil {
			fmt.Printf("Unable to marshal JSON due to %s", err)
		}
	}

	sort.Slice(Users, func(i, j int) bool {
		return Users[i].Points > Users[j].Points
	})

	winA, draw, winB := m.eloChange()

	// only Display 20 last matches
	index := len(Matches)
	if index > 20 {
		index = 20
	}
	sort.Slice(Matches, func(i, j int) bool {
		return Matches[i].DateTime.Unix() > Matches[j].DateTime.Unix()
	})

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"users":        Users,
		"currentMatch": m,
		"winA":         winA,
		"draw":         draw,
		"winB":         winB,
		"matches":      Matches[:index],
	})
}

func PostNewPlayer(c *gin.Context) {
	var newPlayerForm NewPlayerForm
	c.Bind(&newPlayerForm)
	u := &User{newPlayerForm.Name, 800, false}
	Users = append(Users, u)
	f, err := os.OpenFile(Config.DataDir+"/players.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = f.Write([]byte(u.Name + " 800\n")); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func PostPlay(c *gin.Context) {
	var playForm PlayForm
	c.Bind(&playForm)
	var activePlayers []*User
	for _, user := range Users {
		user.Active = false
	}
	for _, selectedUser := range playForm.Users {
		for _, user := range Users {
			if selectedUser == user.Name {
				user.Active = true
				activePlayers = append(activePlayers, user)
			}
		}
	}
	m := Match{}
	m.createTeams(activePlayers)
	session := sessions.Default(c)
	b, err := json.Marshal(&m)
	if err != nil {
		fmt.Println(err)
	}

	session.Set("currentMatch", string(b))
	session.Save()
}

func PostMatch(c *gin.Context) {
	var saveMatchForm SaveMatchForm
	c.Bind(&saveMatchForm)

	session := sessions.Default(c)
	currentMatch := session.Get("currentMatch")
	var m Match
	if currentMatch != nil {
		b := []byte(fmt.Sprint(currentMatch))
		err := json.Unmarshal(b, &m)
		if err != nil {
			fmt.Printf("Unable to marshal JSON due to %s", err)
		}
		m.ScoreTeam1 = saveMatchForm.ScoreTeam1
		m.ScoreTeam2 = saveMatchForm.ScoreTeam2
		m.DateTime = time.Now()
		var team1Points int
		var team2Points int
		for _, u := range m.Team1 {
			team1Points += u.Points
		}
		for _, u := range m.Team2 {
			team2Points += u.Points
		}
		var diff1, diff2 int
		if m.ScoreTeam1 > m.ScoreTeam2 {
			diff1, diff2 = CalculateNewElo(team1Points, team2Points, 1.0)
		} else if m.ScoreTeam1 == m.ScoreTeam2 {
			diff1, diff2 = CalculateNewElo(team1Points, team2Points, 0.5)
		} else {
			diff1, diff2 = CalculateNewElo(team1Points, team2Points, 0.0)
		}
		for _, u := range m.Team1 {
			for _, user := range Users {
				if u.Name == user.Name {
					user.Points += diff1
				}
			}
		}
		for _, u := range m.Team2 {
			for _, user := range Users {
				if u.Name == user.Name {
					user.Points += diff2
				}
			}
		}
		m.PointsChange = int(math.Abs(float64(diff1)))
		Matches = append(Matches, &m)
		session.Set("currentMatch", nil)
		m.saveMatch()
	}
}
