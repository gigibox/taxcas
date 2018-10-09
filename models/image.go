package models

import (
	"bytes"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"math"
	"taxcas/pkg/setting"
	"taxcas/pkg/upload"

	"os"
	"taxcas/pkg/logging"
)

const (
	DefaultFontSize = 14
	DefaultDpi      = 72
)

type Signer struct {
	FontSize   float64
	Dpi        float64
	font       *truetype.Font
	signPoint  image.Point
	drawPoint  image.Point
}

func initFont(fontPath string) (*truetype.Font, error) {
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}
	return font, nil
}

func (this *Signer) SetFont(fontPath string, fontSize float64) (bool) {
	font, err := initFont(fontPath)
	if err != nil {
		return false
	}

	this.font = font

	if fontSize != 0 {
		this.FontSize = fontSize
	}

	return true
}

func (this *Signer) SetDrawPoint(x int, y int) {
	this.drawPoint = image.Pt(x, y)
}

func (this *Signer) SetSignPoint(x int, y int) {
	this.signPoint = image.Pt(x, y)
}

func (this *Signer) Sign(input io.Reader, output io.Writer, design *ImageDesigner) error {
	var (
		origin image.Image
		err    error
	)

	origin, err = png.Decode(input)
	if err != nil {
		logging.Warn("image decode error(%v)", err)
		return err
	}

	dst := image.NewNRGBA(origin.Bounds())
	draw.Draw(dst, dst.Bounds(), origin, image.ZP, draw.Src)

	coords := []Coord{design.Name, design.EnglishName, design.PersonalID, design.SerialNumber, design.Date}
	for i, _ := range coords {
		if coords[i].Str != "" && coords[i].X != 0 && coords[i].Y != 0 {
			// 设置字体, 字号
			fontPath := upload.GetFontPath() + coords[i].Font
			if _, err := os.Stat(fontPath); err != nil {
				fontPath = upload.GetFontPath() + "default.ttc"
			}

			this.SetFont(fontPath, coords[i].FontSize)

			mask := this.drawStringImage(coords[i].Str, coords[i].TextAlign)

			// 设置绘图位置
			X := coords[i].X
			if coords[i].TextAlign == "center" {
				X = (dst.Bounds().Dx() - this.drawPoint.X) / 2
			}
			this.SetSignPoint(X, coords[i].Y)

			draw.Draw(dst, mask.Bounds().Add(this.signPoint), mask, image.ZP, draw.Over)
		}
	}

	err = png.Encode(output, dst)
	if err != nil {
		logging.Warn("image encode error(%v)", err)
		return err
	}

	return nil
}

func (this *Signer) drawStringImage(text, align string) (image.Image) {
	rgba := image.NewRGBA(image.Rect(0, 0, this.drawPoint.X, this.drawPoint.Y))

	draw.Draw(rgba, rgba.Bounds(), image.Transparent, image.ZP, draw.Src)
	painter := &font.Drawer{
		Dst: rgba,
		Src: image.Black,
		Face: truetype.NewFace(this.font, &truetype.Options{
			Size:    this.FontSize,
			DPI:     this.Dpi,
			Hinting: font.HintingNone,
		}),
	}
	y := 10 + int(math.Ceil(this.FontSize * this.Dpi / 72))

	if align == "center" {
		painter.Dot = fixed.Point26_6{
			X: (fixed.I(this.drawPoint.X) - painter.MeasureString(text)) / 2,
			Y: fixed.I(y),
		}
	} else {
		painter.Dot = fixed.P(10, y)
	}

	painter.DrawString(text)

	return rgba
}

func SignImage(imagePath string,  design *ImageDesigner) (error) {
	byteBuff, err := ioutil.ReadFile(setting.AppSetting.RuntimeRootPath + design.ImgName)
	if err != nil{
		logging.Warn(err)
		return err
	}

	srcImage := bytes.NewReader(byteBuff)

	saveImage, err := os.OpenFile(setting.AppSetting.RuntimeRootPath + imagePath,  os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		logging.Warn(err)
		return err
	}
	defer saveImage.Close()

	signWriter :=  &Signer{
		FontSize:   DefaultFontSize,
		Dpi:        DefaultDpi,
		signPoint:  image.ZP,
		drawPoint:  image.Point{X: 500, Y: 500},
	}

	return signWriter.Sign(srcImage, saveImage, design)
}