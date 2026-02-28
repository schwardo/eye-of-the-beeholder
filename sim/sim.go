package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Attribute represents one of the 6 attribute categories
type Attribute int

const (
	Texture Attribute = iota // Fuzzy=false, Shiny=true
	Antennae                 // Feathered=false, Whips=true
	Weapon                   // Stinger=false, Mandibles=true
	Pattern                  // Striped=false, Solid=true
	Wings                    // Sleek=false, Flutter=true
	Payload                  // Honey=false, Pollen=true
)

var attributeNames = []string{"Texture", "Antennae", "Weapon", "Pattern", "Wings", "Payload"}
var attributeValues = [][]string{
	{"Fuzzy", "Shiny"},
	{"Feathered", "Whips"},
	{"Stinger", "Mandibles"},
	{"Striped", "Solid"},
	{"Sleek", "Flutter"},
	{"Honey", "Pollen"},
}

// Card represents a single card with 6 boolean attributes
type Card struct {
	Attributes [6]bool // Each index corresponds to an Attribute
}

// String returns a human-readable representation of the card
func (c Card) String() string {
	result := "["
	for i, attr := range c.Attributes {
		if i > 0 {
			result += ", "
		}
		valueIdx := 0
		if attr {
			valueIdx = 1
		}
		result += attributeValues[i][valueIdx]
	}
	result += "]"
	return result
}

// Matches returns true if the card matches the given attribute value
func (c Card) Matches(attr Attribute, value bool) bool {
	return c.Attributes[attr] == value
}

// Player represents a player in the game
type Player struct {
	ID         int
	Hand       []Card
	TricksWon  int
	ScorePile  []Card
}

// ProtocolBoard represents the 6 slots around the Queen's Favor
type ProtocolBoard struct {
	Slots [6]*AttributeToken // Index 0 = Slot 1 (highest priority), Index 5 = Slot 6 (lowest priority)
}

// AttributeToken represents a token placed on the board
type AttributeToken struct {
	Attribute Attribute
	Value     bool // false or true
}

func (at AttributeToken) String() string {
	valueIdx := 0
	if at.Value {
		valueIdx = 1
	}
	return fmt.Sprintf("%s=%s", attributeNames[at.Attribute], attributeValues[at.Attribute][valueIdx])
}

// Game represents the entire game state
type Game struct {
	Players            []*Player
	NumPlayers         int
	Deck               []Card
	Box                []Card // Cards not dealt this hand
	Board              ProtocolBoard
	CurrentLeader      int
	TrickNumber        int
	HandNumber         int
	Verbose            bool
	SuddenDeath        bool // True when multiple players tied at 10+
	Stats              *GameStats
	LastTrickWinner    int
	ConsecutiveWins    int
}

// GameStats tracks statistics across games
type GameStats struct {
	GamesPlayed      int
	WinsByPlayer     []int
	TricksByPlayer   []int
	TwoStreaksByPlayer []int // Number of times each player won 2+ tricks in a row
	ThreeStreaksByPlayer []int // Number of times each player won 3+ tricks in a row
}

// CreateDeck creates all 64 unique cards
func CreateDeck() []Card {
	deck := make([]Card, 64)
	for i := 0; i < 64; i++ {
		var card Card
		for j := 0; j < 6; j++ {
			card.Attributes[j] = (i & (1 << j)) != 0
		}
		deck[i] = card
	}
	return deck
}

// ShuffleDeck shuffles the deck in place
func ShuffleDeck(deck []Card) {
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
}

// NewGame creates and initializes a new game
func NewGame(numPlayers int, verbose bool) *Game {
	if numPlayers < 2 || numPlayers > 5 {
		panic(fmt.Sprintf("Invalid number of players: %d. Must be 2-5.", numPlayers))
	}

	game := &Game{
		NumPlayers:      numPlayers,
		Players:         make([]*Player, numPlayers),
		Verbose:         verbose,
		CurrentLeader:   rand.Intn(numPlayers),
		HandNumber:      1,
		LastTrickWinner: -1,
		ConsecutiveWins: 0,
	}

	// Initialize players
	for i := 0; i < numPlayers; i++ {
		game.Players[i] = &Player{ID: i}
	}

	// Create and shuffle deck
	game.Deck = CreateDeck()
	game.DealNewHand()

	return game
}

