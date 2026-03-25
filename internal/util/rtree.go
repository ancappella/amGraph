package util


type Rect {
	int x1, y1;
	int x2, y2;
}

type Rtree struct {
    
}

func (r *Rtree) Rect(float64 x1, float64 y1, float64 x2, float64 y2) {
	this.x1 = x1
	this.y1 = y1
	this.x2 = x2
	this.y2 = y2
}

func intersects(Rect other) bool {
	return !(x2 < other.x1 || x1 > other.x2 || y2 < other.y2)
}


