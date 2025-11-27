package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
)

type ImageColor struct {
	Filename string
	R, G, B uint8
	Hue float64
}

func main() {
	
	directory := "./images"
	newDir := "./sorted_images"
	colorDir := "./colorBlocks"
	
	entries, err := os.ReadDir(directory)
	if err != nil{
		log.Fatal(err)
	}

	var images []ImageColor
	for _, entry := range entries {

		if entry.IsDir(){
			continue
		}

		arr :=[]int{0,0,0,0}

		inputFile := filepath.Join(directory, entry.Name())
		file, err := os.Open(inputFile)
		if err != nil {
			log.Fatal(err)
		}
		img, _, err := image.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()

		bounds := img.Bounds()
		width := bounds.Dx()
		height:= bounds.Dy()
		totalPixels := width*height
		
		for i:=0; i < width; i++{
			for j:=0; j < height; j++{

				pixelColor := img.At(i,j)

				r,g,b,_ := pixelColor.RGBA()

				red := int(r >> 8)
				green := int(g >> 8)
				blue := int(b>>8)
				arr[0] += int(red)
				arr[1] += int(green)
				arr[2] += int(blue)
			}
		}
					
		r,g,b := uint8(arr[0] / totalPixels), uint8(arr[1] / totalPixels),uint8(arr[2]/totalPixels)
		hue := RGBtoHue(r,g,b)
		

		images = append(images, ImageColor{
			Filename: entry.Name(),
			R: r,
			G: g,
			B: b,
			Hue: hue,
		})
	}

		sort.Slice(images, func(i,j int) bool {
			return images[i].Hue < images[j].Hue
		})

		for i, imgInfo := range images {

		oldPath := filepath.Join(directory, imgInfo.Filename)
		file, _ := os.Open(oldPath)
		imgDecoded, _, _ := image.Decode(file)
		file.Close()

		newImageName := fmt.Sprintf("%04d.jpg", i+1)
		newImagePath := filepath.Join(newDir, newImageName)

		outFile, err := os.Create(newImagePath)
		if err != nil {
			log.Fatal(err)
		}

		jpeg.Encode(outFile, imgDecoded, &jpeg.Options{Quality: 100})
		outFile.Close()

		avgColor := color.RGBA{R: imgInfo.R, G: imgInfo.G, B: imgInfo.B, A: 255}
		block := image.NewRGBA(image.Rect(0, 0, 100, 100))

		for y := 0; y < 100; y++ {
			for x := 0; x < 100; x++ {
				block.Set(x, y, avgColor)
			}
		}

		blockName := fmt.Sprintf("%04d.jpg", i+1)
		blockPath := filepath.Join(colorDir, blockName)

		blockFile, err := os.Create(blockPath)
		if err != nil {
			log.Fatal(err)
		}

		jpeg.Encode(blockFile, block, &jpeg.Options{Quality: 100})
		blockFile.Close()
	}
}


	func RGBtoHue(r, g, b uint8) float64 {
		rf := float64(r)/255
		gf := float64(g)/255
		bf := float64(b)/255

		max := math.Max(rf, math.Max(gf, bf))
		min := math.Min(rf, math.Min(gf, bf))
		delta := max - min

		
		var h float64
		if delta == 0 {
			h = 0
		} else if max == rf {
			h = math.Mod(((gf-bf)/delta), 6)
		} else if max == gf {
			h = ((bf-rf)/delta) + 2
		} else {
			h = ((rf-gf)/delta) + 4
		}

		h *= 60
		if h < 0 {
			h += 360
		}
		return h
	}