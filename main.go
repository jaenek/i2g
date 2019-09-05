package main

import (
	"bytes"
	"flag"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var path string
	var outputname string
	var delay int
	var loopcount int

	flag.StringVar(&path, "p", "frames/", "Relative path to sequence of images.")
	flag.StringVar(&outputname, "o", "out.gif", "Output file.")
	flag.IntVar(&delay, "d", 4, "Time delay between frames, in 100ths of a second.")
	flag.IntVar(&loopcount, "lc", 0, "Controls the number of times an animation will be restarted during display. Values: 0 - loops forever, -1 - shows each frame once, n - shows each frame n+1 times.")

	flag.Parse()
	run(path, outputname, delay, loopcount)
}

type GifFrame struct {
	Index   int
	Imgpath string
	Frame   *image.Paletted
}

func run(path string, outputname string, delay int, loopcount int) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var GifFrames []GifFrame
	for i, fileInfo := range files {
		imgpath := path + fileInfo.Name()
		ext := filepath.Ext(imgpath)
		switch ext[1:] {
		case "png":
			fallthrough
		case "jpg":
			fallthrough
		case "jpeg":
			log.Print(imgpath)
			GifFrames = append(GifFrames, GifFrame{Index: i, Imgpath: imgpath})
		default:
			log.Println(imgpath + ": Warning: Cannot decode " + ext + " format.")
		}
	}

	input := make(chan GifFrame, len(GifFrames))
	output := make(chan GifFrame, len(GifFrames))

	go worker(input, output)
	go worker(input, output)
	go worker(input, output)
	go worker(input, output)

	for _, img := range GifFrames {
		input <- img
	}

	for range GifFrames {
		img := <-output
		GifFrames[img.Index].Frame = img.Frame
	}

	out := gif.GIF{LoopCount: loopcount}
	for _, img := range GifFrames {
		out.Image = append(out.Image, img.Frame)
		out.Delay = append(out.Delay, delay)
	}

	outfile, err := os.Create(outputname)
	if err != nil {
		log.Println(err)
	}

	err = gif.EncodeAll(outfile, &out)
	if err != nil {
		log.Println(err)
	}
	outfile.Close()
}

func worker(input <-chan GifFrame, output chan<- GifFrame) {
	for img := range input {
		var i image.Image
		ext := filepath.Ext(img.Imgpath)

		file, err := os.Open(img.Imgpath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		switch ext[1:] {
		case "png":
			i, err = png.Decode(file)
		case "jpg":
			fallthrough
		case "jpeg":
			i, err = jpeg.Decode(file)
		}
		if err != nil {
			log.Fatal(err)
		}

		b := new(bytes.Buffer)
		err = gif.Encode(b, i, nil)
		if err != nil {
			log.Fatal(err)
		}

		frame, err := gif.Decode(b)
		if err != nil {
			log.Fatal(err)
		}

		img.Frame = frame.(*image.Paletted)
		output <- img
	}
}
