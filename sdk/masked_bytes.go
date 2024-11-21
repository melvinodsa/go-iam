package sdk

type MaskedBytes []byte

func (m MaskedBytes) MarshalJSON() ([]byte, error) {
	return []byte(`"*****"`), nil
}

func (m MaskedBytes) String() string {
	return "*****"
}
