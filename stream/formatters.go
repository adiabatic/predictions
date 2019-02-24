package stream

func deduplicateTags(ss []string) []string {
	seen := make(map[string]struct{}, len(ss))
	j := 0
	for _, v := range ss {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		ss[j] = v
		j++
	}
	return ss[:j]
}

// TagsUsed returns a list of all tags used in the given slice of stream.Stream.
func TagsUsed(ss []Stream) []string {
	ret := make([]string, 0)
	for _, s := range ss {
		for _, p := range s.Predictions {
			for _, t := range p.Tags {
				ret = append(ret, t)
			}
		}
	}
	return deduplicateTags(ret)
}