// NewStats creates a new statistics tracker
func NewStats(numPlayers int) *GameStats {
	return &GameStats{
		WinsByPlayer:       make([]int, numPlayers),
		TricksByPlayer:     make([]int, numPlayers),
		TwoStreaksByPlayer: make([]int, numPlayers),
		ThreeStreaksByPlayer: make([]int, numPlayers),
	}
}

// DealNewHand shuffles all cards and deals a new hand
func (g *Game) DealNewHand() {
	if g.Verbose {
		fmt.Printf("\n=== HAND %d ===\n", g.HandNumber)
	}

	// Collect all cards
	allCards := make([]Card, 0, 64)
	allCards = append(allCards, g.Deck...)
	allCards = append(allCards, g.Box...)
	for _, player := range g.Players {
		allCards = append(allCards, player.Hand...)
		allCards = append(allCards, player.ScorePile...)
		player.Hand = nil
		player.ScorePile = nil
	}

	// Clear deck and box to prevent duplication in future hands
	g.Deck = nil
	g.Box = nil

	// Shuffle all cards
	ShuffleDeck(allCards)

	// Deal 7 cards to each player regardless of player count
	cardsPerPlayer := 7

	// Deal cards to each player
	cardIndex := 0
	for i := 0; i < cardsPerPlayer; i++ {
		for j := 0; j < g.NumPlayers; j++ {
			g.Players[j].Hand = append(g.Players[j].Hand, allCards[cardIndex])
			cardIndex++
		}
	}

	// Remaining cards go to The Box
	g.Box = allCards[cardIndex:]
	g.TrickNumber = 1

	if g.Verbose {
		cardsPerPlayer := len(g.Players[0].Hand)
		fmt.Printf("Dealt %d cards to each of %d players. %d cards in The Box.\n", cardsPerPlayer, g.NumPlayers, len(g.Box))
		fmt.Printf("Leader for first trick: Player %d\n", g.CurrentLeader)
	}

	// Run draft phase for this hand
	g.RunDraftPhase()
}

// RunDraftPhase executes the draft phase where players place tokens
func (g *Game) RunDraftPhase() {
	if g.Verbose {
		fmt.Printf("\n--- Drafting the Queen's Favor ---\n")
	}

	// Clear the board
	for i := range g.Board.Slots {
		g.Board.Slots[i] = nil
	}

	// Available tokens (all 6 attributes)
	availableTokens := []Attribute{Texture, Antennae, Weapon, Pattern, Wings, Payload}

	// Build counter-clockwise order from QF holder; QF holder drafts last.
	leader := g.CurrentLeader
	ccw := make([]int, g.NumPlayers)
	for i := 1; i < g.NumPlayers; i++ {
		ccw[i-1] = (leader - i + g.NumPlayers) % g.NumPlayers
	}
	ccw[g.NumPlayers-1] = leader

	switch g.NumPlayers {
	case 2:
		// 2-player: alternate from Slot 6 down. First drafter gets 6,4,2; QF holder gets 5,3,1.
		g.draftTokens(ccw[0], &availableTokens, []int{5})
		g.draftTokens(ccw[1], &availableTokens, []int{4})
		g.draftTokens(ccw[0], &availableTokens, []int{3})
		g.draftTokens(ccw[1], &availableTokens, []int{2})
		g.draftTokens(ccw[0], &availableTokens, []int{1})
		g.draftTokens(ccw[1], &availableTokens, []int{0})

	case 3:
		// 3-player: first drafter gets 6,5; next gets 4,3; QF holder gets 2,1.
		g.draftTokens(ccw[0], &availableTokens, []int{5, 4})
		g.draftTokens(ccw[1], &availableTokens, []int{3, 2})
		g.draftTokens(ccw[2], &availableTokens, []int{1, 0})

	case 4:
		// 4-player: first gets 6,5; next gets 4,3; next gets 2; QF holder gets 1.
		g.draftTokens(ccw[0], &availableTokens, []int{5, 4})
		g.draftTokens(ccw[1], &availableTokens, []int{3, 2})
		g.draftTokens(ccw[2], &availableTokens, []int{1})
		g.draftTokens(ccw[3], &availableTokens, []int{0})

	case 5:
		// 5-player: first gets 6,5; remaining each get one slot (4,3,2,1).
		g.draftTokens(ccw[0], &availableTokens, []int{5, 4})
		g.draftTokens(ccw[1], &availableTokens, []int{3})
		g.draftTokens(ccw[2], &availableTokens, []int{2})
		g.draftTokens(ccw[3], &availableTokens, []int{1})
		g.draftTokens(ccw[4], &availableTokens, []int{0})
	}

	if g.Verbose {
		fmt.Println("\nQueen's Favor:")
		for i := 0; i < 6; i++ {
			if g.Board.Slots[i] != nil {
				fmt.Printf("  Slot %d: %s\n", i+1, g.Board.Slots[i])
			}
		}
	}
}

