package utility

import (
	"fmt"
	"io"
	"os"
)

func CloseAndLogError(c io.Closer) {
	if err := c.Close(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "meet error %v, but continue.", err)
	}
}
