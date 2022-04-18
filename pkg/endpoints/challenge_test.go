package endpoints_test

import (
	"testing"
	"time"

	"github.com/Schnides123/wordle-vs/pkg/testutil"
	"github.com/Schnides123/wordle-vs/pkg/wordle"
	"github.com/stretchr/testify/require"
)

func TestBasicChallenge(t *testing.T) {
	defer testutil.StopTestServer(testutil.StartTestServer())
	user1 := testutil.NewUser("")
	user2 := testutil.NewUser("")

	err := user1.NewGame(&wordle.DefaultOptions)
	require.NoError(t, err)

	err = user1.WaitForState(1 * time.Second)
	require.NoError(t, err)

	err = user2.ConnectToGame(user1.State().Id, nil)
	require.NoError(t, err)

	err = user1.WaitForState(1 * time.Second)
	require.NoError(t, err)

	err = user2.WaitForState(1 * time.Second)
	require.NoError(t, err)

	require.EqualValues(t, user1.State(), user2.State())

	t.Log(user1.GetPlayerID())
	t.Log(user2.GetPlayerID())
}
