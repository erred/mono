package envconf

import "os"

func String(name, value string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		return value
	}
	return v
}
