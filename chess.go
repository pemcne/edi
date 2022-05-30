package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-joe/joe"
	"github.com/go-joe/joe/reactions"
	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
)

var Game *chess.Game
var Engine *uci.Engine

var chessEmoji map[string]string = map[string]string{
	"P": "white_pawn",
	"N": "white_knight",
	"B": "white_bishop",
	"R": "white_rook",
	"Q": "white_queen",
	"K": "white_king",
	"p": "black_pawn",
	"n": "black_knight",
	"b": "black_bishop",
	"r": "black_rook",
	"q": "black_queen",
	"k": "black_king",
}

var rankEmoji map[int]string = map[int]string{
	1: ":one:",
	2: ":two:",
	3: ":three:",
	4: ":four:",
	5: ":five:",
	6: ":six:",
	7: ":seven:",
	8: ":eight:",
}

var fileTranslate map[int]string = map[int]string{
	1: "a",
	2: "b",
	3: "c",
	4: "d",
	5: "e",
	6: "f",
	7: "g",
	8: "h",
}

var CHESSROOMS = []string{"C03GV6M95DM", "C03HC5JM28L"}

const chessBrainKey string = "chess.board"

func loadChessProgress() (string, error) {
	progress := ""
	_, err := Edi.Store.Get(chessBrainKey, &progress)
	if err != nil {
		return progress, err
	}
	return progress, nil
}

func storeChessState() error {
	state, err := Game.MarshalText()
	if err != nil {
		return err
	}
	return Edi.Store.Set(chessBrainKey, string(state))
}

func initChess() error {
	// Initialize the engine
	if Engine == nil {
		eng, err := uci.New("stockfish")
		if err != nil {
			return err
		}
		if err := eng.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
			return err
		}
		Engine = eng
	}

	// Load the existing progress if any
	progress, err := loadChessProgress()
	if err != nil {
		return err
	}
	var modules []func(*chess.Game)
	if len(progress) > 0 {
		fnProgress, err := chess.PGN(strings.NewReader(progress))
		if err != nil {
			return err
		}
		modules = append(modules, fnProgress)
	}
	game := chess.NewGame(modules...)
	Game = game
	return nil
}

func emojiChessBoard(fen []byte, lastMove *chess.Move) (string, error) {
	output := ""
	// Start board with files
	output += ":spacer::alphabet-white-a::alphabet-white-b::alphabet-white-c:"
	output += ":alphabet-white-d::alphabet-white-e::alphabet-white-f::alphabet-white-g:"
	output += ":alphabet-white-h:\n"
	rankStrs := strings.Split(string(fen), "/")

	preMove := ""
	postMove := ""
	if lastMove != nil {
		preMove = lastMove.S1().String()
		postMove = lastMove.S2().String()
	}
	count := 8
	for i, rankStr := range rankStrs {
		rank := 8 - i
		output += rankEmoji[count]
		file := 0
		// Figure out if it's a white or black square
		for _, char := range strings.Split(rankStr, "") {
			file++
			color := "black"
			if fmt.Sprintf("%s%d", fileTranslate[file], rank) == postMove {
				color = "active"
			} else if file%2 == (rank)%2 {
				color = "white"
			}
			if val, ok := chessEmoji[char]; ok {
				piece := fmt.Sprintf(":%s_%s:", val, color)
				output += piece
			} else {
				sep, err := strconv.Atoi(char)
				if err != nil {
					return output, err
				}
				for i := 0; i < sep; i++ {
					// Since file is increasing in here, also calc color
					color = "black_large"
					if fmt.Sprintf("%s%d", fileTranslate[file], rank) == preMove {
						color = "large_red"
					} else if file%2 == (rank)%2 {
						color = "white_large"
					}
					square := fmt.Sprintf(":%s_square:", color)
					output += square
					// Don't increment if we're on the last segment
					if i < sep-1 {
						file++
					}
				}
			}
		}
		output += "\n"
		count--
	}
	return output, nil
}

func printChessState(msg joe.Message) error {
	board, err := Game.Position().Board().MarshalText()
	if err != nil {
		return err
	}
	moves := Game.Moves()
	var lastmove *chess.Move
	if len(moves) > 0 {
		lastmove = moves[len(moves)-1]
	}
	txt, err := emojiChessBoard(board, lastmove)
	if err != nil {
		return err
	}
	output := ""
	if Game.Outcome() != chess.NoOutcome {
		Edi.Logger.Debug("Game complete")
		output += fmt.Sprintf("Game complete: %s\n", Game.Method())
	}
	if lastmove != nil {
		output += fmt.Sprintf("Last move: %s\n", lastmove.String())
		if Game.Outcome() != chess.NoOutcome && lastmove.HasTag(chess.Check) {
			output += "CHECK!\n"
		}
	}
	output += txt
	msg.Respond(output)
	return nil
}

func ChessNew(msg joe.Message) error {
	if !correctRoom(msg, CHESSROOMS) {
		return nil
	}
	_, err := Edi.Store.Delete(chessBrainKey)
	if err != nil {
		return err
	}
	initChess()
	return printChessState(msg)
}

func ChessMove(msg joe.Message) error {
	if !correctRoom(msg, CHESSROOMS) {
		return nil
	}
	// Run a move
	move := msg.Matches[0]
	if err := Game.MoveStr(move); err != nil {
		msg.React(reactions.Reaction{
			Shortcode: "x",
		})
		return err
	}
	if Game.Outcome() != chess.NoOutcome {
		printChessState(msg)
		return ChessNew(msg)
	}

	// Still going so CPU move
	cmdPos := uci.CmdPosition{Position: Game.Position()}
	cmdGo := uci.CmdGo{MoveTime: time.Second}
	if err := Engine.Run(cmdPos, cmdGo); err != nil {
		return err
	}
	cpuMove := Engine.SearchResults().BestMove
	if err := Game.Move(cpuMove); err != nil {
		return err
	}
	if Game.Outcome() != chess.NoOutcome {
		printChessState(msg)
		return ChessNew(msg)
	}
	err := storeChessState()
	if err != nil {
		return err
	}
	return printChessState(msg)

}

func ChessState(msg joe.Message) error {
	if !correctRoom(msg, CHESSROOMS) {
		return nil
	}
	return printChessState(msg)
}
