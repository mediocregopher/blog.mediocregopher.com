package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/api"
)

func main() {

	fmt.Fprint(os.Stderr, "Password: ")

	line, err := bufio.NewReader(os.Stdin).ReadString('\n')

	if err != nil {
		panic(err)
	}

	fmt.Println(api.NewPasswordHash(strings.TrimSpace(line)))
}
