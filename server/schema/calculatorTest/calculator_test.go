package calculatorTest

import (
	"awesomeProject2/server/schema"
	"awesomeProject2/testutils"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/logharbour/logharbour"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func registerRoutes() (*gin.Engine, error) {
	// router
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// logger setup
	fallbackWriter := logharbour.NewFallbackWriter(os.Stdout, os.Stdout)
	lctx := logharbour.NewLoggerContext(logharbour.Info)
	l := logharbour.NewLogger(lctx, "crux", fallbackWriter)

	s := service.NewService(r).
		WithLogHarbour(l)

	s.RegisterRoute(http.MethodPost, "/api/v1/calculator", schema.CalculatorFunction)

	return r, nil

}

func TestSchemaNew(t *testing.T) {
	//gin.SetMode(gin.TestMode)
	//r := gin.Default()
	r, _ := registerRoutes()
	testCases := schemaNewTestcase()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Setting up buffer
			payload := bytes.NewBuffer(testutils.MarshalJson(tc.RequestPayload))

			res := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/api/v1/calculator", payload)
			require.NoError(t, err)

			r.ServeHTTP(res, req)

			require.Equal(t, tc.ExpectedHttpCode, res.Code)
			if tc.ExpectedResult != nil {
				jsonData := testutils.MarshalJson(tc.ExpectedResult)
				require.JSONEq(t, string(jsonData), res.Body.String())
			} else {
				jsonData, err := testutils.ReadJsonFromFile(tc.TestJsonFile)
				require.NoError(t, err)
				require.JSONEq(t, string(jsonData), res.Body.String())
			}
		})
	}
}

func schemaNewTestcase() []testutils.TestCasesStruct {
	valTestJson, err := testutils.ReadJsonFromFile("./testData/ok.json")
	if err != nil {
		log.Fatalln("Error reading JSON file:", err)
	}
	var ok1 schema.CalculateRequest
	if err := json.Unmarshal(valTestJson, &ok1); err != nil {
		log.Fatalln("Error unmarshalling JSON:", err)
	}

	errCase, err := testutils.ReadJsonFromFile("./testData/error.json")
	if err != nil {
		log.Fatalln("Error reading JSON file:", err)
	}
	var cusValPayload schema.CalculateRequest
	if err := json.Unmarshal(errCase, &cusValPayload); err != nil {
		log.Fatalln("Error unmarshalling JSON:", err)
	}

	schemaNewTestcase := []testutils.TestCasesStruct{
		{
			Name: "Sunny_Day_Case",
			RequestPayload: wscutils.Request{
				Data: ok1,
			},

			ExpectedHttpCode: http.StatusOK,
			ExpectedResult: &wscutils.Response{
				Status:   wscutils.SuccessStatus,
				Data:     5550,
				Messages: []wscutils.ErrorMessage{},
			},
		},
		{
			Name: "err- standard validation failure ",
			RequestPayload: wscutils.Request{
				Data: cusValPayload,
			},
			TestJsonFile:     "./testData/error_response.json",
			ExpectedHttpCode: http.StatusBadRequest,
		},
	}
	return schemaNewTestcase
}
