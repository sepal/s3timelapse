package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func GenerateTimelapse(files string, video string, speed float64) {
	speed_param := fmt.Sprintf("%f*PTS", speed)

	ffmpeg.Input(files, ffmpeg.KwArgs{"pattern_type": "glob"}).
		Filter("setpts", ffmpeg.Args{speed_param}).
		Output(video, ffmpeg.KwArgs{"c:v": "libx264", "pix_fmt": "yuv420p", "framerate": 30}).
		OverWriteOutput().Run()
}

func ParseUrl(url string) (bucket string, prefix string) {
	parts := strings.Split(url, "/")
	bucket = parts[2]

	prefix = strings.Join(parts[3:], "/")

	return bucket, prefix
}

func ListObjects(bucket string, prefix string, session *session.Session) ([]*s3.Object, error) {
	svc := s3.New(session)
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket), Prefix: aws.String(prefix)})

	if err != nil {
		return nil, err
	}

	return resp.Contents, nil
}

func DownloadImages(session *session.Session, bucket string, objects []*s3.Object) error {

	for _, item := range objects {
		key_parts := strings.Split(*item.Key, "/")
		fn := key_parts[len(key_parts)-1]

		file, err := os.Create(filepath.Join("./images", fn))
		defer file.Close()
		if err != nil {
			return err
		}

		downloader := s3manager.NewDownloader(session)
		numBytes, err := downloader.Download(file, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    item.Key,
		})

		if err != nil {
			return err
		}

		log.Println("Downloaded", file.Name(), numBytes, "bytes")
	}
	return nil
}

func TimeInDay(t time.Time, day time.Time) bool {
	durD, _ := time.ParseDuration("24h")
	end := day.Add(durD)

	if (t.Equal(day) || t.After(day)) && t.Before(end) {
		return true
	} else {
		return false
	}
}

func main() {
	var url string
	var out string
	var speed float64
	var forDay string
	flag.StringVar(&url, "url", "", "An s3 url containing the timelapse images, e.g.: s3://mybucket/images/")
	flag.StringVar(&out, "output", "out.mp4", "The filename of the timelapse video.")
	flag.Float64Var(&speed, "speed", 1.0, "The speed of the timelapse in PTS.")
	flag.StringVar(&forDay, "for", "", "Generate a timelapse for a certain day. The s3 last modified date will be used to pull the relevant images.")

	info, err := os.Stat("./images")

	if os.IsNotExist(err) {
		os.Mkdir("./images", 0750)
	} else if !info.IsDir() {
		log.Fatal("'images' file found, can't create folder")
		return
	}
	flag.Parse()

	bucket, prefix := ParseUrl(url)

	session, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")})

	objects, err := ListObjects(bucket, prefix, session)

	if forDay != "" {
		day, err := time.Parse("2006-01-02", forDay)

		if err != nil {
			log.Fatal(err)
			return
		}

		n := 0
		for _, item := range objects {
			if TimeInDay(*item.LastModified, day) {
				objects[n] = item
				n++
			}
		}

		objects = objects[:n]
	}

	if err != nil {
		log.Fatal(err)
		return
	}

	err = DownloadImages(session, bucket, objects)

	if err != nil {
		log.Fatal(err)
		return
	}

	GenerateTimelapse("images/*.jpg", out, speed)

}