// draftTokens simulates a player drafting tokens (simple AI)
func (g *Game) draftTokens(playerID int, availableTokens *[]Attribute, slots []int) {
	player := g.Players[playerID]

	for _, slotIdx := range slots {
		if len(*availableTokens) == 0 {
			break
		}

		// Simple AI: Analyze hand for best token to place
		attr, value := g.selectBestToken(player, *availableTokens, slotIdx)

		// Place the token
		g.Board.Slots[slotIdx] = &AttributeToken{
			Attribute: attr,
			Value:     value,
		}

		// Remove from available tokens
		for i, a := range *availableTokens {
			if a == attr {
				*availableTokens = append((*availableTokens)[:i], (*availableTokens)[i+1:]...)
				break
			}
		}

		if g.Verbose {
			fmt.Printf("Player %d places %s in Slot %d\n", playerID, g.Board.Slots[slotIdx], slotIdx+1)
		}
	}
}

// selectBestToken uses simple AI to select the best token to place
func (g *Game) selectBestToken(player *Player, availableTokens []Attribute, slotIdx int) (Attribute, bool) {
	// Count how many cards match each attribute value
	bestAttr := availableTokens[0]
	bestValue := false
	bestScore := -1

	for _, attr := range availableTokens {
		trueCount := 0
		falseCount := 0

		for _, card := range player.Hand {
			if card.Attributes[attr] {
				trueCount++
			} else {
				falseCount++
			}
		}

		// For Slot 1 (index 0), prefer attribute values we have more of
		// For other slots, prefer attribute values we have fewer of (to bury weakness)
		var score int
		var value bool

		if slotIdx == 0 { // Slot 1 - most important
			if trueCount > falseCount {
				score = trueCount
				value = true
			} else {
				score = falseCount
				value = false
			}
		} else { // Less important slots - bury weaknesses
			if trueCount < falseCount {
				score = 7 - falseCount // Prefer to bury our stronger side
				value = true           // Place the value we have less of
			} else {
				score = 7 - trueCount
				value = false
			}
		}

		if score > bestScore {
			bestScore = score
			bestAttr = attr
			bestValue = value
		}
	}

	return bestAttr, bestValue
}

// RunPlayPhase executes the reveal phase where each player plays a card
func (g *Game) RunPlayPhase() []struct {
	PlayerID int
	Card     Card
} {
	if g.Verbose {
		fmt.Printf("\n--- Trick %d: Reveal Phase ---\n", g.TrickNumber)
	}

	plays := make([]struct {
		PlayerID int
		Card     Card
	}, g.NumPlayers)

	// Each player selects a card
	for i := 0; i < g.NumPlayers; i++ {
		playerIdx := (g.CurrentLeader + i) % g.NumPlayers
		player := g.Players[playerIdx]

		// Simple AI: select best card based on board
		cardIdx := g.selectBestCard(player)
		card := player.Hand[cardIdx]

		// Remove card from hand
		player.Hand = append(player.Hand[:cardIdx], player.Hand[cardIdx+1:]...)

		plays[i] = struct {
			PlayerID int
			Card     Card
		}{playerIdx, card}

		if g.Verbose {
			fmt.Printf("Player %d plays: %s\n", playerIdx, card)
		}
	}

	return plays
}

