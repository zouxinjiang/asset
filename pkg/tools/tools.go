package tools

func Int64SliceDuplicate(data []int64) []int64 {
	m := map[int64]bool{}
	result := make([]int64, 0, len(data))
	for _, v := range data {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			continue
		}
		m[v] = true
	}
	return result
}
