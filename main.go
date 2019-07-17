package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func imageDecode(format string, r io.Reader) (img image.Image, err error) {
	switch format[1:] {
	case "png":
		img, err = png.Decode(r)
	case "jpeg":
		img, err = jpeg.Decode(r)
	case "jpg":
		img, err = jpeg.Decode(r)
	default:
		img, err = image.Image(nil), errors.New("Warning: Cannot decode "+format+" format.")
	}
	return img, err
}

func createGif(path string, output string, delay int, loopcount int) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}

	var images []*image.Paletted
	var delays []int
	for i := range files {

		file, err := os.Open(path + files[i].Name())
		if err != nil {
			fmt.Println(err)
		}

		format := filepath.Ext(files[i].Name())
		img, err := imageDecode(format, file)
		if err != nil {
			fmt.Println(files[i].Name(), err)
			continue
		}
		fmt.Println(file.Name())

		buf := new(bytes.Buffer)
		err = gif.Encode(buf, img, nil)
		if err != nil {
			fmt.Println(err)
		}

		tempImg, err := gif.Decode(buf)
		images = append(images, tempImg.(*image.Paletted))
		delays = append(delays, delay)
	}

	out, err := os.Create(output)
	if err != nil {
		fmt.Println(err)
	}

	err = gif.EncodeAll(out, &gif.GIF{Image: images, Delay: delays, LoopCount: loopcount})
}

func main() {
	var delay int
	var loopcount int
	var output string
	var path string

	flag.IntVar(&delay, "d", 4, "Time delay between frames, in 100ths of a second.")
	flag.IntVar(&loopcount, "lc", 0, "Controls the number of times an animation will be restarted during display. Values: 0 - loops forever, -1 - shows each frame once, n - shows each frame n+1 times.")
	flag.StringVar(&output, "o", "out.gif", "Output file.")
	flag.StringVar(&path, "p", "frames/", "Relative path to sequence of images.")

	flag.Parse()
	createGif(path, output, delay, loopcount)
}