// selectBestCard uses AI to select the best card to play
func (g *Game) selectBestCard(player *Player) int {
	// Simple strategy: Try to find a card that passes as many filters as possible
	bestIdx := 0
	bestScore := -1

	for i, card := range player.Hand {
		score := g.scoreCard(card)
		if score > bestScore {
			bestScore = score
			bestIdx = i
		}
	}

	return bestIdx
}

// scoreCard scores how well a card matches the protocol board
func (g *Game) scoreCard(card Card) int {
	score := 0
	for i := 0; i < 6; i++ {
		if g.Board.Slots[i] != nil {
			token := g.Board.Slots[i]
			if card.Matches(token.Attribute, token.Value) {
				// Use powers of 2 so earlier slots always dominate
				// Slot 1 (i=0) = 32, Slot 2 = 16, Slot 3 = 8, Slot 4 = 4, Slot 5 = 2, Slot 6 = 1
				score += 1 << (5 - i)
			}
		}
	}
	return score
}

// RunFilterPhase executes the judgement phase to determine the winner
func (g *Game) RunFilterPhase(plays []struct {
	PlayerID int
	Card     Card
}) int {
	if g.Verbose {
		fmt.Printf("\n--- Trick %d: Judgement Phase ---\n", g.TrickNumber)
	}

	// Start with all players active
	active := make([]int, len(plays))
	for i := range plays {
		active[i] = i
	}

	// Process each slot starting from Slot 1
	for slotIdx := 0; slotIdx < 6; slotIdx++ {
		if len(active) == 1 {
			break // Only one card remains
		}

		token := g.Board.Slots[slotIdx]
		if token == nil {
			continue
		}

		if g.Verbose {
			fmt.Printf("\nChecking Slot %d: %s\n", slotIdx+1, token)
		}

		// Check which cards match
		matching := []int{}
		for _, playIdx := range active {
			if plays[playIdx].Card.Matches(token.Attribute, token.Value) {
				matching = append(matching, playIdx)
			}
		}

		// If no cards match this slot, proceed to next slot (per rules)
		if len(matching) == 0 {
			if g.Verbose {
				fmt.Printf("  No cards match Slot %d. Proceeding to next slot.\n", slotIdx+1)
			}
			continue
		}

		// Eliminate non-matching cards
		if len(matching) > 0 {
			eliminated := []int{}
			for _, playIdx := range active {
				found := false
				for _, m := range matching {
					if m == playIdx {
						found = true
						break
					}
				}
				if !found {
					eliminated = append(eliminated, playIdx)
				}
			}

			if g.Verbose && len(eliminated) > 0 {
				for _, playIdx := range eliminated {
					fmt.Printf("  Player %d eliminated\n", plays[playIdx].PlayerID)
				}
			}

			active = matching
		}
	}

	// Determine winner
	var winnerIdx int
	if len(active) == 1 {
		winnerIdx = active[0]
	} else {
		// Tie-breaker: closest to leader in play order
		winnerIdx = active[0]
		if g.Verbose {
			fmt.Printf("\nTie-breaker: Multiple cards survived. Winner is closest to leader.\n")
		}
	}

	winnerPlayerID := plays[winnerIdx].PlayerID

	if g.Verbose {
		fmt.Printf("\nPlayer %d wins the trick!\n", winnerPlayerID)
	}

	// Award the trick
	g.Players[winnerPlayerID].TricksWon++
	for _, play := range plays {
		g.Players[winnerPlayerID].ScorePile = append(g.Players[winnerPlayerID].ScorePile, play.Card)
	}

	// Track statistics
	if g.Stats != nil {
		g.Stats.TricksByPlayer[winnerPlayerID]++

		// Track consecutive wins
		if g.LastTrickWinner == winnerPlayerID {
			g.ConsecutiveWins++
			if g.ConsecutiveWins == 2 {
				g.Stats.TwoStreaksByPlayer[winnerPlayerID]++
			}
			if g.ConsecutiveWins == 3 {
				g.Stats.ThreeStreaksByPlayer[winnerPlayerID]++
			}
		} else {
			g.ConsecutiveWins = 1
			g.LastTrickWinner = winnerPlayerID
		}
	}

	return winnerPlayerID
}

