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
