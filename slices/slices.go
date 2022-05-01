package slices

type UniqueIdSlice []string

func (u UniqueIdSlice) Add(newId string) UniqueIdSlice {
	for _, id := range u {
		if id == newId {
			return u
		}
	}
	return append(u, newId)
}

func ArrayFromMap(m map[string]interface{}) []interface{} {
	var arr []interface{}
	for k, v := range m {
		arr = append(arr, k)
		arr = append(arr, v)
	}
	if len(arr) > 0 {
		return arr
	}
	return nil
}

func ArrayStringFromMap(m map[string]string) []string {
	var arr []string
	for k, v := range m {
		arr = append(arr, k)
		arr = append(arr, v)
	}
	if len(arr) > 0 {
		return arr
	}
	return nil
}

func GetStringSliceIntersection(s1, s2 []string) (inter []string) {
	hash := make(map[string]bool)
	for _, e := range s1 {
		hash[e] = true
	}
	for _, e := range s2 {
		if hash[e] {
			inter = append(inter, e)
			hash[e] = false
		}
	}
	return
}

func IsStrExist(array []string, predicate string) bool {
	for _, e := range array {
		if e == predicate {
			return true
		}
	}
	return false
}
