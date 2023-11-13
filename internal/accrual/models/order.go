package models

type Order struct {
	Number string `json:"order"`
	Goods  []Goods
}

type Goods struct {
	Match string `json:"description"`
	Price int    `json:"price"`
}
