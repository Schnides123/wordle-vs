package endpoints_test

import (
	"testing"
	"time"

	"github.com/Schnides123/wordle-vs/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestMatchmaking(t *testing.T) {
	defer testutil.StopTestServer(testutil.StartTestServer())

	u1 := testutil.NewUser("")
	u2 := testutil.NewUser("")

	u1.JoinMatchmaking()
	u2.JoinMatchmaking()

	err := u1.WaitForState(1 * time.Second)
	assert.NoError(t, err)
	assert.True(t, len(u1.GetPlayerID()) > 0)
	assert.True(t, len(u1.State().Id) > 0)
	err = u2.WaitForState(1 * time.Second)
	assert.NoError(t, err)
	assert.True(t, len(u2.GetPlayerID()) > 0)
	assert.True(t, len(u2.State().Id) > 0)
	t.Log(u1.GetPlayerID())
	t.Log(u2.GetPlayerID())

	assert.Equal(t, u1.State().Id, u2.State().Id)
}

func TestMatchmakingCleanup(t *testing.T) {
	// Ensures that if u1 joins and leaves matchmaking, that u2 does not
	// match with u1 upon joining.
	defer testutil.StopTestServer(testutil.StartTestServer())

	u1 := testutil.NewUser("")
	u2 := testutil.NewUser("")
	u3 := testutil.NewUser("")

	u1.JoinMatchmaking()
	u1.Disconnect()

	u2.JoinMatchmaking()
	u3.JoinMatchmaking()

	err := u2.WaitForState(1 * time.Second)
	assert.NoError(t, err)

	err = u3.WaitForState(1 * time.Second)
	assert.NoError(t, err)

	assert.Equal(t, u2.State().Id, u3.State().Id)
}
