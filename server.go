package main

import (
	"encoding/json"
	"image/png"
	"io"
	"os"
  "fmt"
	"log"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"
)

const captureDir = "./captures"
const capturesToKeep = 2
const captureMismatchThreshold = 4000

type imageChangedResponse struct {
	Changed bool `json:"changed"`
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

func capture() {
	t := time.Now()
	unixTime := t.Unix()

	filename := fmt.Sprintf("%s/capture_%d.png", captureDir, unixTime)
	captureCmd := fmt.Sprintf("screencapture -m -o -t png -x -a %s", filename)

	cmd := exec.Command("bash", "-c", captureCmd)
	err := cmd.Run()
	checkErr(err)
}

func cleanup() {
	files, err := ioutil.ReadDir(captureDir)
	checkErr(err)

	for i, file :=  range files {
		if len(files) - i > capturesToKeep { 
			err := os.Remove(captureFilePath(file))
			checkErr(err)
		}
	}
}

func didImageChange() (bool, error) {
	files, err := ioutil.ReadDir(captureDir)
	checkErr(err)

	readerA, err := os.Open(captureFilePath(files[len(files) - 1]))
	checkErr(err)
	readerB, err := os.Open(captureFilePath(files[len(files) - 2]))
	checkErr(err)

	imageA, err := png.Decode(readerA)
	checkErr(err)
	imageB, err := png.Decode(readerB)
	checkErr(err)

	mismatchPixelCount := 0
	
	rect := imageA.Bounds()
	for x := 0; x < rect.Size().X; x++ {
		for y := 0; y < rect.Size().Y; y++ {
			colorA := imageA.At(x,y)
			colorB := imageB.At(x,y)
			if colorA != colorB {
				mismatchPixelCount++
			}
		}
	}

	changed := mismatchPixelCount > captureMismatchThreshold
	return changed, nil
}

func captureFilePath(file os.FileInfo) string {
	return fmt.Sprintf("%s/%s", captureDir, file.Name())
}

func main() {

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile("./output.html")
		checkErr(err)

		fmt.Fprintf(w, string(data))
  })

	http.HandleFunc("/image_changed", func(w http.ResponseWriter, r *http.Request) {
		changed, err := didImageChange()
		checkErr(err)

		response := imageChangedResponse { Changed: changed }
		responseJSON, err := json.Marshal(response)
		checkErr(err)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(responseJSON))
	})

	http.HandleFunc("/capture.png", func(w http.ResponseWriter, r *http.Request) {
		files, err := ioutil.ReadDir(captureDir)
		checkErr(err)
		lastFile := files[len(files) - 1]
		img, err := os.Open(captureFilePath(lastFile))
		checkErr(err)

		defer img.Close()

		w.Header().Set("Content-Type", "image/png")
		io.Copy(w, img)
	})

	go func() {
		for true {
			capture()
			cleanup()
			time.Sleep(5 * time.Second)
		}
	}()
	
	log.Fatal(http.ListenAndServe(":4000", nil))
}