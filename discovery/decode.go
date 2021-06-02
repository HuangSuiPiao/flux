package discovery

import (
	"fmt"
	"strings"
)

func VerifyJSON(bytes []byte) error {
	size := len(bytes)
	if size < len("{\"k\":0}") {
		return fmt.Errorf("CHECK/JSON/MALFORMED/size: %d", size)
	}
	prefix := strings.TrimSpace(string(bytes[:5]))
	if prefix[0] != '[' && prefix[0] != '{' {
		return fmt.Errorf("CHECK/JSON/MALFORMED/start-token: %s", prefix)
	}
	return nil
}
