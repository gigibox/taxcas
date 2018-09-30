package models

import (
	"bytes"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"taxcas/pkg/setting"

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
	startPoint image.Point
	signPoint  image.Point
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

func (this *Signer) SetStartPoint(x int, y int) {
	this.startPoint = image.Pt(x, y)
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
			this.SetStartPoint(coords[i].X, coords[i].Y)

			// 设置字体, 字号
			fontPath := setting.AppSetting.RuntimeRootPath + setting.AppSetting.FontSavePath + coords[i].Font
			if _, err := os.Stat(fontPath); err != nil {
				fontPath = setting.AppSetting.RuntimeRootPath + setting.AppSetting.FontSavePath + "default.ttc"
			}

			this.SetFont(fontPath, coords[i].FontSize)

			mask, err := this.drawStringImage(coords[i].Str)
			if err != nil {
				logging.Warn("drawStringImage error(%v)", err)
				return err
			}
			draw.Draw(dst, mask.Bounds().Add(this.startPoint), mask, image.ZP, draw.Over)
		}
	}

	err = png.Encode(output, dst)
	if err != nil {
		logging.Warn("image encode error(%v)", err)
		return err
	}
	return nil
}

// draw text image
func (this *Signer) drawStringImage(text string) (image.Image, error) {
	fg, bg := image.Black, image.Transparent
	rgba := image.NewRGBA(image.Rect(0, 0, this.signPoint.X, this.signPoint.Y))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetDPI(this.Dpi)
	c.SetFont(this.font)
	c.SetFontSize(this.FontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)

	// Draw the text.
	pt := freetype.Pt(10, 10 + int(c.PointToFixed(this.FontSize)>>4))
	if _, err := c.DrawString(text, pt); err != nil {
		logging.Warn("c.DrawString(%s) error(%v)", text, err)
		return nil, err
	}

	return rgba, nil
}

func SignImage(imagePath string,  design *ImageDesigner) (error) {
	/*
	srcImage, err := os.Open(setting.AppSetting.RuntimeRootPath + design.ImgName)
	if err != nil {
		logging.Warn(err)
		return err
	}
	defer srcImage.Close()
	*/

	byteBuff, err := ioutil.ReadFile(setting.AppSetting.RuntimeRootPath + design.ImgName)
	if err != nil{
		fmt.Printf("%s\n",err)
		panic(err)
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
		startPoint: image.ZP,
		signPoint:  image.Point{X: 594, Y: 841}, // < A4 72dpi
	}

	if err := signWriter.Sign(srcImage, saveImage, design); err != nil {
		return err
	}

	return nil
}