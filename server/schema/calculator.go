package schema

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"math"
	"strconv"
)

type OperationSymbol int

const (
	Plus OperationSymbol = iota + 1
	Minus
	Multiply
	Division
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
	Operand1  string          `json:"operand1" validate:"required,number"`
	Operand2  string          `json:"operand2" validate:"required,number"`
	Email     string          `json:"email" validate:"required,email"`
	Operation OperationSymbol `json:"operand" validate:"required,oneof=1 2 3 4"`
}

func CalculatorFunction(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Log("SchemaGet request received")

	var req CalculateRequest

	err := wscutils.BindJSON(c, &req)
	if err != nil {
		lh.Debug0().LogActivity("error while binding json request error:", err.Error())
		return
	}

	lh.Debug0().LogDebug("getting request for calculatorFunction with payload", req)

	validationErrors := wscutils.WscValidate(req, func(err validator.FieldError) []string { return []string{} })
	validationErrors = append(validationErrors, customValidationErrors(&req)...)
	if len(validationErrors) > 0 {
		lh.Debug0().LogDebug("validations failed with the validation failure", req)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}

	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: func(req *CalculateRequest) int64 {
		a, _ := strconv.Atoi(req.Operand1)
		b, _ := strconv.Atoi(req.Operand2)

		switch OperationSymbol(req.Operation) {
		case Plus:
			return int64(a + b)
		case Minus:
			return int64(a - b)
		case Multiply:
			return int64(a * b)
		case Division:
			return int64(a / b)
		default:
			return math.MaxInt64 // this would never come as we are already doing sanity over the request body
		}
	}(&req), Messages: []wscutils.ErrorMessage{}})

}

func customValidationErrors(req *CalculateRequest) []wscutils.ErrorMessage {
	errResponseList := make([]wscutils.ErrorMessage, 0)

	operand2, _ := strconv.Atoi(req.Operand2)
	operand2StringName := "Operand2"

	operand1, _ := strconv.Atoi(req.Operand1)
	operand1StringName := "Operand1"

	if operand2 == 0 && req.Operation == Division {
		errCodeInvalid := "invalid"
		errResponseList = append(errResponseList, wscutils.BuildErrorMessage(MsgIdZeroNotAllowed, &errCodeInvalid, operand2StringName, Zero, DivisionKeyString))
	}

	// testing multiple cases // not a valid use case though // dont divide big numbers
	if operand1 > 1000 && req.Operation == Division {
		errCodeInvalid := "invalid"
		errResponseList = append(errResponseList, wscutils.BuildErrorMessage(MsgIdZeroNotAllowed, &errCodeInvalid, operand1StringName, "1000", DivisionKeyString))
	}

	return errResponseList
}
