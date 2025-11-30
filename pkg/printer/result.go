package printer

import (
	"fmt"
	"io"
	"strings"
)

func Print(output io.Writer, result map[string][]string) error {
	sb := strings.Builder{}

	sb.WriteString("------------------\n")
	for i := range result {
		if len(result[i]) == 1 {
			continue
		}

		sb.WriteString("--- Same files ---\n")
		for j := range result[i] {
			sb.WriteString(fmt.Sprintf("file: %s\n", result[i][j]))
		}
		sb.WriteString("\n")
	}
	sb.WriteString("------------------\n")

	_, err := output.Write([]byte(sb.String()))
	if err != nil {
		return err
	}

	return nil
}
