package cli

// ArrayFlags implements flag.Value for multiple string values
type ArrayFlags []string

func (i *ArrayFlags) String() string {
	return "array flags"
}

func (i *ArrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}