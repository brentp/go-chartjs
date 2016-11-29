package types

import (
	"fmt"
	"image/color"
)

// RGBA amends image/color.RGBA to have a MarshalJSON that meets the expectations of chartjs.
type RGBA color.RGBA

// MarshalJSON satisfies the json.Marshaler interface.
func (c RGBA) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"rgba(%d, %d, %d, %.3f)\"", c.R, c.G, c.B, float64(c.A)/255)), nil
}

// Bool is a convenience typedef for pointer to bool so that we can differentiate between unset
// and false.
type Bool *bool

var (
	t = true
	f = false
	// True is a convenience for pointer to true
	True = Bool(&t)

	// False is a convenience for pointer to false
	False = Bool(&f)
)
