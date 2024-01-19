package utils

import "fmt"

func ToAddress(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}
