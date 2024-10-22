package media

import (
	"os"
	"path"
	"strings"
	"time"
)

type Image struct {
	ID        string
	Ref       string
	Modified  string
	Size      int64
	Path      string
	Timestamp time.Time
	FileName  string
	FileType  string
	MIMEType  string
	Prefix    string
}

func ParseImage(p *Path, name string) (*Image, error) {
	pth := path.Join(p.Path, name)
	f, err := os.Open(pth)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := info.Size()
	timestamp := info.ModTime()
	modified := timestamp.Format("2006-01-02 03:04 PM")
	id := name
	if len(p.Prefix) > 0 {
		// if there's a prefix prepend it to the ID
		id = path.Join(p.Prefix, name)
	}
	ref := replacer.Replace(id)
	mimeType := "image"
	switch strings.ToUpper(name)[strings.LastIndex(name, ".")+1:] {
	case "PNG":
		mimeType = "image/png"
	case "BMP":
		mimeType = "image/bmp"
	case "GIF":
		mimeType = "image/gif"
	case "ICO":
		mimeType = "image/vnd.microsoft.icon"
	case "JPEG", "JPG":
		mimeType = "image/jpeg"
	case "SVG":
		mimeType = "image/svg+xml"
	case "TIF", "TIFF":
		mimeType = "image/tiff"
	case "WEBP":
		mimeType = "image/webp"
	}
	v := &Image{
		ID:        id,
		Ref:       ref,
		Modified:  modified,
		Size:      size,
		Path:      pth,
		Timestamp: timestamp,
		FileName:  name,
		MIMEType:  mimeType,
		Prefix:    p.Prefix,
	}
	return v, nil
}
