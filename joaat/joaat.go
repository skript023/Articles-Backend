package joaat

import "strings"

func Hash(name string) uint32 {
	var hash uint32 = 0
	str := strings.ToLower(name)
	for _, c := range str {
		hash += uint32(c)
		hash += hash << 10
		hash ^= hash >> 6
	}

	hash += hash << 3
	hash ^= hash >> 11
	hash += hash << 15

	return hash
}
