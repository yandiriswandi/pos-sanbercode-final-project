package main

import (
	"fmt"

	"github.com/yandiriswandi/pos-sanbercode-final-project/config"
	"github.com/yandiriswandi/pos-sanbercode-final-project/routers"
)

func main() {
	config.InitDB()

	PORT := ":8080"

	routers.StartSever().Run(PORT)

	fmt.Println("Server running at http://localhost:8080")

}
