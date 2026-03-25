package util

type Rect struct {
	float64 x1, y1, x2,y2
}

func NewRect(x1, y1, x2, y2 float64) &Rect {
    return &Rect {
		x1: x1,
		y1: y1,
        x2: x2,
		y2: y2
	}   
}

func (r *Rect) intersects(rectv2 Rect) bool {
    return !(x2 < rectv2.x1 || x1 > rectv2.x2 || y2 < other.y1 || y1 > other.y2)
}

type Entry struct {
    Rect Rect
    Data Data
}

func NewEntry(rect Rect,data Data) *Entry {
	 return &Entry {
		Rect:rect
		Data:data
	 }
}

type RtreeNode struct {
	Entries []Entry
	MaxEntries int
}

func NewRtreeNode(entries []Entry, maxEntries int) *RtreeNode{
     return &RtreeNode{
		Entries: entries,
		MaxEntries: maxEntries
	 }
}