package model

type TokenType uint8

const (
	PacketFee = iota + 1
	PacketValue
)

func (tt TokenType) String() string {
	switch tt {
	case PacketFee:
		return "packet_fee"
	case PacketValue:
		return "packet_value"
	default:
		return "unknown"
	}
}
