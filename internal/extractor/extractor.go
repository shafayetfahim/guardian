package extractor

import (
	"os"
	"github.com/rwcarlsen/goexif/exif"
)

// ExtractPhotoData opens a file and pulls the 'Holy Grail' settings
func ExtractPhotoData(path string) (map[string]interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		return nil, err
	}

	// Helper to extract specific tags safely
	data := make(map[string]interface{})
	
	if cam, err := x.Get(exif.Model); err == nil {
		data["camera"] = cam.String()
	}
	if fstop, err := x.Get(exif.FNumber); err == nil {
		data["aperture"] = fstop.String()
	}
	if iso, err := x.Get(exif.ISOSpeedRatings); err == nil {
		data["iso"] = iso.String()
	}
	if lens, err := x.Get(exif.LensModel); err == nil {
		data["lens"] = lens.String()
	}

	return data, nil
}