// RunActionPhase executes the action phase after each trick
func (g *Game) RunActionPhase() {
	if g.Verbose {
		fmt.Printf("\n--- Trick %d: Action Phase ---\n", g.TrickNumber)
		fmt.Println("Current Queen's Favor:")
		for i := 0; i < 6; i++ {
			if g.Board.Slots[i] != nil {
				fmt.Printf("  Slot %d: %s\n", i+1, g.Board.Slots[i])
			}
		}
		fmt.Println()
	}

	// Track the previous action to prevent immediate repeats
	var previousAction *Action

	// Each player takes one action, starting with the leader
	for i := 0; i < g.NumPlayers; i++ {
		playerIdx := (g.CurrentLeader + i) % g.NumPlayers
		player := g.Players[playerIdx]

		// AI decides which action to take (excluding previous action)
		action := g.selectBestAction(player, previousAction)

		switch action.Type {
		case "flip":
			slot := action.SlotIndex
			g.Board.Slots[slot].Value = !g.Board.Slots[slot].Value
			if g.Verbose {
				fmt.Printf("Player %d flips Slot %d to %s\n", playerIdx, slot+1, g.Board.Slots[slot])
			}

		case "swap":
			slot1, slot2 := action.SlotIndex, action.SlotIndex2
			g.Board.Slots[slot1], g.Board.Slots[slot2] = g.Board.Slots[slot2], g.Board.Slots[slot1]
			if g.Verbose {
				fmt.Printf("Player %d swaps Slot %d and Slot %d\n", playerIdx, slot1+1, slot2+1)
			}
		}

		// Save this action to prevent the next player from repeating it
		previousAction = &action
	}

	if g.Verbose {
		fmt.Println("\nQueen's Favor after actions:")
		for i := 0; i < 6; i++ {
			if g.Board.Slots[i] != nil {
				fmt.Printf("  Slot %d: %s\n", i+1, g.Board.Slots[i])
			}
		}
	}
}

// Action represents a player's action choice
type Action struct {
	Type       string // "flip" or "swap"
	SlotIndex  int    // For flip: the slot to flip. For swap: first slot
	SlotIndex2 int    // For swap: second slot
}

// actionsMatch returns true if two actions are the same
func actionsMatch(a1, a2 Action) bool {
	if a1.Type != a2.Type {
		return false
	}
	switch a1.Type {
	case "flip":
		return a1.SlotIndex == a2.SlotIndex
	case "swap":
		// Swaps match if they involve the same two slots (in either order)
		return (a1.SlotIndex == a2.SlotIndex && a1.SlotIndex2 == a2.SlotIndex2) ||
			(a1.SlotIndex == a2.SlotIndex2 && a1.SlotIndex2 == a2.SlotIndex)
	}
	return false
}

