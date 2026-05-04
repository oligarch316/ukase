package uktodo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// TODO: Clean up/formalize all this ArgRange junk

type ArgRange struct{ Low, High int }

func (ar ArgRange) Contains(position int) bool {
	lowOk := ar.Low < 0 || ar.Low <= position
	highOk := ar.High < 0 || ar.High >= position
	return lowOk && highOk
}

func FormatArgRange(ar ArgRange) string {
	// TODO: This especially is currently terrible

	switch noLow, noHigh := ar.Low < 0, ar.High < 0; {
	case noLow && noHigh:
		return "[...]"
	case noLow:
		return fmt.Sprintf("[...%d]", ar.High)
	case noHigh:
		return fmt.Sprintf("[%d...]", ar.Low)
	case ar.Low == ar.High:
		return fmt.Sprintf("[%d]", ar.Low)
	default:
		return fmt.Sprintf("[%d...%d]", ar.Low, ar.High)
	}
}

func ParseArgRange(s string) (ArgRange, error) {
	// ""    ⇒ [-∞,∞]
	// "#"   ⇒ [#,#]
	// ":"   ⇒ [-∞,∞]
	// ":#"  ⇒ [-∞,#-1]
	// "#:"  ⇒ [#,∞]
	// "#:#" ⇒ [#,#-1]

	const lowMin = 0
	const highMin = 1

	s = strings.TrimSpace(s)
	lowS, highS, sepFound := strings.Cut(s, ":")

	if !sepFound {
		val, err := parseArgVal(s, lowMin)
		argRange := ArgRange{Low: val, High: val}
		return argRange, err
	}

	lowVal, lowErr := parseArgVal(lowS, lowMin)
	highVal, highErr := parseArgVal(highS, highMin)
	if err := errors.Join(lowErr, highErr); err != nil {
		return ArgRange{}, err
	}

	if lowVal >= lowMin && highVal >= highMin && lowVal >= highVal {
		return ArgRange{}, errors.New("invalid range (low >= high)")
	}

	argRange := ArgRange{Low: lowVal, High: highVal - 1}
	return argRange, nil
}

func parseArgVal(s string, min int) (int, error) {
	if s == "" {
		return min - 1, nil
	}

	val, err := strconv.Atoi(s)
	if err != nil {
		return -2, err
	}

	if val < min {
		return -2, errors.New("below minimum")
	}

	return val, nil
}
