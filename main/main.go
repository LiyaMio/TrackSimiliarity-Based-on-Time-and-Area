package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"track/compute"
	"track/dbs"
	"track/send"
	"track/user"
)

func main() {
	//encryption.Test_homomorphicCrypto()
	ans := []string{}
	tid := []string{}
	r := gin.Default()
	r.LoadHTMLGlob(
		"views/HTML/*",
	)
	//r.StaticFS("/views", http.Dir("/views"))
	r.StaticFS("/JS", http.Dir("views/JS"))

	r.StaticFS("/CSS", http.Dir("views/CSS"))

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/map")
	})
	r.GET("/map", func(c *gin.Context) {
		c.HTML(http.StatusOK, "map.html", nil)
	})
	r.POST("/test", func(c *gin.Context) {
		var db dbs.Database
		var err error
		db.MysqlDb, err= sql.Open("mysql", "root:09070810.@tcp(localhost:3306)/track")
		if err != nil {
			fmt.Printf("connect to db 127.0.0.1:3306 error: %v\n", err)
		}
		var track =make(map[int][][]string)
		sqlStr:="select * from trajectory"
		rows,ok:=db.QueryDB(sqlStr)
		if ok{
			fmt.Printf("len(rows) %d\n",len(rows))
			for i:=0;i<len(tid);i++{
				var location [][]string
				for j:=0;j<len(rows);j++{
					if(tid[i]==rows[j]["tid"] ){
						var site []string
						site= append(site, rows[j]["x"])
						site= append(site, rows[j]["y"])
						location= append(location,site)
						//fmt.Println("local %v\n",location)
					}
				}
				//fmt.Printf("locar length %d,%d",i,len(location))
				track[i]= location
			}
		}else{
			c.JSON(200,gin.H{
				"msg":"相似轨迹查询失败",
			})
		}
		c.JSON(200,gin.H{
			"track":track,
		})
		//c.JSON(302, gin.H{"location": "/sate3"})

	})
	r.POST("/map", func(c *gin.Context) {
		var now send.Date
		err := c.BindJSON(&now)
		fmt.Println(err)
		log.Printf("/map post error %v",err)
		//fmt.Println(len(now.Node))
		fmt.Printf("/map post now %v",now)
		ans,tid = compute.Count(now)

		var db dbs.Database
		db.MysqlDb, err= sql.Open("mysql", "root:09070810.@tcp(localhost:3306)/track")
		if err != nil {
			fmt.Printf("connect to db 127.0.0.1:3306 error: %v\n", err)
		}
		var track =make(map[string][][]string)
		sqlStr:="select * from trajectory"
		rows,ok:=db.QueryDB(sqlStr)
		if ok{
			fmt.Printf("len(rows) %d\n",len(rows))
			for i:=0;i<len(tid);i++{
				var location [][]string
				for j:=0;j<len(rows);j++{
					if(tid[i]==rows[j]["tid"] ){
						var site []string
						site= append(site, rows[j]["x"])
						site= append(site, rows[j]["y"])
						location= append(location,site)
					}
				}
				track["N"+strconv.Itoa(i)]= location
			}
		}else{
			c.JSON(200,gin.H{
				"msg":"相似轨迹查询失败",
			})
		}
		fmt.Println(track)
		c.JSON(200,gin.H{
			"track":track,
			"len":len(tid),
		})
	})
	r.GET("/json", func(c *gin.Context) {
		var db dbs.Database
		var err error
		db.MysqlDb, err = sql.Open("mysql", "root:09070810.@tcp(localhost:3306)/track")

		if err != nil {
			fmt.Printf("connect to db 127.0.0.1:3306 error: %v\n", err)
		}

		length := len(ans)
		var d []send.Driver
		sqlStr:="select * from driver"
		res,ok:=db.QueryDB(sqlStr)
		if ok{
			for i:=0;i<length;i++{
				for j:=0;j<len(res);j++{
					if ans[i]==res[j]["id"]{
						log.Println(res[j])
						var tmp send.Driver
						tmp.Uid=ans[i]
						tmp.Name=res[j]["name"]
						tmp.Phone=res[j]["phone"]
						log.Println(tmp)
						d= append(d, tmp)
					}
				}
			}
		}else{
			log.Println("查询失败")
			c.JSON(200,gin.H{
				"msg":"无相似轨迹1",
			})
		}
		c.JSON(200,gin.H{
			"data":d,
		})


	})
	r.GET("/sate3", func(c *gin.Context) {
		//c.JSON(200, gin.H{"data": "1"})
		c.HTML(http.StatusOK, "state_green.html", nil)

	})
	r.POST("/register", user.Register)
	r.POST("/login",user.Login)
	r.POST("/cookie",user.SetCookie)
	r.POST("/addtrack",user.AddTrack)
	r.POST("/postTrack",user.PostTrack)
	r.POST("/choosetrack",user.ChooseTrack)
	r.Run(":80")
}
