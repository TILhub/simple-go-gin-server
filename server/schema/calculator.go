package schema

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"math"
	"strconv"
)

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
	if len(validationErrors) > 0 {

		lh.Debug0().LogDebug("validations failed with the validation failure", struct {
			request    CalculateRequest
			errMessage []wscutils.ErrorMessage
		}{
			request:    req,
			errMessage: validationErrors,
		})

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
		default:
			return math.MaxInt64 // this would never come as we are already doing sanity over the request body
		}
	}(&req), Messages: []wscutils.ErrorMessage{}})

}
