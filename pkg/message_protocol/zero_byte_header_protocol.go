package message_protocol

import "fmt"

const (
	magicByte = 0x7F
)

var (
	ErrInvalidMagicByte = fmt.Errorf("invalid message format, missing magic byte")
)

type ZeroByteHeaderProtocol struct{}

func NewZeroByteHeaderProtocol() *ZeroByteHeaderProtocol {
	return new(ZeroByteHeaderProtocol)
}

func (r *ZeroByteHeaderProtocol) Parse(data []byte) (int, []byte, error) {
	// magic byte + message type + message)
	if len(data) == 0 || data[0] != magicByte {
		return 0, nil, ErrInvalidMagicByte
	}
	return int(data[1]), data[2 : len(data)-1], nil
}

func (r *ZeroByteHeaderProtocol) Encode(messageType int, message []byte) []byte {
	message = append(message, '\n')
	return append([]byte{magicByte, byte(messageType)}, message...)
}
