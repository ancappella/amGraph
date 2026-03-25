package util


type HP struct{
	Global []int
	sort.IntSlice
}

func NewHp(nums []int) *hp {
   return &hp{
	 Global: nums,
   }
}

func (h HP) Less(i, j int) bool {
	return h.Global[h.IntSlice[i]] > h.Global[h.IntSlice[j]]
}

func (h *HP) Push(v interface{}) {
	h.IntSlice = append(h.IntSlice, v.(int))
}

func (h *HP) Pop() interface{} {
	a := h.IntSlice
	v := a[len(a)-1]
	h.IntSlice = a[:len(a) - 1]
	return v
}


