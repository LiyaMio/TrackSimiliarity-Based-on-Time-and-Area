package user
type user struct {
	Name string `form:"name"`
	Pwd string `form:"pwd"`
}
type Node struct {
	X float64 `json:"lng"`
	Y float64 `json:"lat"`
}
type Track struct {
	Tname string `form:"tname"`
	Uid string `form:"uid"`
	Tid string  `form:"tid"`
	Node  []Node `json:"path"`
}

