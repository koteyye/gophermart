package monetary

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

// Precision определяет точность денежной единицы.
const Precision = 2

var base = math.Pow10(Precision)

var (
	_ json.Marshaler   = (*Unit)(nil)
	_ json.Unmarshaler = (*Unit)(nil)
)

// Unit определяет денежную единицу с точностью Precision.
type Unit int64

// Format форматирует f в денежную единицу с точностью precision и возвращает её.
//
// Специальные кейсы:
//
//	Format(±0) = 0
//	Format(±Inf) = 0
//	Format(NaN) = 0
//
// FIXME: не работает с большими числами, такими как math.MaxFloat64;
// требуется доработка.
func Format(f float64) Unit {
	if math.IsInf(f, 0) || math.IsNaN(f) || f == 0 || f == -0 || f == +0 {
		return 0
	}
	return Unit(math.Round(f * base))
}

func (u Unit) String() string {
	return strconv.FormatFloat(u.Float64(), 'f', -1, 64)
}

func (u Unit) Float64() float64 {
	return float64(u) / base
}

func (u Unit) MarshalJSON() ([]byte, error) {
	f := u.Float64()
	return json.Marshal(f)
}

func (u *Unit) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" {
		*u = 0
		return nil
	}
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("parse float64: %w", err)
	}
	*u = Format(f)
	return nil
}

var _ driver.Valuer = (*NullUnit)(nil)

// NullUnit определяет денежную единицу для использования в пакете
// `database/sql`.
type NullUnit struct {
	Unit  Unit
	Valid bool

	value sql.NullInt64
}

func (n *NullUnit) Scan(value any) error {
	err := n.value.Scan(value)
	if err != nil {
		return err
	}
	n.Unit, n.Valid = Unit(n.value.Int64), n.value.Valid
	return nil
}

func (n *NullUnit) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return int64(n.Unit), nil
}
