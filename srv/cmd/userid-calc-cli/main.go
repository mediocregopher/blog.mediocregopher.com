package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/chat"
)

func main() {

	secret := flag.String("secret", "", "Secret to use when calculating UserIDs")
	name := flag.String("name", "", "")
	password := flag.String("password", "", "")
	flag.Parse()

	calc := chat.NewUserIDCalculator([]byte(*secret))
	userID := calc.Calculate(*name, *password)

	b, err := json.Marshal(userID)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
