package main

import (
	"fmt"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func generate_timelapse(files string, video string, speed float32) {
	speed_param := fmt.Sprintf("%f*PTS", speed)

	ffmpeg.Input(files, ffmpeg.KwArgs{"pattern_type": "glob"}).
		Filter("setpts", ffmpeg.Args{speed_param}).
		Output(video, ffmpeg.KwArgs{"c:v": "libx264", "pix_fmt": "yuv420p", "framerate": 30}).
		OverWriteOutput().Run()
}

func main() {
	generate_timelapse("basilpi/*.jpg", "test.mp4", 2.0)
}
