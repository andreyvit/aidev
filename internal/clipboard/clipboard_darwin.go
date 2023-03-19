package clipboard

import (
	"fmt"
	"os/exec"
	"strings"
)

func CopyText(content string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(content)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("pbcopy: %w", err)
	}
	return nil
}
