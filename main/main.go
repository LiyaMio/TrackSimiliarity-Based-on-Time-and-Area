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
		e,_:=strconv.ParseFloat(now.E,64)
		t,_:=strconv.ParseFloat(now.T,64)
		log.Printf("/map post error %v\n",err)
		fmt.Printf("/map post now %v\n",now)
		fmt.Printf("阈值：%s \n",now.E)
		ans,tid = compute.Count(now,e,t)

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

	r.POST("/register", user.Register)
	r.POST("/login",user.Login)
	r.POST("/cookie",user.SetCookie)
	r.POST("/addtrack",user.AddTrack)
	r.POST("/postTrack",user.PostTrack)
	r.POST("/choosetrack",user.ChooseTrack)
	r.POST("/deletetrack",user.DeleteTrack)
	r.POST("/contrastcommon",user.ChooseTrack)
	r.POST("/contrastpath",func(c *gin.Context) {
		var now send.Date
		err := c.BindJSON(&now)
		e, _ := strconv.ParseFloat(now.E, 64)
		t, _ := strconv.ParseFloat(now.T, 64)
		log.Printf("/map post error %v\n", err)
		fmt.Printf("/map post now %v\n", now)
		ans, tid = compute.Count(now, e, t)

		var db dbs.Database
		db.MysqlDb, err = sql.Open("mysql", "root:09070810.@tcp(localhost:3306)/track")
		if err != nil {
			fmt.Printf("connect to db 127.0.0.1:3306 error: %v\n", err)
		}
		var track = make(map[string][][]string)
		sqlStr := "select * from trajectory"
		rows, ok := db.QueryDB(sqlStr)
		if ok {
			fmt.Printf("len(rows) %d\n", len(rows))
			for i := 0; i < len(tid); i++ {
				var location [][]string
				for j := 0; j < len(rows); j++ {
					if (tid[i] == rows[j]["tid"]) {
						var site []string
						site = append(site, rows[j]["x"])
						site = append(site, rows[j]["y"])
						location = append(location, site)
					}
				}
				track["N"+strconv.Itoa(i)] = location
			}
		} else {
			c.JSON(200, gin.H{
				"msg": "相似轨迹查询失败",
			})
		}
		c.JSON(200, gin.H{
			"track": track,
			"len":   len(tid),
		})
	})
	r.Run(":80")
}
