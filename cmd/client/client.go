package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Schnides123/wordle-vs/pkg/testutil"
	"github.com/Schnides123/wordle-vs/pkg/util"
	"github.com/Schnides123/wordle-vs/pkg/wordle"
)

/*

client new
client join <gameID>
client join queue

TURN
honery badger
-> fat couch

R ? ? S ?

BOARD
R A I S E
X X X X X
X X X X X
X X X X X
X X X X X

> sound

TURN
-> honery badger
fat couch

R O ? S ?

BOARD
R A I S E
S O U N D
X X X X X
X X X X X
X X X X X

Waiting for turn...
Attempt Steal?
>

*/

func main() {
	ur, _ := url.Parse("ws://localhost:42069/")
	testutil.SetTestURL(ur)
	parseArgs()
}

func parseArgs() {
	if len(os.Args) < 2 {
		fmt.Println("bad args yo")
		return
	}
	var err error

	newCmd := flag.NewFlagSet("new", flag.ExitOnError)
	joinCmd := flag.NewFlagSet("join", flag.ExitOnError)

	user := testutil.NewUser("")
	defer user.Disconnect()

	switch os.Args[1] {
	case "new":
		newCmd.Parse(os.Args[2:])
		opts := wordle.DefaultOptions
		opts.Word = "hello"
		err = user.NewGame(&opts)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		err = user.WaitForState(3 * time.Second)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Println("new game:", user.State().Id)
		fmt.Println("waiting for other players...")
		waitForTurn(user)
		repl(user)
	case "join":
		joinCmd.Parse(os.Args[2:])
		joinGameID := "queue"
		if joinCmd.NArg() > 1 {
			os.Exit(1)
		} else if joinCmd.NArg() == 1 {
			joinGameID = joinCmd.Arg(0)
		}

		fmt.Println("join", joinGameID)

		if joinGameID == "queue" {
			err = user.JoinMatchmaking()
		} else {
			err = user.ConnectToGame(joinGameID, nil)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		waitForTurn(user)
		repl(user)
	default:
		fmt.Println("bad args yo")
		os.Exit(1)
	}
}

func waitForTurn(u *testutil.User) error {
	for u.State() == nil ||
		!u.CanPlay() ||
		u.State().Players == nil ||
		u.State().Players.Length() < 2 {
		err := u.WaitForState(3 * time.Second)
		if err != nil {
			if err.Error() == "timeout" {
				continue
			}
			return err
		}
	}
	return nil
}

func drainState(u *testutil.User) error {
	for {
		err := u.WaitForState(1 * time.Second)
		if err != nil {
			if err.Error() == "timeout" {
				return nil
			}
			return err
		}
	}
}

func status(g *wordle.Game) string {
	progress := []rune{}
	for i := 0; i < g.Word.Length; i++ {
		progress = append(progress, '?')
	}

	for _, guess := range g.Guessed {
		fmt.Println(guess)
		runes := []rune(guess.Word)

		for i, hint := range guess.Result {
			if hint == 2 {
				progress[i] = runes[i]
			}
		}
	}

	return string(progress)
}

func repl(u *testutil.User) {
	for {
		// Print TURN
		fmt.Println("TURN")
		for _, p := range u.State().Players.Values() {
			if len(u.State().Guessed) == 0 ||
				u.State().Guessed[len(u.State().Guessed)-1].Player != p {
				fmt.Println("\t", "->", p)
			} else {
				fmt.Println("\t", p)
			}
		}

		// Print Game status
		fmt.Println("WORD")
		fmt.Println(strings.ToLower(status(u.State())))

		// Print BOARD
		fmt.Println("GUESSES")
		if len(u.State().Guessed) == 0 {
			fmt.Println("no moves yet")
		} else {
			for _, g := range u.State().Guessed {
				fmt.Println("\t", g.Word)
			}
		}

		// If not turn, wait
		if err := waitForTurn(u); err != nil {
			fmt.Println(err)
			return
		}

		// Prompt turn
		command := util.ScanLine(context.TODO())

		if err := u.Guess(command); err != nil {
			fmt.Println(err)
			return
		}

		err := drainState(u)
		for err != nil {
			fmt.Println(err)
			err = drainState(u)
		}
	}
}
