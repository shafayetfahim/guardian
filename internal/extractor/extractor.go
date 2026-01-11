package extractor

import (
	"os"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

type PhotoMetadata struct {
	Path     string `json:"path"`
	Camera   string `json:"camera"`
	Aperture string `json:"aperture"`
	ISO      string `json:"iso"`
	Lens     string `json:"lens"`
}

func ExtractPhotoData(path string) (*PhotoMetadata, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		return nil, err
	}

	data := &PhotoMetadata{
		Path: path,
	}

	clean := func(tag exif.FieldName) string {
		val, err := x.Get(tag)
		if err != nil {
			return "Unknown"
		}
		return strings.Trim(val.String(), "\"")
	}

	data.Camera = clean(exif.Model)
	data.Aperture = clean(exif.FNumber)
	data.ISO = clean(exif.ISOSpeedRatings)
	data.Lens = clean(exif.LensModel)

	return data, nil
}
