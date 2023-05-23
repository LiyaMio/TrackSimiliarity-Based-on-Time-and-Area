package compute

import (
	"database/sql"
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"
	"track/send"
)

type point struct {
	x1     [100]float64
	x2     [100]float64
	year   [100]int
	month  [100]int
	day    [100]int
	hour   [100]int
	minute [100]int
}

var m int
var patient point
var n int
var p point
var extrap point
var db *sql.DB
func init(){
	var err error
	db, err = sql.Open("mysql", "root:09070810.@tcp(localhost:3306)/track")
	if err != nil {
		fmt.Printf("connect to db 127.0.0.1:3306 error: %v\n", err)
	}
}
//经纬度距离计算
func EarthDistance(lat1, lng1, lat2, lng2 float64) float64 {
	radius := 6378.137
	rad := math.Pi / 180.0
	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad
	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	return dist * radius
}

//e取0.4，点距离比较
func distance(n int, p point,e float64) int {
	if n==0 {return -100}
	var m int
	//var (
	//	c, e float64
	//)
	var f int
	var c float64
	f = 0
	m = 90%n
	//e = 0.4
	arr:=90/n
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {

			c = EarthDistance(p.x1[i], p.x2[i], patient.x1[j], patient.x2[j])
			if c <= e {
				f = f + 1
			}
		}
		if i != n-1 {
			hanshu(p.x1[i], p.x1[i+1], p.x2[i], p.x2[i+1],arr)

			for k := 0; k < arr; k++ {
				for j := 0; j < m; j++ {
					c = EarthDistance(extrap.x1[k], extrap.x2[k], patient.x1[j], patient.x2[j])
					if c <= e {
						f = f + 1
					}
				}
			}
		}
	}
	return f
}

//函数取点
func hanshu(x1, x2, y1, y2 float64,arr int) {
	var (
		k, b, r float64
	)
	k = (y1 - y2) / (x1 - x2)
	b = y1 - k*x1
	r = (x2 - x1) / float64(arr)
	for i := 0; i < arr; i++ {
		x1 = x1 + r
		extrap.x1[i] = x1
		extrap.x2[i] = k*extrap.x1[i] + b
	}

}

//函数取时间点
func timehanshu(a1, a2, b1, b2,arr int) {
	var ho, mi int
	ho = a2 - a1
	mi = b2 - b1
	if ho > 0 {
		mi = mi + 60*ho
	}
	mi = mi / 6
	for i := 0; i < arr; i++ {
		b1 = b1 + mi
		if b1 >= 60 {
			b1 = b1 - 60
			a1++
		}
		extrap.hour[i] = a1
		extrap.minute[i] = b1
	}
}

//时间预判断
func timejudgef(n int, p point,t float64) int {
	var judge int
	judge = 0

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			if p.year[i] != patient.year[j] || p.month[i] != patient.month[j] || p.day[i] != patient.day[j]  {
				judge = -1
				break
			} else{
				hou := p.hour[0] - patient.hour[0]
				min := (p.minute[0] - patient.minute[0]) + hou*60
				if float64(min) < 0 {
					min = 0 - min
				}
				if float64(min) < t {
					judge=0

				}else{
					return -1
				}
				hou = p.hour[1] - patient.hour[1]
				min = (p.minute[1] - patient.minute[1]) + hou*60
				if min < 0 {
					min = 0 - min
				}
				if float64(min)< t {
					judge=0

				}else{
					return -1
				}
			}
		}
	}
	if judge == 0 {
		return 1
	} else {
		return -1
	}

}

