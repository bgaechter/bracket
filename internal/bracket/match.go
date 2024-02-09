package bracket

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
)

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

		// Limit team size to two
		if len(m.Team2) == 2 && len(m.Team1) == 2 {
			break
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
