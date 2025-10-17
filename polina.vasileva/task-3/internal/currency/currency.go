package currency

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrEmptyNumber = errors.New("empty number")
	ErrMultipleSeparators = errors.New("multiple decimal separators")
	ErrInvalidNumber = errors.New("invalid number")
)

type Rates struct {
	Data []Currency `xml:"Valute"`
}

type (
	CommaFloat float64
	Currency struct {
		NumCode  int        `json:"num_code"  xml:"NumCode"`
		CharCode string     `json:"char_code" xml:"CharCode"`
		Value    CommaFloat `json:"value"     xml:"Value"`
	}
)

func (cf *CommaFloat) UnmarshalText(text []byte) error {
	str := strings.TrimSpace(string(text))
	if str == "" {
		return ErrEmptyNumber
	}

	str = strings.Replace(str, ",", ".", 1)
	if strings.Count(str, ",") > 0 {
		return fmt.Errorf("%w: %q", ErrMultipleSeparators, text)
	}

	v, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("%w: %q: %w", ErrInvalidNumber, text, err)
	}

	*cf = CommaFloat(v)

	return nil
}

func DescendingComparatorCurrency(a, b Currency) int {
	floatA, floatB := float64(a.Value), float64(b.Value)

	switch {
	case floatB < floatA:
		return -1
	case floatB > floatA:
		return 1
	default:
		return 0
	}
}
