package challenge_generator

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateChallenge(t *testing.T) {
	generator := NewHashcashChallengeGenerator(4, 1)
	id := "test"
	challenge, err := generator.GenerateChallenge(id)
	require.NoError(t, err, "GenerateChallenge")
	require.NotNil(t, challenge, "GenerateChallenge")

	var hashcashChallenge HashcashChallenge
	err = json.Unmarshal(challenge, &hashcashChallenge)
	require.NoError(t, err, "Unmarshal")
}

func TestValidateChallengeResponse(t *testing.T) {
	generator := NewHashcashChallengeGenerator(4, 1)
	id := "test"
	challengeBytes, err := generator.GenerateChallenge(id)
	require.NoError(t, err, "GenerateChallenge")
	require.NotNil(t, challengeBytes, "GenerateChallenge")

	proof, err := generator.SolveChallenge(challengeBytes)
	require.NoError(t, err, "SolveChallenge")
	ok, err := generator.ValidateChallengeResponse(challengeBytes, proof)
	require.NoError(t, err, "ValidateChallengeResponse correct proof")
	require.True(t, ok, "ValidateChallengeResponse correct proof")

	ok, err = generator.ValidateChallengeResponse(challengeBytes, challengeBytes)
	require.NoError(t, err, "ValidateChallengeResponse incorrect proof")
	require.False(t, ok, "ValidateChallengeResponse incorrect proof")
}
