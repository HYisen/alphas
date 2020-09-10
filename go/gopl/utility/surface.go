package utility

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func RequireInput(hint string) string {
	fmt.Print(hint)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text()
	}
	panic(fmt.Errorf("can not get input with hint %s:%v", hint, scanner.Err()))
}

func ExecAndGetStdOut(args ...string) string {
	cmd := exec.Command(args[0], args[1:]...)
	var sb strings.Builder
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, &sb, os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	return sb.String()
}
