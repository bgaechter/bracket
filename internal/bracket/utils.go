package bracket

import "math"

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
