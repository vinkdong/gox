package slice

type Runtime struct {
}

func Difference(sliceA []string, sliceB []string) []string {
	diff := make([]string, 0)
	diffMap := make(map[string]int)

	for _, v := range sliceA {
		diffMap[v] = 1
	}
	for _, v := range sliceB {
		diffMap[v] = diffMap[v] - 1
	}

	for k, v := range diffMap {
		if v > 0 {
			diff = append(diff, k)
		}
	}
	return diff
}

func UnionString(slices ...[]string) []string {
	m := make(map[string]bool)
	for _, s := range slices {
		for _, val := range s {
			m[val] = true
		}
	}
	s := make([]string, len(m))
	i := 0
	for k, _ := range m {
		s[i] = k
		i++
	}
	return s
}

func IsContain(sliceA []string, element string) bool {
	for _, v := range sliceA {
		if v == element {
			return true
		}
	}
	return false
}
