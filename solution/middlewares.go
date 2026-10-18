package rates

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrTemporaryFailure = errors.New("temporary failure")
	ErrValidatorPanic   = errors.New("validator panic")
)

// StaleDataError сообщает, что данные об изменении курса устарели.
type StaleDataError struct {
	Rate     ExchangeRate
	StaleFor time.Duration
}

func (e *StaleDataError) Error() string {
	return fmt.Sprintf(
		"exchange-rate data for %s is stale by %s",
		e.Rate.Currency,
		e.StaleFor,
	)
}

// FetchFunc возвращает значение из произвольного источника данных.
type FetchFunc[T any] func() (T, error)

// Validator проверяет значение, возвращённое FetchFunc.
type Validator[T any] interface {
	Validate(T) error
}

// Logger записывает форматированные сообщения.
type Logger interface {
	Printf(format string, args ...any)
}

func validate[T any](validator Validator[T], value T) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("%w: %v", ErrValidatorPanic, recovered)
		}
	}()

	return validator.Validate(value)
}

// WithValidate добавляет к next последовательную проверку валидаторами.
func WithValidate[T any](
	next FetchFunc[T],
	validators ...Validator[T],
) FetchFunc[T] {
	return func() (T, error) {
		value, err := next()
		if err != nil {
			var zero T
			return zero, err
		}

		for index, validator := range validators {
			if err := validate(validator, value); err != nil {
				var zero T
				return zero, fmt.Errorf("validator %d: %w", index, err)
			}
		}

		return value, nil
	}
}

// WithRetry добавляет к next повторные попытки для временных ошибок.
func WithRetry[T any](next FetchFunc[T], retries int) FetchFunc[T] {
	return func() (T, error) {
		for attempt := 0; attempt <= retries; attempt++ {
			value, err := next()
			if err == nil {
				return value, nil
			}

			if attempt == retries || !shouldRetry(err) {
				var zero T
				return zero, err
			}
		}

		panic("unreachable")
	}
}

func shouldRetry(err error) bool {
	if errors.Is(err, ErrValidatorPanic) {
		return false
	}

	var staleErr *StaleDataError
	if errors.As(err, &staleErr) {
		return staleErr.StaleFor < 5*time.Minute
	}

	return errors.Is(err, ErrTemporaryFailure)
}

// WithLogs записывает начало и завершение вызова next.
func WithLogs[T any](next FetchFunc[T], logger Logger) FetchFunc[T] {
	return func() (T, error) {
		logger.Printf("функция начата в %s", time.Now().Format(time.RFC3339Nano))
		defer func() {
			logger.Printf("функция завершена в %s", time.Now().Format(time.RFC3339Nano))
		}()

		return next()
	}
}
