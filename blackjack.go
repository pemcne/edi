package main

import (
	"fmt"

	"github.com/go-joe/joe"
	"github.com/pemcne/carddeck/v2"
)

const blackjackStoreKey string = "blackjack"

type BlackjackHand struct {
	Cards []carddeck.Card `json:"hand"`
	Total int             `json:"total"`
	Bust  bool            `json:"bust"`
	Stand bool            `json:"stand"`
}

type BlackjackState struct {
	Bet        int             `json:"bet"`
	Bank       int             `json:"bank"`
	Games      int             `json:"games"`
	Wins       int             `json:"wins"`
	DealerHand BlackjackHand   `json:"dealer"`
	PlayerHand []BlackjackHand `json:"player"`
}

var blackjackTranslate map[int]string = map[int]string{
	1:  ":alphabet-white-a:",
	2:  ":two:",
	3:  ":three:",
	4:  ":four:",
	5:  ":five:",
	6:  ":six:",
	7:  ":seven:",
	8:  ":eight:",
	9:  ":nine:",
	10: ":keycap_ten:",
	11: ":alphabet-white-j:",
	12: ":alphabet-white-q:",
	13: ":alphabet-white-k:",
}

var deck carddeck.Deck
var defaultBet = 50

func loadBlackjack() {
	deckConfig := carddeck.DeckConfig{Size: 4}
	deck.Initialize(deckConfig)
	deck.Shuffle()
}

func getBlackjackState() (BlackjackState, error) {
	state := BlackjackState{}
	ok, err := Edi.Store.Get(blackjackStoreKey, &state)
	if err != nil {
		return state, err
	}
	if !ok {
		// Some sane initial values
		state.Bet = defaultBet
	}
	return state, nil
}

func setBlackjackState(state BlackjackState) error {
	return Edi.Store.Set(blackjackStoreKey, state)
}

func handTotal(hand []carddeck.Card) int {
	var total int
	var aces = 0
	for _, c := range hand {
		cval := c.Value.Value
		if cval >= 2 && cval <= 10 {
			// Normal number
			total += cval
		} else if cval > 10 {
			// Face cards are all 10
			total += 10
		} else {
			// Ace, store and process later
			aces++
		}
	}
	// Now we've done everything but aces
	for range aces {
		if total+11 <= 21 {
			// See if we can process an 11
			total += 11
		} else {
			// Has to be a 1 instead
			total += 1
		}
	}
	return total
}

func runBlackjackDealer(state *BlackjackState) error {
	hand := &state.DealerHand
	for {
		if hand.Total >= 17 {
			if hand.Total > 21 {
				hand.Bust = true
			} else {
				hand.Stand = true
			}
			break
		}
		cards, err := deck.Draw(1)
		if err != nil {
			return err
		}
		hand.Cards = append(hand.Cards, cards[0])
		hand.Total = handTotal(hand.Cards)
	}
	return nil
}

func printBlackjack(state *BlackjackState, done bool) string {
	output := "Dealer hand: "
	for i, c := range state.DealerHand.Cards {
		if i > 0 && !done {
			output += ":alphabet-white-question:"
		} else {
			output += blackjackTranslate[c.Value.Value]
		}
	}
	if done {
		output += fmt.Sprintf(" | Total: %d", state.DealerHand.Total)
		if state.DealerHand.Bust {
			output += " :x: "
		}
	}
	for ihand, phand := range state.PlayerHand {
		if len(state.PlayerHand) > 1 {
			output += fmt.Sprintf("\n Player hand %d: ", ihand+1)
		} else {
			output += "\nPlayer hand: "
		}
		for _, card := range phand.Cards {
			output += blackjackTranslate[card.Value.Value]
		}
		total := handTotal(phand.Cards)
		output += fmt.Sprintf(" | Total: %d", total)
		if done {
			if phand.Bust {
				output += " :x: "
			} else if state.DealerHand.Bust {
				output += " :tada: "
			} else if state.DealerHand.Total < phand.Total {
				output += " :tada: "
			}
		}
	}
	return output
}

func BlackjackGame(msg joe.Message) error {
	// Reset the state
	state, err := getBlackjackState()
	if err != nil {
		return err
	}
	state.Games++
	Edi.Logger.Debug(fmt.Sprintf("Blackjack state: %v", state))
	dealerHand, err := deck.Draw(2)
	if err != nil {
		return err
	}
	state.DealerHand = BlackjackHand{
		Bust:  false,
		Stand: false,
		Cards: dealerHand,
		Total: handTotal(dealerHand),
	}
	Edi.Logger.Debug(fmt.Sprintf("Blackjack dealer hand: %q", dealerHand))
	playerHand, err := deck.Draw(2)
	if err != nil {
		return err
	}
	state.PlayerHand = []BlackjackHand{{
		Bust:  false,
		Stand: false,
		Cards: playerHand,
		Total: handTotal(playerHand),
	}}
	Edi.Logger.Debug(fmt.Sprintf("Blackjack player hand: %q", playerHand))
	msg.Respond(printBlackjack(&state, false))
	return setBlackjackState(state)
}

func BlackjackHit(msg joe.Message) error {
	state, err := getBlackjackState()
	if err != nil {
		return err
	}
	finished := true
	for i := range state.PlayerHand {
		phand := &state.PlayerHand[i]
		if phand.Bust || phand.Stand {
			continue
		} else {
			finished = false
			draws, err := deck.Draw(1)
			if err != nil {
				return err
			}
			phand.Cards = append(phand.Cards, draws[0])
			phand.Total = handTotal(phand.Cards)
			if phand.Total > 21 {
				finished = true
				phand.Bust = true
			}
			// We only process one hand at a time
			break
		}
	}
	msg.Respond(printBlackjack(&state, finished))
	if finished {
		return BlackjackGame(msg)
	}
	return setBlackjackState(state)
}

func BlackjackStand(msg joe.Message) error {
	state, err := getBlackjackState()
	if err != nil {
		return err
	}
	processed := 0
	for i := range state.PlayerHand {
		phand := &state.PlayerHand[i]
		processed++
		if phand.Bust || phand.Stand {
			continue
		}
		phand.Stand = true
		break
	}
	finished := processed == len(state.PlayerHand)
	if finished {
		err = runBlackjackDealer(&state)
		if err != nil {
			return err
		}
	}
	msg.Respond(printBlackjack(&state, finished))
	return setBlackjackState(state)
}
