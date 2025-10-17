package currency

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrEmptyNumber        = errors.New("empty number")
	ErrMultipleSeparators = errors.New("multiple decimal separators")
	ErrInvalidNumber      = errors.New("invalid number")
)

type Rates struct {
	Data []Currency `xml:"Valute"`
}

type (
	CommaFloat float64
	Currency   struct {
		NumCode  int        `json:"num_code"  xml:"NumCode"`
		CharCode string     `json:"char_code" xml:"CharCode"`
		Value    CommaFloat `json:"value"     xml:"Value"`
	}
)

func (cf *CommaFloat) UnmarshalText(text []byte) error {
	input := strings.TrimSpace(string(text))
	if input == "" {
		return ErrEmptyNumber
	}

	normalized := strings.Replace(input, ",", ".", 1)
	if strings.Contains(normalized, ",") {
		return fmt.Errorf("%w: %q", ErrMultipleSeparators, text)
	}

	value, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return fmt.Errorf("%w: %q: %w", ErrInvalidNumber, text, err)
	}

	*cf = CommaFloat(value)
	
	return nil
}

func ComparatorCurrency(a, b Currency) int {
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