// selectBestAction uses AI to select the best action for a player
func (g *Game) selectBestAction(player *Player, previousAction *Action) Action {
	// First, find the best card with the CURRENT board (no action)
	currentBestScore := -1
	currentBestCard := -1
	for idx, card := range player.Hand {
		score := g.scoreCard(card)
		if score > currentBestScore {
			currentBestScore = score
			currentBestCard = idx
		}
	}

	// Evaluate all legal actions and pick the best
	bestAction := Action{Type: "flip", SlotIndex: 0}
	bestScore := -1
	bestCard := -1

	// Try all flip actions
	for slot := 0; slot < 6; slot++ {
		action := Action{Type: "flip", SlotIndex: slot}

		// Skip if this matches the previous action
		if previousAction != nil && actionsMatch(action, *previousAction) {
			continue
		}

		score, cardIdx := g.scoreActionOutcome(player, action)
		if score > bestScore {
			bestScore = score
			bestCard = cardIdx
			bestAction = action
		}
	}

	// Try all swap actions
	for slot1 := 0; slot1 < 6; slot1++ {
		for slot2 := slot1 + 1; slot2 < 6; slot2++ {
			action := Action{Type: "swap", SlotIndex: slot1, SlotIndex2: slot2}

			// Skip if this matches the previous action
			if previousAction != nil && actionsMatch(action, *previousAction) {
				continue
			}

			score, cardIdx := g.scoreActionOutcome(player, action)
			if score > bestScore {
				bestScore = score
				bestCard = cardIdx
				bestAction = action
			}
		}
	}

	// Log the decision
	if g.Verbose {
		actionType := bestAction.Type
		if actionType == "flip" {
			actionType = fmt.Sprintf("flip slot %d", bestAction.SlotIndex+1)
		} else if actionType == "swap" {
			actionType = fmt.Sprintf("swap slots %d+%d", bestAction.SlotIndex+1, bestAction.SlotIndex2+1)
		}
		cardChange := ""
		if bestCard != currentBestCard {
			cardChange = fmt.Sprintf(" (changes best card from #%d to #%d)", currentBestCard, bestCard)
		}

		blocked := ""
		if previousAction != nil {
			if previousAction.Type == "flip" {
				blocked = fmt.Sprintf(" [blocked: flip slot %d]", previousAction.SlotIndex+1)
			} else if previousAction.Type == "swap" {
				blocked = fmt.Sprintf(" [blocked: swap %d+%d]", previousAction.SlotIndex+1, previousAction.SlotIndex2+1)
			}
		}

		fmt.Printf("  (Player %d AI: best=%d, choosing %s%s%s)\n",
			player.ID, bestScore, actionType, cardChange, blocked)
	}

	return bestAction
}

// scoreActionOutcome simulates an action and scores how good it is for the player
func (g *Game) scoreActionOutcome(player *Player, action Action) (int, int) {
	// Save current board state
	savedBoard := make([]*AttributeToken, 6)
	for i := 0; i < 6; i++ {
		if g.Board.Slots[i] != nil {
			savedBoard[i] = &AttributeToken{
				Attribute: g.Board.Slots[i].Attribute,
				Value:     g.Board.Slots[i].Value,
			}
		}
	}

	// Apply the action temporarily
	switch action.Type {
	case "flip":
		g.Board.Slots[action.SlotIndex].Value = !g.Board.Slots[action.SlotIndex].Value

	case "swap":
		g.Board.Slots[action.SlotIndex], g.Board.Slots[action.SlotIndex2] = g.Board.Slots[action.SlotIndex2], g.Board.Slots[action.SlotIndex]
	}

	// Find the best card score with the new board
	bestScore := -1
	bestCardIdx := -1
	for idx, card := range player.Hand {
		score := g.scoreCard(card)
		if score > bestScore {
			bestScore = score
			bestCardIdx = idx
		}
	}

	// Restore board state
	for i := 0; i < 6; i++ {
		g.Board.Slots[i] = savedBoard[i]
	}

	return bestScore, bestCardIdx
}

// PlayTrick executes one complete trick.
// isLastTrick indicates this is round 7 (final round of a hand), in which case
// the Manipulate phase is skipped since tiles are about to be re-drafted.
// Returns the player ID of a sudden death winner (if the game ends after the
// Judge phase), or -1 if play continues.
func (g *Game) PlayTrick(isLastTrick bool) int {
	// Each trick: Present Phase -> Judge Phase -> (maybe) Manipulate Phase
	plays := g.RunPlayPhase()
	winner := g.RunFilterPhase(plays)

	// In sudden death, check for an outright leader after the Judge phase
	// (before Manipulation). If someone leads, the game ends immediately.
	if g.SuddenDeath {
		sdWinner := g.CheckWinner()
		if sdWinner != -1 {
			g.CurrentLeader = winner
			g.TrickNumber++
			return sdWinner
		}
	}

	// Update leader to round winner before manipulation (rules: manipulation
	// starts with the player that won this round, proceeding clockwise)
	g.CurrentLeader = winner

	// Skip Manipulate phase on the final round of a hand
	if !isLastTrick {
		g.RunActionPhase()
	} else if g.Verbose {
		fmt.Println("\n(Skipping Manipulate phase — final round of hand)")
	}

	g.TrickNumber++
	return -1
}

