package utils

import (
	"bytes"
	"fmt"
	"os"
)

func TrimEnv(i string) string {
	return string(bytes.TrimSpace([]byte(fmt.Sprintf("%s", os.Getenv(i)))))
}
