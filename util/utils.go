package util

func MaxLength(strings ...string) int {
	max := 0
	for _, s := range strings {
		if len(s) > max {
			max = len(s)
		}
	}
	return max
}
