package models

type Order struct {
	Number string `json:"order"`
	Goods  []Goods
}

type Goods struct {
	Match string `json:"description"`
	Price int    `json:"price"`
}

type OrderOut struct {
	Number  string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}