////时间判断
func timejudge(n int, p point) int {

	var hou, min, f int
	arr:=90/n
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			hou = p.hour[i] - patient.hour[j]
			min = (p.minute[i] - patient.minute[j]) + hou*60
			if min < 0 {
				min = 0 - min
			}
			if min < 20 {
				f=1
			}else{
				f=-1
				break
			}
			fmt.Println(min)
		}
		if i == 0 {
			timehanshu(p.hour[i], p.hour[i+1], p.minute[i], p.minute[i+1],arr)
			for k := 0; k < arr; k++ {
				for j := 0; j < m; j++ {
					hou = extrap.hour[j] - patient.hour[j]
					min = (p.minute[i] - patient.minute[j]) + hou*60
					if min < 0 {
						min = 0 - min
					}
					if min < 20 {
						f=1
					}else {
						f=-1
						break
					}
				}
			}
		}

	}
	//	fmt.Printf("%d\n",f)
	return f
}
func parse_timestr_to_datetime(time_str string) time.Time {

	t, error4 := time.Parse("2006-01-02 15:04:05", time_str)
	if error4 != nil {
		panic(error4)
	}
	return t
}
func Count(now send.Date,e float64,t float64) ([]string,[]string) {
	fmt.Println("请输入您当日所到的地点数")
	n = len(now.Node)
	fmt.Println(n)
	fmt.Println("请依次输入您当日所到的地点的起止时间")
	st := parse_timestr_to_datetime(now.Stime)
	ed := parse_timestr_to_datetime(now.Etime)
	fmt.Println(int(st.Month()))
	fmt.Println(ed.Year())
	p.year[0] = st.Year()
	p.month[0] = int(st.Month())
	p.day[0] = st.Day()
	p.hour[0] = st.Hour()
	p.minute[0] = st.Minute()
	p.year[1] = ed.Year()
	p.month[1] = int(ed.Month())
	p.day[1] = ed.Day()
	p.hour[1] = ed.Hour()
	p.minute[1] = ed.Minute()

	fmt.Println("请依次输入您当日所到的地点")
	for i := 0; i < n; i++ {
		p.x1[i] = now.Node[i].X
		p.x2[i] = now.Node[i].Y
	}
	cnt := 1
	var dict = make(map[string]int, 10)
	var Tid []string
	for cnt <= 100 {
		sqlStr := "select trajectory.id, x, y,starttime,endtime from trajectory left join link l on trajectory.tid = l.tid where trajectory.tid = ?"
		rows, _:= db.Query(sqlStr, cnt)
		i := 0
		var stime,etime string
		defer rows.Close()
		for rows.Next() {
			var u send.Trajectory
			err := rows.Scan(&u.Id, &u.X, &u.Y,&stime,&etime)
			if err != nil {
				fmt.Printf("scan failed, err:%v\n", err)
			}
			patient.x1[i] = u.X
			patient.x2[i] = u.Y
			i++
		}
		st := parse_timestr_to_datetime(stime)
		ed := parse_timestr_to_datetime(etime)
		patient.year[0] = st.Year()
		patient.month[0] = int(st.Month())
		patient.day[0] = st.Day()
		patient.hour[0] = st.Hour()
		patient.minute[0] = st.Minute()
		patient.year[1] = ed.Year()
		patient.month[1] = int(ed.Month())
		patient.day[1] = ed.Day()
		patient.hour[1] = ed.Hour()
		patient.minute[1] = ed.Minute()
		var fen float64
		fen = 0
		fen = float64(timejudgef(n,p,t))
		if fen >= 0 {
			//fmt.Printf("%v %d\n",Tid,fen)
			fen = float64(distance(n, p,e))
			//fmt.Printf("fen is %v\n",fen)
			//fmt.Printf("tid is %d\n",cnt)
			//fmt.Printf("dis is %v",distance(n,p))
			//fmt.Printf("area is %v",area(n,p)/100)

			if fen != 0 {
				//dict[name] = fen
				fmt.Printf("轨迹序号 %v\n",cnt)
				Tid= append(Tid, strconv.Itoa(cnt))
			}


		}
		cnt++
	}
	type peroson struct {
		Name string
		Fun  int//相似度
	}

	var lstPerson []peroson
	for k, v := range dict {
		lstPerson = append(lstPerson, peroson{k, v})
	}

	sort.Slice(lstPerson, func(i, j int) bool {
		return lstPerson[i].Fun > lstPerson[j].Fun // 升序
	})
		var ans []string
		for _,v :=range  lstPerson{
			ans= append(ans, v.Name)
		}
		//fmt.Println(Tid)
	return ans,Tid
}
