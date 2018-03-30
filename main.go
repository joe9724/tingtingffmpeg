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
	"github.com/vansante/go-ffprobe"
	"time"
	"encoding/json"
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
	//exec_shell("/root/go/src/ffmpeg-git-20180314-64bit-static/ffmpeg")

	db,err := OpenConnection()
	if err!=nil{
		fmt.Println(err.Error())
	}
	defer db.Close()

	ticker := time.NewTicker(600 * time.Second)
	for _= range ticker.C{
		fmt.Println(time.Now())
		var chapters []Chapter
		db.Raw("select * from chapters where duration=100000").Find(&chapters)
		//fmt.Println("name is",chapters[38].Name)
		///root/go/src/resource/mp3
		for k:=0;k<len(chapters) ;k++  {
			file := "/root/go/src/resource/mp3/"+strings.Replace(*(chapters[k].URL),"http://tingting-resource.bitekun.xin/resource/mp3/","",-1)
			fmt.Println("file is",file)
			path := file
			data, err := ffprobe.GetProbeData(path, 1000 * time.Millisecond)
			if err != nil {
				continue
				//log.Panicf("Error getting data: %v", err)
				fmt.Println("err3 is",err.Error())
			}

			//buf, err := json.MarshalIndent(data, "", "  ")
			//fmt.Print(string(buf))

			buf, err := json.MarshalIndent(data.GetFirstAudioStream(), "", "  ")
			if err != nil{
				continue
				fmt.Println("err1 is",err.Error())
			}
			log.Println(string(buf))

			fmt.Println("duration is",data.Format.DurationSeconds)
			db.Exec("update chapters set duration=? where id=?",data.Format.DurationSeconds,chapters[k].Id)
		}
	}



	//fmt.Printf("\nDuration: %v\n", data.Format.Duration())
	//fmt.Printf("\nStartTime: %v\n", data.Format.StartTime())
	/*c := exec.Command("/root/go/src/ffmpeg-git-20180314-64bit-static/ffprobe", file)
	//c := exec.Command("root/go/src/ffmpeg-git-20180314-64bit-static/ffmpeg", "-i", "/root/go/src/resource/mp3/1522215547.m4a" ,"2>&1", "|", "grep 'Duration' | cut -d ' ' -f 4 | sed s/,//")
	out,err := c.CombinedOutput()
	if err!= nil{
		fmt.Println("err is",err.Error())
	}
	log.Println("out is ",string(out))*/
    // 解析出Duration字段

	/*if error!=nil{
		fmt.Println("error: ",error.Error())
	}*/
	/*buf := new(bytes.Buffer)
	c.Stdout = buf
	fmt.Printf("%s", buf.String())*/





	/**
	ffmpeg version N-45325-gb173e0353-static https://johnvansickle.com/ffmpeg/  Copyright (c) 2000-2018 the FFmpeg developers
  built with gcc 6.3.0 (Debian 6.3.0-18+deb9u1) 20170516
  configuration: --enable-gpl --enable-version3 --enable-static --disable-debug --disable-ffplay --disable-indev=sndio --disable-outdev=sndio --cc=gcc-6 --enable-fontconfig --enable-frei0r --enable-gnutls --enable-gray --enable-libfribidi --enable-libass --enable-libfreetype --enable-libmp3lame --enable-libopencore-amrnb --enable-libopencore-amrwb --enable-libopenjpeg --enable-librubberband --enable-libsoxr --enable-libspeex --enable-libvorbis --enable-libopus --enable-libtheora --enable-libvidstab --enable-libvo-amrwbenc --enable-libvpx --enable-libwebp --enable-libx264 --enable-libx265 --enable-libxvid --enable-libzimg
  libavutil      56.  9.100 / 56.  9.100
  libavcodec     58. 14.100 / 58. 14.100
  libavformat    58. 10.100 / 58. 10.100
  libavdevice    58.  2.100 / 58.  2.100
  libavfilter     7. 13.100 /  7. 13.100
  libswscale      5.  0.102 /  5.  0.102
  libswresample   3.  0.101 /  3.  0.101
  libpostproc    55.  0.100 / 55.  0.100
Input #0, mov,mp4,m4a,3gp,3g2,mj2, from '/root/go/src/resource/mp3/1520209590.m4a':
  Metadata:
    major_brand     : M4A
    minor_version   : 512
    compatible_brands: isomiso2
    title           : 1ˀ½芏خ´𘮴󶅓µ±§
    artist          : [bbs.12fou.com]
    album           : µݗӹ嚢bs.12fou.com]
    encoder         : Lavf55.19.100
    comment         : 12·񷺍¯؊ԴÛ̳[bbs.12fou.com
  Duration: 00:04:55.01, start: 0.000000, bitrate: 97 kb/s
    Stream #0:0(und): Audio: aac (LC) (mp4a / 0x6134706D), 44100 Hz, stereo, fltp, 95 kb/s (default)
    Metadata:
      handler_name    : SoundHandler
At least one output file must be specified

	 */
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