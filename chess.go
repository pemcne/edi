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
	"P": ":white_pawn:",
	"N": ":white_knight:",
	"B": ":white_bishop:",
	"R": ":white_rook:",
	"Q": ":white_queen:",
	"K": ":white_king:",
	"p": ":black_pawn:",
	"n": ":black_knight:",
	"b": ":black_bishop:",
	"r": ":black_rook:",
	"q": ":black_queen:",
	"k": ":black_king:",
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
	eng, err := uci.New("stockfish")
	if err != nil {
		return err
	}
	if err := eng.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
		return err
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
	Engine = eng
	return nil
}

func emojiChessBoard(fen []byte) (string, error) {
	output := ""
	// Start board with columns
	output += ":spacer::alphabet-white-a::alphabet-white-b::alphabet-white-c:"
	output += ":alphabet-white-d::alphabet-white-e::alphabet-white-f::alphabet-white-g:"
	output += ":alphabet-white-h:\n"
	rankStrs := strings.Split(string(fen), "/")
	count := 8
	for rank, rankStr := range rankStrs {
		output += rankEmoji[count]
		file := 0
		for _, char := range strings.Split(rankStr, "") {
			file++
			if val, ok := chessEmoji[char]; ok {
				output += val
			} else {
				sep, err := strconv.Atoi(char)
				if err != nil {
					return output, err
				}
				for i := 0; i < sep; i++ {
					square := ":black_medium_square:"
					if file%2 == (rank+1)%2 {
						square = ":white_medium_square:"
					}
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

func printChessState(msg joe.Message, cpu *chess.Move) error {
	position, err := Game.Position().MarshalText()
	if err != nil {
		return err
	}
	fmt.Printf("Position: %s\n", string(position))
	Edi.Logger.Debug("Getting game board text")
	board, err := Game.Position().Board().MarshalText()
	if err != nil {
		return err
	}
	fmt.Printf("Board: %s\n", string(board))
	Edi.Logger.Debug("Getting emoji board")
	txt, err := emojiChessBoard(board)
	if err != nil {
		return err
	}
	output := ""
	Edi.Logger.Debug("Setting cpu move if any")
	if cpu != nil {
		output = fmt.Sprintf("CPU move: %s\n", cpu.String())
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
	return printChessState(msg, nil)
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
		out := fmt.Sprintf("Game complete %s by %s", Game.Outcome(), Game.Method())
		msg.Respond(out)
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
		out := fmt.Sprintf("Game complete %s by %s", Game.Outcome(), Game.Method())
		msg.Respond(out)
	} else {
		err := storeChessState()
		if err != nil {
			return err
		}
		return printChessState(msg, cpuMove)
	}
	return ChessNew(msg)
}

func ChessState(msg joe.Message) error {
	if !correctRoom(msg, CHESSROOMS) {
		return nil
	}
	return printChessState(msg, nil)
}
