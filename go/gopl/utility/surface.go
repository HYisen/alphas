package utility

import (
	"bufio"
	"fmt"
	"os"
)

func RequireInput(hint string) string {
	fmt.Print(hint)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text()
	}
	panic(fmt.Errorf("can not get input with hint %s:%v", hint, scanner.Err()))
}
