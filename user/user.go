package user

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"strings"
	"track/compute"
	"track/send"
)

func Register(c*gin.Context){
	var u user
	c.ShouldBind(&u)
	sqlStr:="select  from user where name =?"
	var tmp string
	db.QueryRow(sqlStr,u.Name).Scan(&tmp)
	if tmp!=""{
		c.JSON(200,gin.H{
			"msg":"该用户名已存在",
		})
	}else {
		sqlStr = "insert into user(password, name) values(?,?)"
		_, err := db.Exec(sqlStr, u.Pwd, u.Name)
		if err != nil {
			c.JSON(200, gin.H{
				"msg": "插入失败",
			})
		} else {
			c.JSON(200, gin.H{
				"msg": "success",
			})
		}
	}
}
var u user
func Login(c*gin.Context){
	c.ShouldBind(&u)
	fmt.Println(u)
	sqlStr:="select password from user where name = ?"
	var tmp string
	db.QueryRow(sqlStr,u.Name).Scan(&tmp)
	if tmp==""{
		c.JSON(200, gin.H{
			"msg": "用户不存在",
		})
	}else {
		if tmp==u.Pwd{
			c.JSON(200, gin.H{
				"msg": "success",
			})
			c.SetCookie("cookie", u.Name, 3600, "/", "127.0.0.1", false, true)
			//c.JSON(200,"success")
		}else{
			c.JSON(200, gin.H{
				"msg": "密码不正确",
			})
		}
	}
}
func AddTrack(c*gin.Context){
	var t Track
	c.ShouldBind(&t)
	if t.Node==nil{
		return
	}
	fmt.Println(t)
	name,err := c.Cookie("cookie")
	if err != nil {
		log.Println(err.Error())
	} else {
		sqlStr:="select id from user where name =?"
		db.QueryRow(sqlStr,name).Scan(&t.Uid)
		sqlStr="select max(tid) from usertrack where uid=?"
		db.QueryRow(sqlStr,t.Uid).Scan(&t.Tid)
		if t.Tid==""{
			t.Tid="0"
		}else{
			tmp,_:=strconv.Atoi(t.Tid)
			t.Tid=strconv.Itoa(tmp+1)
		}
		fmt.Println(t.Node[0])
		fmt.Println(t.Node[1])
		InsertUserTrack(0,t)
		fmt.Println(t)
	}
}
func InsertUserTrack(i int,t Track){
	sqlStr:="insert into usertrack (uid,tid,x,y,tname) values "
	len :=len(t.Node)
	if i+100<=len{
		len=i*100+100
	}
	fmt.Println(len)
	for j:=i*100;j<len;j++{
		fmt.Println(j)
		x := fmt.Sprintf("%.14f",t.Node[j].X)
		y:= fmt.Sprintf("%.14f",t.Node[j].Y)
		sqlStr+=" ('"+t.Uid+"','"+t.Tid+"','"+x+"','"+y+"','"+t.Tname+"'),"
		//fmt.Printf(sqlStr)
	}
	sqlStr = strings.TrimRight(sqlStr, ",")
	_,err:=db.Exec(sqlStr)
	if err!=nil{
		log.Println(err.Error())
	}
}
func SetCookie(c*gin.Context){
	fmt.Println(u.Name)
	c.SetCookie("cookie", u.Name, 3600, "/", "127.0.0.1", false, true)
	c.JSON(200,"success")
}
func PostTrack(c*gin.Context){
	name,err := c.Cookie("cookie")
	if err != nil {
		log.Println(err.Error())
	} else {
		var uid string
		sqlStr:="select id from user where name =?"
		db.QueryRow(sqlStr,name).Scan(&uid)
		sqlStr="select distinct (tid),tname from usertrack where uid="+uid
		var track []Track
		res,ok:=Db.QueryDB(sqlStr)
		if ok{
			for i:=0;i<len(res);i++{
				var tmp Track
				tmp.Tid=res[i]["tid"]
				tmp.Tname=res[i]["tname"]
				track= append(track, tmp)
			}
			c.JSON(200,track)
		}

	}

}
func ChooseTrack(c*gin.Context){
	d,_:=c.GetRawData()
	var body map[string]string
	_=json.Unmarshal(d,&body)
	var data send.Date
	data.Stime=body["stime"]
	data.Etime=body["etime"]
	name:=body["name"]
	uname,err:= c.Cookie("cookie")
	tmp:=body["距离阈值"]
	e,_:=strconv.ParseFloat(tmp,64)
	tmp=body["时间阈值"]
	fmt.Printf("时间阈值%s\n",tmp)
	t,_:=strconv.ParseFloat(tmp,64)
	if err!=nil{
		c.JSON(200,gin.H{
			"msg":"登录已经过期，重新登录",
		})
	}else{
		var uid string
		sqlStr:="select id from user where name =?"
		db.QueryRow(sqlStr,uname).Scan(&uid)
		sqlStr="select * from usertrack where uid="+uid+" and tname = '"+name+"'"
		res,ok :=Db.QueryDB(sqlStr)
		if ok {
			for i:=0;i<len(res);i++{
				var tmp send.Node
				x,_:= strconv.ParseFloat(res[i]["x"],64)
				y,_:= strconv.ParseFloat(res[i]["y"],64)
				tmp.X=x
				tmp.Y=y
				data.Node= append(data.Node, tmp)
			}
			_,tid := compute.Count(data,e,t)
			fmt.Println(tid)
			var track =make(map[string][][]string)
			sqlStr="select * from trajectory"
			rows,ok:=Db.QueryDB(sqlStr)
			if ok{
				fmt.Printf("len(rows) %d\n",len(rows))
				for i:=0;i<len(tid);i++{
					var location [][]string
					for j:=0;j<len(rows);j++ {
						if (tid[i] == rows[j]["tid"]) {
							var site []string
							site = append(site, rows[j]["x"])
							site = append(site, rows[j]["y"])
							location = append(location, site)
						}
					}
					track["N"+strconv.Itoa(i)]= location
				}
			}else{
				c.JSON(200,gin.H{
					"msg":"相似轨迹查询失败",
				})
			}
			c.JSON(200,gin.H{
				"orgin":data.Node,
				"track":track,
				"len":len(tid),
			})
		}
		}

}
func DeleteTrack(c*gin.Context){
	d,_:=c.GetRawData()
	var body map[string]string
	_=json.Unmarshal(d,&body)
	name:=body["name"]
	uname,err:= c.Cookie("cookie")
	if err!=nil{
		c.JSON(200,gin.H{
			"msg":"登录已经过期，重新登录",
		})
	}else{
		var uid string
		sqlStr:="select id from user where name =?"
		db.QueryRow(sqlStr,uname).Scan(&uid)
		sqlStr="delete from usertrack where uid=? and tname=?"
		_,err=db.Exec(sqlStr,uid,name)
		if err!=nil {
			c.JSON(200, gin.H{
				"msg": "删除失败",
			})
		}else{
			c.JSON(200,gin.H{
				"msg":"删除成功",
			})
		}
	}
}