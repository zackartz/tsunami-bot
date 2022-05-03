package lib

import (
	"encoding/json"
	"fmt"
	"github.com/fogleman/gg"
	"golang.org/x/image/draw"
	"image"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

var used = map[int]bool{
	12: true,
}

type BingoData struct {
	BingoText []string `json:"thing"`
}

type BingoCell struct {
	ID   string
	Text string
}

func generateBingoImage() {
	WIDTH := 2000
	HEIGHT := 2000

	var TOTAL_COLUMNS = 5
	var TOTAL_ROWS = 5
	var CEL_ID = 0
	var CEL_WIDTH = float64(WIDTH / TOTAL_COLUMNS)
	var CEL_HEIGHT = float64(HEIGHT / TOTAL_ROWS)

	dc := gg.NewContext(WIDTH, HEIGHT)
	dc.SetRGB(0, 0, 0)

	im, err := gg.LoadPNG("image.png")
	if err != nil {
		log.Panic(err)
	}

	resized := image.NewRGBA64(image.Rect(0, 0, im.Bounds().Max.X*2, im.Bounds().Max.Y*2))

	draw.NearestNeighbor.Scale(resized, resized.Rect, im, im.Bounds(), draw.Over, nil)

	dc.DrawImageAnchored(resized, WIDTH/2, HEIGHT/2, 0.5, 0.5)

	data := getJsonData(TOTAL_COLUMNS, TOTAL_ROWS)

	if err := dc.LoadFontFace("font.ttf", 34); err != nil {
		panic(err)
	}

	for col := 0; col < TOTAL_COLUMNS; col++ {
		for row := 0; row < TOTAL_ROWS; row++ {
			dc.DrawRectangle(float64(col)*CEL_WIDTH, float64(row)*CEL_HEIGHT, CEL_WIDTH, CEL_HEIGHT)
			dc.SetLineWidth(6)
			dc.SetRGBA(1, 1, 1, 0.45)
			dc.FillPreserve()
			dc.SetRGB(0, 0, 0)
			dc.Stroke()

			dc.DrawStringWrapped(
				data[CEL_ID].Text,
				float64(col)*CEL_WIDTH+CEL_WIDTH/2,
				float64(row)*CEL_HEIGHT+CEL_HEIGHT/2,
				0.5,
				0.5,
				260,
				1.3,
				gg.AlignCenter,
			)

			CEL_ID++
		}
	}

	err = dc.SavePNG("out.png")
	if err != nil {
		return
	}
}

func getJsonData(width, height int) []BingoCell {
	fileContent, err := os.Open("things.json")
	var ret []BingoCell

	if err != nil {
		log.Fatal(err)
		return nil
	}

	fmt.Println("The File is opened successfully...")

	defer fileContent.Close()

	byteResult, _ := ioutil.ReadAll(fileContent)

	var res BingoData
	json.Unmarshal(byteResult, &res)

	for i := 0; i < (width * height); i++ {
		var num int
		num = getRandomNumber(len(res.BingoText))

		ret = append(ret, BingoCell{
			ID:   fmt.Sprintf("1d%d", i),
			Text: res.BingoText[num],
		})
	}

	ret[12] = BingoCell{
		ID:   "1d13",
		Text: res.BingoText[12],
	}

	return ret
}

func getRandomNumber(total int) int {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(total)

	if used[num] == true {
		num = getRandomNumber(total)
	}

	used[num] = true

	return num
}
