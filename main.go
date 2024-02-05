package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

type User struct {
	Name   string `json:"name"`
	Points int    `json:"points"`
	Active bool   `json:"active"`
}

var users = []*User{}

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

func loadPlayers() {
	readFile, err := os.Open("./data/players.txt")

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		player := strings.Split(line, " ")
		fmt.Println(line)
		points, err := strconv.Atoi(player[1])
		if err != nil {
			fmt.Println(err)
		}
		users = append(users, &User{player[0], points, false})
	}
	readFile.Close()
}

func loadMatches() {
	readFile, err := os.Open("./data/matches.txt")

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		var m Match
		line := fileScanner.Text()
		b := []byte(line)
		err := json.Unmarshal(b, &m)
		if err != nil {
			fmt.Printf("Unable to marshal JSON due to %s", err)
		}
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
		} else {
			diff1, diff2 = CalculateNewElo(team1Points, team2Points, 0.0)
		}

		for _, u := range m.Team1 {
			for _, user := range users {
				if u.Name == user.Name {
					user.Points += diff1
				}
			}
		}
		for _, u := range m.Team2 {
			for _, user := range users {
				if u.Name == user.Name {
					user.Points += diff2
				}
			}
		}

	}
	readFile.Close()
}

// CalculateNewElo calculates and returns the new ELO ratings for two players.
func CalculateNewElo(ratingA, ratingB int, scoreA float64) (int, int) {
	// Player A wins
	// 1 for win, 0.5 for draw, 0 for loss
	const K = 32 // K-factor, which controls the rate at which ELO ratings can change.
	expectedScoreA := 1 / (1 + math.Pow(10, float64(ratingB-ratingA)/400))
	expectedScoreB := 1 - expectedScoreA

	// newRatingA := ratingA + int(float64(K)*(scoreA-expectedScoreA))
	// newRatingB := ratingB + int(float64(K)*((1-scoreA)-expectedScoreB))
	ratingDiffA := int(float64(K) * (scoreA - expectedScoreA))
	ratingDiffB := int(float64(K) * ((1 - scoreA) - expectedScoreB))

	return ratingDiffA, ratingDiffB
}

func (m *Match) createTeams(activePlayers []*User) {
	rand.Shuffle(len(activePlayers), func(i, j int) {
		activePlayers[i], activePlayers[j] = activePlayers[j], activePlayers[i]
	})

	for len(activePlayers) > 0 {
		var player *User
		player, activePlayers = activePlayers[0], activePlayers[1:]
		if len(activePlayers)%2 == 0 {
			m.Team1 = append(m.Team1, player)
		} else {
			m.Team2 = append(m.Team2, player)
		}
	}
}

func (m *Match) saveMatch() error {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile("./data/matches.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	b, err := json.Marshal(&m)
	if err != nil {
		fmt.Println(err)
	}
	if _, err := f.Write(b); err != nil {
		log.Fatal(err)
	}
	if _, err = f.Write([]byte("\n")); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func getUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, users)
}

func postUsers(c *gin.Context) {
	var u *User
	if err := c.BindJSON(&u); err != nil {
		return
	}
	users = append(users, u)
	c.IndentedJSON(http.StatusCreated, u)
}

func postPlay(c *gin.Context) {
	var playForm PlayForm
	c.Bind(&playForm)
	var activePlayers []*User
	for _, selectedUser := range playForm.Users {
		for _, user := range users {
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

func postMatch(c *gin.Context) {
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
		} else {
			diff1, diff2 = CalculateNewElo(team1Points, team2Points, 0.0)
		}

		for _, u := range m.Team1 {
			for _, user := range users {
				if u.Name == user.Name {
					user.Points += diff1
				}
			}
		}
		for _, u := range m.Team2 {
			for _, user := range users {
				if u.Name == user.Name {
					user.Points += diff2
				}
			}
		}
		session.Set("currentMatch", nil)
		m.saveMatch()
	}
}

func main() {
	loadPlayers()
	loadMatches()
	router := gin.Default()
	store := memstore.NewStore([]byte("ranker"))
	router.Use(sessions.Sessions("session", store))

	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)
		currentMatch := session.Get("currentMatch")
		var m Match
		if currentMatch != nil {

			b := []byte(fmt.Sprint(currentMatch))
			err := json.Unmarshal(b, &m)
			if err != nil {
				fmt.Printf("Unable to marshal JSON due to %s", err)
			}
			fmt.Println("Deserialized struct:", m)
		}
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":        "Main website",
			"users":        users,
			"currentMatch": m,
		})
	})
	router.POST("/", func(c *gin.Context) {
		postPlay(c)
		c.Request.URL.Path = "/"
		c.Request.Method = "GET"
		router.HandleContext(c)
	})
	router.POST("/postMatch", func(c *gin.Context) {
		postMatch(c)
		c.Request.URL.Path = "/"
		c.Request.Method = "GET"
		router.HandleContext(c)
	})
	router.GET("/users", getUsers)
	router.POST("/users", postUsers)

	router.Run("localhost:8080")
}
