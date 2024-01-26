package zbaction

type ErrRequiredArgument struct {
	Key string
}

func NewErrRequiredArgument(key string) ErrRequiredArgument {
	return ErrRequiredArgument{
		Key: key,
	}
}

func (r ErrRequiredArgument) Error() string {
	return "missing required argument: " + r.Key
}
