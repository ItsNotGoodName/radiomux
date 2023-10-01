package diff

import "math"

type Changed int

const (
	ChangedNone = 0
	ChangedAll  = math.MaxInt
)

func (c Changed) Is(check Changed) bool {
	return c&check == check
}

func (c Changed) IsAll() bool {
	return c == ChangedAll
}

func (c Changed) Merge(merge Changed) Changed {
	return c | merge
}
