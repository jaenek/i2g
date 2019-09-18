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
	i2g(path, outputname, delay, loopcount)
}

func decode(img *os.File) (*image.Paletted, error) {
	ext := filepath.Ext(img.Name())

	var i image.Image
	var err error
	switch ext[1:] {
	case "png":
		i, err = png.Decode(img)
	case "jpg":
		fallthrough
	case "jpeg":
		i, err = jpeg.Decode(img)
	}
	if err != nil {
		return &image.Paletted{}, err
	}

	b := new(bytes.Buffer)
	err = gif.Encode(b, i, nil)
	if err != nil {
		return &image.Paletted{}, err
	}

	frame, err := gif.Decode(b)
	if err != nil {
		return &image.Paletted{}, err
	}

	return frame.(*image.Paletted), nil
}

func i2g(path string, outputname string, delay int, loopcount int) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var imgFrames []*os.File
	for _, fileInfo := range files {
		imgpath := path + fileInfo.Name()
		ext := filepath.Ext(imgpath)
		switch ext[1:] {
		case "png":
			fallthrough
		case "jpg":
			fallthrough
		case "jpeg":
			img, err := os.Open(imgpath)
			if err != nil {
				log.Fatal(err)
			}

			log.Println(fileInfo.Name() + ": Opened")
			imgFrames = append(imgFrames, img)
		default:
			log.Println(imgpath + ": Not Opened, Warning: Cannot decode " + ext + " format.")
		}
	}

	frames := make([]*image.Paletted, len(imgFrames))
	sem := make(chan error, len(files))
	for i, img := range imgFrames {
		go func(i int, img *os.File) {
			frames[i], err = decode(img)
			sem <- err
		}(i, img)
	}

	for _, img := range imgFrames {
		if <-sem != nil {
			log.Fatal(err)
		} else {
			log.Println(img.Name() + ": Decoded")
		}
	}

	out := gif.GIF{LoopCount: loopcount}
	for _, frame := range frames {
		out.Image = append(out.Image, frame)
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
