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

func initChess() error {
	eng, err := uci.New("stockfish")
	if err != nil {
		return err
	}
	if err := eng.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
		return err
	}
	game := chess.NewGame()
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
	for _, rankStr := range rankStrs {
		output += rankEmoji[count]
		for _, char := range strings.Split(rankStr, "") {
			if val, ok := chessEmoji[char]; ok {
				output += val
			} else {
				sep, err := strconv.Atoi(char)
				if err != nil {
					return output, err
				}
				for i := 0; i < sep; i++ {
					output += ":black_small_square:"
				}
			}
		}
		output += "\n"
		count--
	}
	return output, nil
}

func printChessState(msg joe.Message, cpu *chess.Move) error {
	board, err := Game.Position().Board().MarshalText()
	if err != nil {
		return err
	}
	txt, err := emojiChessBoard(board)
	if err != nil {
		return err
	}
	output := ""
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
		return printChessState(msg, cpuMove)
	}
	return nil
}

func ChessState(msg joe.Message) error {
	if !correctRoom(msg, CHESSROOMS) {
		return nil
	}
	return printChessState(msg, nil)
}
