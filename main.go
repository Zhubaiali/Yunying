package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	http.HandleFunc("/cut-video", cutVideoHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
	fmt.Println("hello world") 
}


func cutVideoHandler(w http.ResponseWriter, r *http.Request) {
	// 解析请求参数
	videoURL := r.URL.Query().Get("url")
	startTime := r.URL.Query().Get("start")
	endTime := r.URL.Query().Get("end")

	// 下载视频
	videoFile, err := downloadVideo(videoURL)
	if err != nil {
		log.Printf("下载视频失败: %v\n", err)
		http.Error(w, "下载视频失败", http.StatusInternalServerError)
		return
	}
	defer os.Remove(videoFile.Name())

	// 剪辑视频
	outputFile, err := cutVideo(videoFile.Name(), startTime, endTime)
	if err != nil {
		log.Printf("剪辑视频失败: %v\n", err)
		http.Error(w, "剪辑视频失败", http.StatusInternalServerError)
		return
	}

	// 返回剪辑后的视频URL
	fmt.Fprintf(w, "http://%s/%s", r.Host, outputFile)
}

func downloadVideo(url string) (*os.File, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	videoFile, err := os.CreateTemp("", "input-*.mp4")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(videoFile, resp.Body)
	return videoFile, err
}

func cutVideo(inputFile, startTime, endTime string) (string, error) {
	outputFile := fmt.Sprintf("output-%s-%s.mp4", startTime, endTime)

	cmd := exec.Command("ffmpeg", "-i", inputFile, "-ss", startTime, "-to", endTime, "-c", "copy", outputFile)
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return outputFile, nil
}
