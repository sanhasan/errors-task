package rates

import "time"

type ExchangeRate struct {
	Currency         string
	ChangeMultiplier float64
	ChangedAt        time.Time
}

type freshnessValidator struct {
	maxAge time.Duration
}

func NewFreshnessValidator(maxAge time.Duration) Validator[ExchangeRate] {
	panic("not implemented")
}

func (v freshnessValidator) Validate(rate ExchangeRate) error {
	panic("not implemented")
}

type fieldsValidator struct {
	supportedCurrencies map[string]struct{}
}

func NewFieldsValidator(supportedCurrencies map[string]struct{}) Validator[ExchangeRate] {
	panic("not implemented")
}

func (v fieldsValidator) Validate(rate ExchangeRate) error {
	panic("not implemented")
}
