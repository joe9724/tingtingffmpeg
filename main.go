package main

import "fmt"
import (
	"log"
	_"time"
	"flag"
	"os"
	_"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"os/exec"
	"strings"
	"bytes"
)

var (
	logFileName = flag.String("log", "/var/log/ffmpeg.log", "Log file name")
)

func exec_shell(s string) {
	cmd := exec.Command(s, "-i", "/root/go/src/resource/mp3/1522215547.m4a 2>&1 | grep 'Duration' | cut -d ' ' -f 4 | sed s/,// ")
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out.String())
}


func main() {
	logFile, logErr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("no log file access...exit...")
		//os.Exit(1)
	}
	log.SetOutput(logFile)
	log.Printf("start ")

	//exec_shell("/root/go/src/ffmpeg-git-20180314-64bit-static/ffmpeg -i /root/go/src/resource/mp3/1522215547.m4a 2>&1 | grep 'Duration' | cut -d ' ' -f 4 | sed s/,// ")
	exec_shell("/root/go/src/ffmpeg-git-20180314-64bit-static/ffmpeg")

	db,err := OpenConnection()
	if err!=nil{
		fmt.Println(err.Error())
	}
	defer db.Close()

	var chapters []Chapter
	db.Raw("select * from chapters").Find(&chapters)
	fmt.Println("name is",chapters[38].Name)
	///root/go/src/resource/mp3
	file := "/root/go/src/resource/mp3/"+strings.Replace(*(chapters[38].URL),"http://tingting-resource.bitekun.xin/resource/mp3/","",-1)
	fmt.Println("file is",file)
	c := exec.Command("/root/go/src/ffmpeg-git-20180314-64bit-static/ffmpeg","-i", file, "/root/go/src/resource/mp3/"+"test123.mp3")
	error:= c.Run()
	if error!=nil{
		fmt.Println("error: ",error.Error())
	}
	buf := new(bytes.Buffer)
	c.Stdout = buf
	fmt.Printf("%s", buf.String())




}



func OpenConnection() (db *gorm.DB, err error) {
	db, err = gorm.Open("mysql", "root:Tingtingyuedu654321!!!@tcp(47.104.131.147:3306)/tingting?charset=utf8&parseTime=True&loc=Local")
	return db,err
}

type Chapter struct {
	Id int64 `json:"id" gorm:"AUTO_INCREMENT"`
	// big cover
	// Required: true
	BigCover *string `json:"bigCover"`

	// duration
	// Required: true
	Duration *int64 `json:"duration"`

	// icon
	// Required: true
	Icon *string `json:"icon"`

	// name
	// Required: true
	Name *string `json:"name"`

	// order
	// Required: true
	Order *int64 `json:"order"`

	// play count
	// Required: true
	PlayCount *int64 `json:"playCount"`

	// show icon
	// Required: true
	ShowIcon *bool `json:"showIcon"`

	// status
	Status int64 `json:"status"`

	// sub title
	// Required: true
	SubTitle *string `json:"subTitle"`

	// time
	Time int64 `json:"time"`

	// update tips
	// Required: true
	UpdateTips *string `json:"updateTips"`

	// url
	// Required: true
	URL *string `json:"url"`

	Summary *string `json:"summary"`
}