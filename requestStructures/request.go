package requestStructures

type APPInfo struct {
	Name string
}

type OperationSymbol int

const (
	Plus OperationSymbol = iota + 1
	Minus
	Multiply
)

func (u OperationSymbol) Validate() bool {
	switch u {
	case Plus:
		return true
	case Minus:
		return true
	case Multiply:
		return true
	default:
		return false
	}
}

type CalculateRequest struct {
	Operand1  string `json:"operand1" validate:"required,number"`
	Operand2  string `json:"operand2" validate:"required,number"`
	Email     string `json:"email" validate:"required,email"`
	Operation int    `json:"operand" validate:"required,oneof=1 2 3"`
}
