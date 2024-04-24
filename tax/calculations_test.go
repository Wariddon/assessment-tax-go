package tax

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_taxCalculations(t *testing.T) {

	cases := []Calculations{
		{
			TotalIncome: 500000.0, WHT: 0.0,
			Allowances: []Allowance{
				{AllowanceType: "donation", Amount: 0.0},
			},
		},
	}
	
	reqBodyJSON, _ := json.Marshal(cases[0])
	t.Run("should return 200 status ok", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBuffer(reqBodyJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)

		if assert.NoError(t, CalculationTax(c)) {
			assert.Equal(t, http.StatusOK, res.Code)
			var resp Tax
			err := json.Unmarshal(res.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, float64(29000), resp.Tax)
		}

	})
}
