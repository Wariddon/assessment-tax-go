package tax

type Tax struct {
	Tax      float64    `json:"tax"`
	TaxLevel []TaxLevel `json:"taxLevel"`
}

type TaxRefund struct {
	TaxRefund float64    `json:"taxRefund"`
	TaxLevel  []TaxLevel `json:"taxLevel"`
}

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type TaxesCSV struct {
	Taxes []interface{} `json:"taxes"`
}

type TaxCSV struct {
	TotalIncome float64 `json:"totalIncome"`
	Tax         float64 `json:"tax"`
}

type TaxRefundCSV struct {
	TotalIncome float64 `json:"totalIncome"`
	TaxRefund   float64 `json:"taxRefund"`
}
