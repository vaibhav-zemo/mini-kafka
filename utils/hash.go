package utils

func Hash(key string) int {
	h := 0
	for _, c := range key {
		h = (h*31 + int(c)) % 1000000007
	}
	return h
}
