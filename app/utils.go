package main

func unpack(s []string, vars ...*string) {
	for i := range vars {
		*vars[i] = s[i]
	}
}