// PlayHand plays all tricks in a hand, returns true if game ended during hand
func (g *Game) PlayHand() bool {
	tricksPerHand := len(g.Players[0].Hand)
	for i := 0; i < tricksPerHand; i++ {
		isLastTrick := (i == tricksPerHand-1)
		sdWinner := g.PlayTrick(isLastTrick)

		// During sudden death, PlayTrick checks for a winner after the
		// Judge phase (before Manipulation) and returns the winner ID.
		if sdWinner != -1 {
			if g.Verbose {
				fmt.Printf("\n🎉 SUDDEN DEATH WINNER! Player %d breaks ahead with %d tricks! 🎉\n", sdWinner, g.Players[sdWinner].TricksWon)
				g.PrintScores()
			}
			return true // Game ended
		}
	}

	if g.Verbose {
		fmt.Printf("\n=== End of Hand %d ===\n", g.HandNumber)
		g.PrintScores()
	}

	g.HandNumber++
	return false // Hand completed, game continues
}

// CheckWinner returns the player ID if someone has won, otherwise -1
func (g *Game) CheckWinner() int {
	// Check if at least one player has 10+ tricks
	maxTricks := 0
	for _, player := range g.Players {
		if player.TricksWon > maxTricks {
			maxTricks = player.TricksWon
		}
	}

	// If no one has 10+ tricks yet, continue playing
	if maxTricks < 10 {
		return -1
	}

	// Find all players with the maximum number of tricks
	playersWithMax := []int{}
	for _, player := range g.Players {
		if player.TricksWon == maxTricks {
			playersWithMax = append(playersWithMax, player.ID)
		}
	}

	// If only one player has the max, they win
	if len(playersWithMax) == 1 {
		return playersWithMax[0]
	}

	// Multiple players tied at the max - sudden death continues
	return -1
}

// DetermineNextLeader sets the leader for the next hand
func (g *Game) DetermineNextLeader() {
	maxTricks := -1
	leaderID := 0

	for _, player := range g.Players {
		if player.TricksWon > maxTricks {
			maxTricks = player.TricksWon
			leaderID = player.ID
		}
	}

	g.CurrentLeader = leaderID
}

// PrintScores prints the current score
func (g *Game) PrintScores() {
	fmt.Println("\nCurrent Scores:")
	for _, player := range g.Players {
		fmt.Printf("  Player %d: %d tricks\n", player.ID, player.TricksWon)
	}
}

// Run executes the complete game simulation
func (g *Game) Run() int {
	if g.Verbose {
		fmt.Println("=== EYE OF THE BEE-HOLDER SIMULATION ===")
		fmt.Println("First player to 10 tricks wins!\n")
	}

	for {
		// Play a hand (7 tricks), may end early in sudden death
		gameEnded := g.PlayHand()
		if gameEnded {
			// Game ended in sudden death
			winner := -1
			maxTricks := 0
			for _, player := range g.Players {
				if player.TricksWon > maxTricks {
					maxTricks = player.TricksWon
					winner = player.ID
				}
			}
			if g.Stats != nil {
				g.Stats.WinsByPlayer[winner]++
			}
			return winner
		}

		// Check if anyone has 10+ tricks
		maxTricks := 0
		for _, player := range g.Players {
			if player.TricksWon > maxTricks {
				maxTricks = player.TricksWon
			}
		}

		// Check for winner after hand completes (not in sudden death)
		winner := g.CheckWinner()
		if winner != -1 {
			if g.Verbose {
				fmt.Printf("\n🎉 GAME OVER! Player %d wins with %d tricks! 🎉\n", winner, g.Players[winner].TricksWon)
				g.PrintScores()
			}
			if g.Stats != nil {
				g.Stats.WinsByPlayer[winner]++
			}
			return winner
		}

		// Check if we should enter sudden death mode
		if maxTricks >= 10 && !g.SuddenDeath {
			g.SuddenDeath = true
			if g.Verbose {
				fmt.Println("\n⚡ SUDDEN DEATH! Multiple players tied at the top. Playing until someone breaks ahead! ⚡")
				g.PrintScores()
			}
		}

		// Deal new hand
		g.DetermineNextLeader()
		g.DealNewHand()
	}
}

