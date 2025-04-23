package pkg

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
	"github.com/pixiv/go-libjpeg/jpeg"
)

func ImageEncode(w io.Writer, m image.Image, im_type int) error {
	jpeg_options := &jpeg.EncoderOptions{
		Quality:         85,
		ProgressiveMode: true,
	}

	webp_options := &webp.Options{Lossless: false, Quality: 85}

	if im_type == 1 {
		return webp.Encode(w, m, webp_options)
	} else if im_type == 2 {
		return jpeg.Encode(w, m, jpeg_options)
	}
	return nil
}

func Resize(width uint, height uint, img image.Image) image.Image {
	return resize.Resize(width, height, img, resize.Lanczos2)
}
