package challenge_generator

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// HashcashChallenge holds the challenge data
type HashcashChallenge struct {
	ID            string `json:"id"`
	Prefix        int    `json:"prefix"` // Number of leading zeros required
	RandomEntropy string `json:"random"`
	Version       int    `json:"version"`
	Time          int64  `json:"time"`
	Proof         int    `json:"proof"`
}

func (h HashcashChallenge) String() string {
	return fmt.Sprintf("%s:%d:%s:%d:%d", h.ID, h.Prefix, h.RandomEntropy, h.Version, h.Time)
}

type HashcashChallengeGenerator struct {
	prefix  int
	version int
}

func NewHashcashChallengeGenerator(prefix, version int) *HashcashChallengeGenerator {
	return &HashcashChallengeGenerator{
		prefix:  prefix,
		version: version,
	}
}

func (h *HashcashChallengeGenerator) GenerateChallenge(id string) ([]byte, error) {
	randomUUID := uuid.New().String()
	base64Entrophy := base64.StdEncoding.EncodeToString([]byte(randomUUID))
	timeNow := time.Now().Unix()
	challenge := HashcashChallenge{
		ID:            id,
		Prefix:        h.prefix,
		RandomEntropy: base64Entrophy,
		Version:       h.version,
		Time:          timeNow,
	}
	marshaled, err := json.Marshal(challenge)
	if err != nil {
		return nil, err
	}

	return marshaled, nil
}

// verifies if the proof is correct for the given challenge
func (h *HashcashChallengeGenerator) ValidateChallengeResponse(challengeBytes, proofBytes []byte) (bool, error) {
	challenge := HashcashChallenge{}
	err := json.Unmarshal(challengeBytes, &challenge)
	if err != nil {
		return false, err
	}

	proofChallenge := HashcashChallenge{}
	err = json.Unmarshal(proofBytes, &proofChallenge)
	if err != nil {
		return false, err
	}

	proof := strconv.Itoa(proofChallenge.Proof)
	hash := sha256.New()
	hash.Write([]byte(challenge.String() + proof))
	hashValue := hash.Sum(nil)
	hashBigInt := new(big.Int).SetBytes(hashValue)

	// Create a big integer representing the prefix of leading zeros
	requiredPrefix := new(big.Int).Lsh(big.NewInt(1), uint(256-challenge.Prefix))

	return hashBigInt.Cmp(requiredPrefix) == -1, nil
}

func (h *HashcashChallengeGenerator) DeserializeChallengeID(challengeBytes []byte) (string, error) {
	challenge := HashcashChallenge{}
	err := json.Unmarshal(challengeBytes, &challenge)
	if err != nil {
		return "", err
	}

	return challenge.ID, nil
}

// try to find a valid proof for the challenge
func (h *HashcashChallengeGenerator) SolveChallenge(challengeBytes []byte) ([]byte, error) {
	challenge := HashcashChallenge{}
	err := json.Unmarshal(challengeBytes, &challenge)
	if err != nil {
		return nil, err
	}

	// Attempt different proof values until one is found that satisfies the challenge
	for i := 0; ; i++ {
		proof := strconv.Itoa(i)
		hash := sha256.New()
		hash.Write([]byte(challenge.String() + proof))
		hashValue := hash.Sum(nil)
		hashBigInt := new(big.Int).SetBytes(hashValue)

		// Create a big integer representing the prefix of leading zeros
		requiredPrefix := new(big.Int).Lsh(big.NewInt(1), uint(256-challenge.Prefix))

		if hashBigInt.Cmp(requiredPrefix) == -1 {
			challenge.Proof = i
			proof, err := json.Marshal(challenge)
			if err != nil {
				return nil, err
			}
			return proof, nil
		}
	}
}
