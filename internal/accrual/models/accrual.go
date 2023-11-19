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
	Accrual monetary.NullUnit    `json:"accrual"`
}

type Match struct {
	MatchName string `json:"match"`
	Reward monetary.Unit `json:"reward"`
	RewardType RewardType `json:"reward_type"`
}

type RewardType string

func (r RewardType) Uint() uint8 {
	switch r {
	case "%":
		return 0
	case "pt":
		return 1
	}
	return 2
}