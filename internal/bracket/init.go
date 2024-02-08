package bracket

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var Users = []*User{}

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
		Users = append(Users, &User{player[0], points, false})
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

	}
	readFile.Close()
}

func init() {
	loadPlayers()
	loadMatches()
}
