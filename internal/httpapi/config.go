package httpapi

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

func Address(addr string) string {
	if addr != "" {
		return addr
	}
	if p := os.Getenv("PORT"); p != "" {
		if n, e := strconv.Atoi(p); e == nil && n > 0 && n < 65536 {
			return fmt.Sprintf("127.0.0.1:%d", n)
		}
	}
	return "127.0.0.1:19081"
}
func AddrFlag() *string { return flag.String("addr", "127.0.0.1:19081", "监听地址") }
