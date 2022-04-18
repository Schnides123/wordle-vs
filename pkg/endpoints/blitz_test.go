package endpoints_test

import (
	"strings"
	"testing"
	"time"

	"github.com/Schnides123/wordle-vs/pkg/wordle"

	"github.com/Schnides123/wordle-vs/pkg/testutil"
	"github.com/stretchr/testify/require"
)

// Tests that when two participants are in a game, and one attempts a guess,
// they are both notified
func TestBroadcastGuess(t *testing.T) {
	defer testutil.StopTestServer(testutil.StartTestServer())

	user1 := testutil.NewUser("abcde")
	user2 := testutil.NewUser("abcdf")

	err := user1.ConnectToGame("12345", &wordle.DefaultOptions)
	require.NoError(t, err)

	err = user2.ConnectToGame("12345", nil)
	require.NoError(t, err)

	err = user1.WaitForState(1 * time.Second)
	require.NoError(t, err, "waiting for state")
	require.Equal(t, "12345", user1.State().Id)

	err = user1.WaitForState(2 * time.Second)
	require.NoError(t, err, "user1 waiting for state after guess")

	err = user2.WaitForState(1 * time.Second)
	require.NoError(t, err, "waiting for state")

	guessWord := strings.Repeat("a", user1.State().Word.Length)
	err = user1.Guess(guessWord)
	require.NoError(t, err, "user1 makes guess")

	err = user1.WaitForState(2 * time.Second)
	require.NoError(t, err, "user1 waiting for state after guess")

	err = user2.WaitForState(1 * time.Second)
	require.NoError(t, err, "user2 waiting for state after guess")

	require.EqualValues(t, user1.State(), user2.State(),
		"user1 and user2 expect identical state after user1 guesses")

	require.Len(t, user1.State().Guessed, 1)
	require.Equal(t, user1.State().Guessed[0].Word, guessWord)
	require.Equal(t, user1.State().Guessed[0].Player, user1.GetPlayerID())

	err = user1.Disconnect()
	require.NoError(t, err)

	err = user2.Disconnect()
	require.NoError(t, err)
}

//Tests that two players can play a complete game together
func TestGame(t *testing.T) {
	defer testutil.StopTestServer(testutil.StartTestServer())

	user1 := testutil.NewUser("abcde")
	user2 := testutil.NewUser("bcdef")

	options := wordle.GameOptions{
		TurnLength: 30 * time.Second,
		WordLength: 0,
		MaxGuesses: -1,
		Word:       "hello",
	}

	err := user1.ConnectToGame("12346", &options)
	require.NoError(t, err)

	err = user2.ConnectToGame("12346", nil)
	require.NoError(t, err)

	err = user1.WaitForState(1 * time.Second)
	require.NoError(t, err, "waiting for state")

	err = user2.WaitForState(1 * time.Second)
	require.NoError(t, err, "waiting for state")

	err = user1.WaitForState(1 * time.Second)
	require.NoError(t, err, "waiting for state")

	require.EqualValues(t, user1.State(), user2.State())
	require.Len(t, user1.State().Guessed, 0)
	require.Equal(t, user1.State().Id, "12346")

	err = user1.Guess("hella")
	require.NoError(t, err, "player 1 guesses")

	err = user1.WaitForState(2 * time.Second)
	require.NoError(t, err, "wait u1 after player 1 guesses")

	err = user2.WaitForState(2 * time.Second)
	require.NoError(t, err, "wait u2 after player 1 guesses")

	require.Len(t, user1.State().Guessed, 1)
	require.Equal(t, user1.State().Guessed[0].Word, "hella")
	require.Equal(t, user1.State().Guessed[0].Player, "abcde")
	require.Equal(t, user1.State().Guessed[0].Result, []int{2, 2, 2, 2, 0})
	require.Equal(t, user2.State(), user1.State())

	err = user2.Guess("hello")
	require.NoError(t, err)
	err = user1.WaitForState(1 * time.Second)
	require.NoError(t, err)
	err = user2.WaitForState(1 * time.Second)
	require.NoError(t, err, "player 2 guesses, wins")
	require.Len(t, user1.State().Guessed, 2)
	require.Equal(t, user1.State().Guessed[1].Word, "hello")
	require.Equal(t, user1.State().Guessed[1].Player, "bcdef")
	require.Equal(t, user1.State().Guessed[1].Result, []int{2, 2, 2, 2, 2})
	require.True(t, user1.State().Winner)
	require.Equal(t, user2.State(), user1.State())

	err = user1.Guess("hella")
	require.NoError(t, err)
	err = user1.WaitForState(1 * time.Second)
	require.Error(t, err)

	err = user1.Disconnect()
	require.NoError(t, err)
	err = user2.Disconnect()
	require.NoError(t, err)

}
