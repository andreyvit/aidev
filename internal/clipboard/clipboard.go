//go:build !darwin

package clipboard

import (
	"errors"
)

func CopyText(content string) error {
	return errors.New("copying to clipboard not supported")
}
