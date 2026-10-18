package rates

import "errors"

var (
	ErrTemporaryFailure = errors.New("temporary failure")
	ErrValidatorPanic   = errors.New("validator panic")
)

// StaleDataError сообщает, что данные об изменении курса устарели.
type StaleDataError struct{}

func (e *StaleDataError) Error() string {
	panic("not implemented")
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

// WithValidate добавляет к next последовательную проверку валидаторами.
func WithValidate[T any](
	next FetchFunc[T],
	validators ...Validator[T],
) FetchFunc[T] {
	panic("not implemented")
}

// WithRetry добавляет к next повторные попытки для временных ошибок.
func WithRetry[T any](next FetchFunc[T], retries int) FetchFunc[T] {
	panic("not implemented")
}

// WithLogs записывает начало и завершение вызова next.
func WithLogs[T any](next FetchFunc[T], logger Logger) FetchFunc[T] {
	panic("not implemented")
}
