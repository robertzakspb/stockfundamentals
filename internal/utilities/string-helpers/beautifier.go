package stringhelpers

import (
	"fmt"
	"math"
	"strings"

	typeconverter "github.com/compoundinvest/stockfundamentals/internal/utilities/converters"
)

// Returns a number in a concise human-readable format (e.g. 124502.32 -> +124 тыс.)
func BeautifyNumber(number any) (string, error) {
	float, err := typeconverter.GetFloat(number)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	if float >= 0 {
		sb.WriteString("+")
	}

	formattedFloat := ""
	if math.Abs(float) < 1000 {
		formattedFloat = fmt.Sprintf("%.2f", float)
	} else if math.Abs(float) >= 1000 && math.Abs(float) < 1_000_000 {
		formattedFloat = fmt.Sprintf("%.2f", float/1000)
	} else if math.Abs(float) >= 1_000_000 {
		formattedFloat = fmt.Sprintf("%.2f", float/1_000_000)
	}

	//Removing a potential and redundant ".0 " at the end
	if len(formattedFloat) > 0 && formattedFloat[len(formattedFloat)-3:] == ".00" {
		formattedFloat = formattedFloat[:len(formattedFloat)-3]
	}
	sb.WriteString(formattedFloat)

	if math.Abs(float) >= 1000 && math.Abs(float) < 1_000_000 {
		sb.WriteString(" тыс.")
	} else if math.Abs(float) >= 1_000_000 {
		sb.WriteString(" млн.")
	}

	return sb.String(), nil
}

// Converts values like 0.12 into concise +12%.
func BeatufityPercentage(percentage float64) (string, error) {
	var sb strings.Builder

	if percentage < 1.0 {
		if percentage >= 0 {
			sb.WriteString("+")
		}
		fmt.Fprintf(&sb, "%.1f", percentage*100)
		sb.WriteString("%")
	} else {
		sb.WriteString("x" + fmt.Sprintf("%.1f", percentage+1.0))
	}

	return sb.String(), nil
}
