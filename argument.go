package zbaction

type ExpanderFn func(string) string
type Argument[T any] interface {
	Value(expander ExpanderFn) T
	Raw() string
}

type argContainer struct {
	raw string
}

func (c argContainer) Value(expander ExpanderFn) string {
	return expander(c.raw)
}

func (c argContainer) Raw() string {
	return c.raw
}

type mappableArgContainer[T any] struct {
	argContainer
	mapper func(s string) T
}

func (c mappableArgContainer[T]) Value(expander ExpanderFn) T {
	return c.mapper(c.argContainer.Value(expander))
}

func (c mappableArgContainer[T]) Raw() string {
	return c.raw
}

func NewArgument[T any](value string, mapper func(s string) T) Argument[T] {
	return mappableArgContainer[T]{
		argContainer: argContainer{raw: value},
		mapper:       mapper,
	}
}

func NewArgumentStr(value string) Argument[string] {
	return argContainer{raw: value}
}

func NewArgumentBool(value string) Argument[bool] {
	return mappableArgContainer[bool]{
		argContainer: argContainer{raw: value},
		mapper: func(s string) bool {
			return s == "true" || s == "1" || s == "True" || s == "TRUE"
		},
	}
}
