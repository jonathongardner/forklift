package helpers

func MatchSig(raw []byte, toMatch []byte, offset int) bool {
	if len(raw) < offset+len(toMatch) {
		return false
	}

	for i := 0; i < len(toMatch); i++ {
		if raw[offset+i] != toMatch[i] {
			return false
		}
	}

	return true
}

func MatchSigFunc(toMatch []byte, offset int) func(raw []byte, limit uint32) bool {
	return func(raw []byte, limit uint32) bool {
		return MatchSig(raw, toMatch, offset)
	}
}

func MatchSigMultiOffsetFunc(toMatch []byte, offsets []int) func(raw []byte, limit uint32) bool {
	return func(raw []byte, limit uint32) bool {
		for i := 0; i < len(offsets); i++ {
			if MatchSig(raw, toMatch, offsets[i]) {
				return true
			}
		}
		return false
	}
}
