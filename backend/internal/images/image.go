package images

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"

	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/disintegration/imaging"
)

// CompressedJPEG adjusts the resolution to fit within the given bounds and compresses the image so the file size is less than maxKB
func CompressedJPEG(imgByte []byte, maxWidth int, maxHeight int, maxKB int) ([]byte, error) {
	if maxWidth == 0 || maxHeight == 0 || maxKB == 0 {
		panic("must provide all params")
	}

	img, err := imaging.Decode(bytes.NewReader(imgByte), imaging.AutoOrientation(true))
	if err != nil {
		return nil, &errors.Error{Code: errors.Invalid, Message: "error processing image", Err: err}
	}
	fittedImg := imaging.Fit(img, maxWidth, maxHeight, imaging.Linear)

	quality := 75
	result, err := ToJPEG(fittedImg, quality)
	if err != nil {
		return nil, err
	}

	for len(result)/1000 > maxKB {
		if quality < 25 {
			return nil, &errors.Error{Code: errors.Invalid, Message: "image is too large, try using a smaller image"}
		}
		quality = int(float64(quality) * 0.75) // Reduce quality by 25% per step (e.g. 75, 56, 42, 31, 23)
		result, err = ToJPEG(fittedImg, quality)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func ToJPEG(img image.Image, quality int) ([]byte, error) {
	buf := bytes.Buffer{}
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, &errors.Error{Code: errors.Internal, Err: err}
	}

	result, err := io.ReadAll(&buf)
	if err != nil {
		return nil, &errors.Error{Code: errors.Internal, Err: err}
	}
	return result, nil
}
