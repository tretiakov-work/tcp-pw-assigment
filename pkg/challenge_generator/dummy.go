package challenge_generator

type Dummy struct{}

func NewDummy() *Dummy {
	return new(Dummy)
}

const dummyID = "dummy"

func (d *Dummy) GenerateChallenge(_ string) ([]byte, error) {
	return []byte("dummy challenge"), nil
}

func (d *Dummy) ValidateChallengeResponse(_, proof []byte) (bool, error) {
	return string(proof) == "dummy challenge\n", nil
}

func (d *Dummy) DeserializeChallengeID(_ []byte) (string, error) {
	return dummyID, nil
}
