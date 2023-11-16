package models

import "github.com/sergeizaitcev/gophermart/pkg/monetary"

type Order struct {
	Number string  `json:"order"`
	Goods  []Goods `json:"goods"`
}

type Goods struct {
	Match string        `json:"description"`
	Price monetary.Unit `json:"price"`
}

type OrderOut struct {
	Number  string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}
