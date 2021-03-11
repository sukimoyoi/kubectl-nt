package resourcerelationer

type ResourceQueue []ResourceRelationer

func (r ResourceQueue) Len() int { return len(r) }

func (r *ResourceQueue) Push(r2 ResourceRelationer) { *r = append(*r, r2) }

func (r *ResourceQueue) PushSlice(rList []ResourceRelationer) { *r = append(*r, rList...) }

func (r *ResourceQueue) Pop() ResourceRelationer {
	old := *r
	n := old.Len()
	x := old[0]
	*r = old[1:n]
	return x
}
