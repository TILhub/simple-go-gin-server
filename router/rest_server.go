package router

import (
	"awesomeProject2/requestStructures"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/logharbour/logharbour"
	"math"
	"net/http"
	"strconv"
)

type RESTServer interface {
	Init()
}
type sampleRESTServer struct {
	router *gin.Engine
	logger *logharbour.Logger
}

func NewRESTServer(router *gin.Engine, logger *logharbour.Logger) RESTServer {
	return &sampleRESTServer{
		router: router,
		logger: logger,
	}
}

func (server *sampleRESTServer) Init() {
	server.router.GET("/", server.healthCheck)
	server.router.POST("/api/v1/calculator", server.calculatorFunction)
}

func (server *sampleRESTServer) healthCheck(c *gin.Context) {
	server.logger.Debug0().LogDebug("DEBUG healthCheck", "")
	c.JSON(http.StatusOK, requestStructures.APPInfo{
		Name: "WatchDog Bark",
	})
}

func (server *sampleRESTServer) calculatorFunction(c *gin.Context) {

	var req requestStructures.CalculateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		wscutils.SendErrorResponse(c,
			wscutils.NewResponse(
				wscutils.ErrorStatus,
				nil,
				[]wscutils.ErrorMessage{
					wscutils.ErrorMessage{
						MsgID:   1,
						ErrCode: wscutils.ErrorStatus,
						Vals:    []string{"unable bind json to a structure", err.Error()},
					},
				}))
		return
	}

	server.logger.Debug0().LogDebug("getting request for calculatorFunction with payload", req)

	validationErrors := wscutils.WscValidate(req, func(err validator.FieldError) []string { return []string{} })
	if len(validationErrors) > 0 {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	} else {
		wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: func(req *requestStructures.CalculateRequest) int64 {
			a, _ := strconv.Atoi(req.Operand1)
			b, _ := strconv.Atoi(req.Operand2)

			switch requestStructures.OperationSymbol(req.Operation) {
			case requestStructures.Plus:
				return int64(a + b)
			case requestStructures.Minus:
				return int64(a - b)
			case requestStructures.Multiply:
				return int64(a * b)
			default:
				return math.MaxInt64 // this would never come as we are already doing sanity over the request body
			}
		}(&req), Messages: []wscutils.ErrorMessage{}})
	}
}
