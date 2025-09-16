package sdk

// MaskedBytes is a byte slice type that masks its content when serialized.
// This is used for sensitive data like passwords or secrets that should
// never be exposed in JSON responses or logs, even when accidentally serialized.
type MaskedBytes []byte

// MarshalJSON implements the json.Marshaler interface.
// It always returns a masked string regardless of the actual content,
// ensuring sensitive data is never accidentally exposed in JSON output.
func (m MaskedBytes) MarshalJSON() ([]byte, error) {
	return []byte(`"*****"`), nil
}

// String implements the fmt.Stringer interface.
// It returns a masked string to prevent sensitive data from being
// accidentally logged or displayed.
func (m MaskedBytes) String() string {
	return "*****"
}
