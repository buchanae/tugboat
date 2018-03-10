package tugboat

import (
	"fmt"
)

func debug(msg string, args ...interface{}) {
	fmt.Println(msg, args)
}
