package rates

import (
	"errors"
	"fmt"
	"time"
)

type ExchangeRate struct {
	Currency         string
	ChangeMultiplier float64
	ChangedAt        time.Time
}

type freshnessValidator struct {
	maxAge time.Duration
}

func NewFreshnessValidator(maxAge time.Duration) Validator[ExchangeRate] {
	return freshnessValidator{maxAge: maxAge}
}

func (v freshnessValidator) Validate(rate ExchangeRate) error {
	age := time.Since(rate.ChangedAt)
	if age <= v.maxAge {
		return nil
	}

	return &StaleDataError{
		Rate:     rate,
		StaleFor: age - v.maxAge,
	}
}

type fieldsValidator struct {
	supportedCurrencies map[string]struct{}
}

func NewFieldsValidator(supportedCurrencies map[string]struct{}) Validator[ExchangeRate] {
	return fieldsValidator{supportedCurrencies: supportedCurrencies}
}

func (v fieldsValidator) Validate(rate ExchangeRate) error {
	switch {
	case rate.Currency == "":
		return errors.New("currency is empty")
	case !contains(v.supportedCurrencies, rate.Currency):
		return fmt.Errorf("currency %q is not supported", rate.Currency)
	case rate.ChangeMultiplier <= 0:
		return errors.New("change multiplier must be positive")
	default:
		return nil
	}
}

func contains(values map[string]struct{}, value string) bool {
	_, ok := values[value]
	return ok
}
