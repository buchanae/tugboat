package docker

import (
	"fmt"
	"math/rand"
	"time"
)

func formatVolumeArg(host, container string, readonly bool) string {
	mode := "rw"
	if readonly {
		mode = "ro"
	}
	return fmt.Sprintf("%s:%s:%s", host, container, mode)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
