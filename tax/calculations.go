package tax

import (
	"encoding/csv"
	"io"
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Calculations struct {
	TotalIncome float64     `json:"totalIncome"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type CalculationsCSV struct {
	TotalIncome float64 `json:"totalIncome"`
	WHT         float64 `json:"wht"`
	Allowances  float64 `json:"allowances"`
}

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type UpdateDeduction struct {
	Amount float64 `json:"amount" validate:"required,min=0,max=100000"`
}

func DeductionKreceipt(c echo.Context) error {
	de := UpdateDeduction{}
	err := c.Bind(&de)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if de.Amount > 100000 || de.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "amount must between 0 and 100000"})
	}

	err = UpdateAllowances("k-receipt", de.Amount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]float64{
		"kReceipt": de.Amount,
	})

}

func DeductionPersonal(c echo.Context) error {
	de := UpdateDeduction{}
	err := c.Bind(&de)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if de.Amount > 100000 || de.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "amount must between 0 and 100000"})
	}

	err = UpdateAllowances("personal", de.Amount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]float64{
		"personalDeduction": de.Amount,
	})

}

func CalculationTax(c echo.Context) error {
	cal := Calculations{}
	err := c.Bind(&cal)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	totalIncome := cal.TotalIncome

	totalIncome, _ = deductionAllowance(totalIncome, cal.Allowances)

	calTax := 0.0

	var taxDetails []TaxLevel
	taxDetails, calTax = calculationTax(totalIncome)

	// Cal WHT
	wht := cal.WHT
	calTax -= wht
	if calTax > 0 {
		return c.JSON(http.StatusOK, Tax{Tax: calTax, TaxLevel: taxDetails})
	} else {

		return c.JSON(http.StatusOK, TaxRefund{TaxRefund: math.Abs(calTax), TaxLevel: taxDetails})
	}

	//fmt.Println(cal.Allowances[0].Amount)

}

func CalculationTaxCSV(c echo.Context) error {
	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["taxFile"]

	src, err := files[0].Open()
	if err != nil {
		return err
	}
	defer src.Close()

	taxes, err := processTaxCSV(src)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, taxes)

}
func deductionAllowance(income float64, allowance []Allowance) (float64, error) {
	//personal deduction tax
	allowances, err := GetAllowances()
	if err != nil {
		return 0, err
	}

	// deduction personal
	income -= allowances["personal"]
	delete(allowances, "personal")

	// other deduction if you have
	for _, allowance := range allowance {
		if MaxAmount, ok := allowances[allowance.AllowanceType]; ok {

			calAmount := allowance.Amount
			if allowance.Amount > MaxAmount {
				calAmount = MaxAmount
			}
			income -= calAmount
		}
	}

	return income, nil
}

/*
รายได้ 0 - 150,000 ได้รับการยกเว้น
150,001 - 500,000 อัตราภาษี 10%
500,001 - 1,000,000 อัตราภาษี 15%
1,000,001 - 2,000,000 อัตราภาษี 20%
มากกว่า 2,000,000 อัตราภาษี 35%
*/
func calculationTax(income float64) ([]TaxLevel, float64) {

	taxDetails := []TaxLevel{
		{Level: "0-150,000", Tax: 0.0},
		{Level: "150,001-500,000", Tax: 0.0},
		{Level: "500,001-1,000,000", Tax: 0.0},
		{Level: "1,000,001-2,000,000", Tax: 0.0},
		{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
	}

	totalTax := 0.0
	if income > 2000000 {
		tax := (income - 2000000) * 0.35
		taxDetails[4].Tax = tax
		totalTax += tax
		income = 2000000
	}
	if income > 1000000 {
		tax := (income - 1000000) * 0.20
		taxDetails[3].Tax = tax
		totalTax += tax
		income = 1000000
	}
	if income > 500000 {
		tax := (income - 500000) * 0.15
		taxDetails[2].Tax = tax
		totalTax += tax
		income = 500000
	}
	if income > 150000 {
		tax := (income - 150000) * 0.10
		taxDetails[1].Tax = tax
		totalTax += tax
		income = 150000
	}

	return taxDetails, totalTax
}

func processTaxCSV(file io.Reader) (TaxesCSV, error) {

	var taxesInterface []interface{}
	reader := csv.NewReader(file)

	// Skip the header row
	header, err := reader.Read()
	if err != nil {
		return TaxesCSV{}, err
	}

	// Read each record
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return TaxesCSV{}, err
		}

		totalIncome, _ := strconv.ParseFloat(record[0], 64)
		wht, _ := strconv.ParseFloat(record[1], 64)
		amount, _ := strconv.ParseFloat(record[2], 64)

		income := totalIncome

		//Make Object
		allowance := []Allowance{
			{AllowanceType: header[2], Amount: amount},
		}

		income, _ = deductionAllowance(income, allowance)

		calTax := 0.0

		_, calTax = calculationTax(income)

		// Cal WHT
		calTax -= wht
		if calTax > 0 {
			taxData := TaxCSV{
				TotalIncome: totalIncome,
				Tax:         calTax,
			}
			taxesInterface = append(taxesInterface, taxData)
		} else {
			taxData := TaxRefundCSV{
				TotalIncome: totalIncome,
				TaxRefund:   math.Abs(calTax),
			}
			taxesInterface = append(taxesInterface, taxData)
		}
	}

	taxes := TaxesCSV{Taxes: taxesInterface}

	return taxes, nil
}
