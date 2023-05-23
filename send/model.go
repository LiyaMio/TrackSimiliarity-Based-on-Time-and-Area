package send

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
type Node struct {
	X float64 `json:"lng"`
	Y float64 `json:"lat"`
}
type Trajectory struct {
	Id  int     `json:"id"`
	Tid int     `json:"tid"`
	X   float64 `json:"x"`
	Y   float64 `json:"y"`
}
type Link struct {
	Id        int    `json:"id"`
	Uid       string `json:"uid"`
	Tid       int    `json:"tid"`
	Starttime string `json:"starttime"`
	Endtime   string `json:"endtime"`
}
type Driver struct {
	Uid   string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Node []Node `json:"path"`
}
type Date struct {
	Stime string `json:"stime"`
	Etime string `json:"etime"`
	Node  []Node `json:"path"`
	E string `json:"距离阈值"`
	T string `json:"时间阈值"`
}