// RunStatistics runs multiple games and collects statistics
func RunStatistics(numPlayers int, numGames int) {
	stats := NewStats(numPlayers)

	fmt.Printf("Running %d games with %d players...\n", numGames, numPlayers)

	for i := 0; i < numGames; i++ {
		game := NewGame(numPlayers, false)
		game.Stats = stats
		game.Run()
		stats.GamesPlayed++

		if (i+1) % 100 == 0 {
			fmt.Printf("  Completed %d/%d games\n", i+1, numGames)
		}
	}

	// Print statistics
	fmt.Printf("\n=== STATISTICS FOR %d-PLAYER GAMES (%d games) ===\n", numPlayers, numGames)
	fmt.Println()

	// Win rates
	fmt.Println("Game Wins by Player:")
	for i := 0; i < numPlayers; i++ {
		winRate := float64(stats.WinsByPlayer[i]) / float64(numGames) * 100
		fmt.Printf("  Player %d: %d wins (%.1f%%)\n", i, stats.WinsByPlayer[i], winRate)
	}
	fmt.Println()

	// Trick totals
	totalTricks := 0
	for i := 0; i < numPlayers; i++ {
		totalTricks += stats.TricksByPlayer[i]
	}
	fmt.Println("Total Tricks Won by Player:")
	for i := 0; i < numPlayers; i++ {
		trickRate := float64(stats.TricksByPlayer[i]) / float64(totalTricks) * 100
		fmt.Printf("  Player %d: %d tricks (%.1f%%)\n", i, stats.TricksByPlayer[i], trickRate)
	}
	fmt.Println()

	// Streak analysis
	fmt.Println("2+ Trick Winning Streaks:")
	for i := 0; i < numPlayers; i++ {
		fmt.Printf("  Player %d: %d streaks\n", i, stats.TwoStreaksByPlayer[i])
	}
	fmt.Println()

	fmt.Println("3+ Trick Winning Streaks:")
	for i := 0; i < numPlayers; i++ {
		fmt.Printf("  Player %d: %d streaks\n", i, stats.ThreeStreaksByPlayer[i])
	}
	fmt.Println()
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Parse command line arguments
	if len(os.Args) > 1 && os.Args[1] == "stats" {
		// Statistics mode
		numGames := 1000
		if len(os.Args) > 2 {
			var err error
			numGames, err = strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Println("Usage: go run sim.go stats [num_games]")
				os.Exit(1)
			}
		}

		// Run statistics for all player counts
		for numPlayers := 2; numPlayers <= 5; numPlayers++ {
			RunStatistics(numPlayers, numGames)
		}
	} else {
		// Single game mode
		numPlayers := 4 // Default to 4 players
		if len(os.Args) > 1 {
			var err error
			numPlayers, err = strconv.Atoi(os.Args[1])
			if err != nil || numPlayers < 2 || numPlayers > 5 {
				fmt.Println("Usage: go run sim.go [num_players]")
				fmt.Println("       go run sim.go stats [num_games]")
				fmt.Println()
				fmt.Println("  num_players: 2-5 (default: 4)")
				fmt.Println("  stats: Run statistical analysis across all player counts")
				fmt.Println("  num_games: Number of games per player count (default: 1000)")
				os.Exit(1)
			}
		}

		fmt.Printf("Running simulation with %d players\n\n", numPlayers)

		// Create and run game with verbose output
		game := NewGame(numPlayers, true)
		game.Run()
	}
}
