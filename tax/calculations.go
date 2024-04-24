package tax

import (
	"fmt"
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Calculations struct {
	TotalIncome float64     `json:"totalIncome"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

func CalculationTax(c echo.Context) error {
	cal := Calculations{}
	err := c.Bind(&cal)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	totalIncome := cal.TotalIncome

	calTax := 0.0

	//personal deduction tax

	allowances, err := GetAllowances()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	// deduction personal
	totalIncome -= allowances["personal"]

	delete(allowances, "personal")

	// other deduction if you have
	for _, allowance := range cal.Allowances {
		if MaxAmount, ok := allowances[allowance.AllowanceType]; ok {

			calAmount := allowance.Amount
			if allowance.Amount > MaxAmount {
				calAmount = MaxAmount
			}
			fmt.Println("calAmount ", calAmount)
			totalIncome -= calAmount
		}

	}

	_, calTax = calculationTax(totalIncome)

	// Cal WHT
	wht := cal.WHT
	calTax -= wht
	if calTax > 0 {
		return c.JSON(http.StatusOK, map[string]float64{"tax": calTax})
	} else {
		return c.JSON(http.StatusOK, map[string]float64{"taxRefund": math.Abs(calTax)})
	}

	//fmt.Println(cal.Allowances[0].Amount)

}

/*
รายได้ 0 - 150,000 ได้รับการยกเว้น
150,001 - 500,000 อัตราภาษี 10%
500,001 - 1,000,000 อัตราภาษี 15%
1,000,001 - 2,000,000 อัตราภาษี 20%
มากกว่า 2,000,000 อัตราภาษี 35%
*/
func calculationTax(income float64) ([]TaxLevel, float64) {

	fmt.Println("income", income)

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
