package customtype

const (
	InvalidNum NumberType = "invalid"
	OddNum     NumberType = "odd"
	EvenNum    NumberType = "even"
)

// NumberType represents any number type. e.g. Odd, Even, etc.
type NumberType string

func (t NumberType) String() string {
	return string(t)
}

// NewNumType returns the NumberType associate to the given string number type.
// Returns InvalidNum if the number type is not supported.
func NewNumType(nType string) NumberType {
	switch nType {
	case OddNum.String():
		return OddNum
	case EvenNum.String():
		return EvenNum
	default:
		return InvalidNum
	}
}
