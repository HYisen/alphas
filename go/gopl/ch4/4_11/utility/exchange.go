package utility

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func GetInputFromTextEditor(filename, message string) string {
	path := os.TempDir() + "/" + filename

	createFile(path, message)

	modifyFile(path)

	return readFile(path)
}

func modifyFile(path string) {
	editor := os.Getenv("EDITOR")
	if len(editor) == 0 {
		editor = "vim"
	}

	cmd := exec.Command(editor, path)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

func readFile(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func createFile(path string, message string) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.WriteString(message)
	if err != nil {
		log.Fatal(err)
	}

	_ = f.Close()
}
