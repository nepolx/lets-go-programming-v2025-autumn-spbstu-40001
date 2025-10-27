package currency

import (
	"fmt"
	"strconv"
	"strings"
)

type Rates struct {
	Data []Currency `xml:"Valute"`
}

type (
	FloatforCur float64
	Currency    struct {
		NumCode  int         `json:"num_code"  xml:"NumCode"`
		CharCode string      `json:"char_code" xml:"CharCode"`
		Value    FloatforCur `json:"value"     xml:"Value"`
	}
)

func (cf *FloatforCur) UnmarshalText(text []byte) error {
	input := strings.TrimSpace(string(text))
	if input == "" {
		return fmt.Errorf("empty number")
	}

	normalized := strings.Replace(input, ",", ".", 1)
	if strings.Contains(normalized, ",") {
		return fmt.Errorf("multiple decimal separators in %q", text)
	}

	value, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return fmt.Errorf("invalid number %q: %v", text, err)
	}

	*cf = FloatforCur(value)

	return nil
}
