package main

import (
	"testing"
	"time"

	"github.com/Schnides123/wordle-vs/pkg/wordle"

	"github.com/Schnides123/wordle-vs/pkg/testutil"
	"github.com/stretchr/testify/require"
)

func TestSinglePeerConnect(t *testing.T) {
	defer testutil.StopTestServer(testutil.StartTestServer())

	user := testutil.NewUser("abcde")
	err := user.ConnectToGame("12345", &wordle.DefaultOptions)
	require.NoError(t, err)

	err = user.Disconnect()
	require.NoError(t, err)
}

func TestMultiplePeersConnect(t *testing.T) {
	defer testutil.StopTestServer(testutil.StartTestServer())

	user1 := testutil.NewUser("abcde")
	user2 := testutil.NewUser("abcde")

	err := user1.ConnectToGame("12345", &wordle.DefaultOptions)
	require.NoError(t, err)

	err = user2.ConnectToGame("12345", nil)
	require.NoError(t, err)

	err = user1.WaitForState(1 * time.Second)
	require.NoError(t, err, "waiting for state")

	err = user2.WaitForState(1 * time.Second)
	require.NoError(t, err, "waiting for state")

	require.EqualValues(t, user1.State(), user2.State(), "should be in same game")

	err = user1.Disconnect()
	require.NoError(t, err)

	err = user2.Disconnect()
	require.NoError(t, err)
}
