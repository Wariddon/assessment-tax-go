package tax

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
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
		{
			TotalIncome: 3000000.0, WHT: 0.0,
			Allowances: []Allowance{
				{AllowanceType: "donation", Amount: 0.0},
			},
		},
	}

	reqBodyJSON, _ := json.Marshal(cases[1])
	t.Run("should return tax : 639000", func(t *testing.T) {
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
			assert.Equal(t, float64(639000), resp.Tax)
		}

	})

	reqBodyJSON, _ = json.Marshal(cases[0])
	t.Run("should return tax : 29000", func(t *testing.T) {
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

func Test_UpdateAmountAllowance(t *testing.T) {

	cases := []UpdateDeduction{
		{Amount: 50000},
		{Amount: 60000},
	}

	reqBodyJSON, _ := json.Marshal(cases[0])
	t.Run("update personal amount is ok", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/k-receipt", bytes.NewBuffer(reqBodyJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)

		if assert.NoError(t, DeductionKreceipt(c)) {
			assert.Equal(t, http.StatusOK, res.Code)
			var resp map[string]float64
			err := json.Unmarshal(res.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, float64(50000), resp["kReceipt"])
		}

	})

	reqBodyJSON, _ = json.Marshal(cases[1])
	t.Run("update personal amount is ok", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", bytes.NewBuffer(reqBodyJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)

		if assert.NoError(t, DeductionPersonal(c)) {
			assert.Equal(t, http.StatusOK, res.Code)
			var resp map[string]float64
			json.Unmarshal(res.Body.Bytes(), &resp)
			assert.Equal(t, float64(60000), resp["personalDeduction"])
		}

	})

	t.Run("update personal amount is empty", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)

		if assert.NoError(t, DeductionPersonal(c)) {
			assert.Equal(t, http.StatusBadRequest, res.Code)
			var resp map[string]interface{}
			json.Unmarshal(res.Body.Bytes(), &resp)
			assert.Equal(t, "amount must between 0 and 100000", resp["message"])
		}

	})

	t.Run("update k-receipt amount is empty", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/k-receipt", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)

		if assert.NoError(t, DeductionKreceipt(c)) {
			assert.Equal(t, http.StatusBadRequest, res.Code)
			var resp map[string]interface{}
			json.Unmarshal(res.Body.Bytes(), &resp)
			assert.Equal(t, "amount must between 0 and 100000", resp["message"])
		}

	})

	reqBodyJSON, _ = json.Marshal(map[string]string{"amount": "error"})
	t.Run("update k-receipt format error", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/k-receipt", bytes.NewBuffer(reqBodyJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)

		if assert.NoError(t, DeductionKreceipt(c)) {
			assert.Equal(t, http.StatusBadRequest, res.Code)
		}
	})

	t.Run("update personal format error", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", bytes.NewBuffer(reqBodyJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)

		if assert.NoError(t, DeductionPersonal(c)) {
			assert.Equal(t, http.StatusBadRequest, res.Code)
		}
	})

}

func Test_CSV(t *testing.T) {

	csvContent :=
		`totalIncome,wht,donation
500000,0,0
600000,40000,20000
750000,50000,15000
`

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("taxFile", "csv.csv")
	part.Write([]byte(csvContent))
	writer.Close()

	t.Run("update personal format error", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)

		if assert.NoError(t, CalculationTaxCSV(c)) {
			assert.Equal(t, http.StatusOK, res.Code)
		}
	})

}
