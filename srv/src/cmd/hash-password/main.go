package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/http"
)

func main() {

	fmt.Fprint(os.Stderr, "Password: ")

	line, err := bufio.NewReader(os.Stdin).ReadString('\n')

	if err != nil {
		panic(err)
	}

	fmt.Println(http.NewPasswordHash(strings.TrimSpace(line)))
}
