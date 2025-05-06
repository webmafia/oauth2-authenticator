package auth

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// An integer that encodes to (and decodes from) a JSON string.
type TextualInt int

var _ json.Marshaler = TextualInt(0)
var _ json.Unmarshaler = (*TextualInt)(nil)

func (t TextualInt) String() string {
	return strconv.Itoa(int(t))
}

// MarshalJSON implements json.Marshaler.
func (t TextualInt) MarshalJSON() ([]byte, error) {
	// Convert the int to string, then marshal it as a JSON string
	return json.Marshal(t.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *TextualInt) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("TextualInt expects a JSON string: %w", err)
	}

	n, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("invalid integer in string: %w", err)
	}

	*t = TextualInt(n)
	return nil
}
