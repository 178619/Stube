package media

import (
	"os"
	"path"
	"time"
	"github.com/dhowden/tag"
)

// Video represents metadata for a single video.
type Video struct {
	ID          string
	Title       string
	Album       string
	Description string
	Thumb       []byte
	ThumbType   string
	Modified    string
	Size        int64
	Path        string
	Timestamp   time.Time
	FileName	string
	Format		tag.Format
	FileType	tag.FileType
	MIMEType	string
}

// ParseVideo parses a video file's metadata and returns a Video.
func ParseVideo(p *Path, name string) (*Video, error) {
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
	m, err := tag.ReadFrom(f)
	if err != nil {
		return nil, err
	}
	title := m.Title()
	// Default title is filename
	if title == "" {
		title = name
	}
	format := m.Format()
	filetype := m.FileType()
	mimeType := ""
	switch filetype {
		case "MP3":
			mimeType = "audio/mpeg"
		case "M4A", "M4B", "M4P":
			mimeType = "audio/aac"
		case "FLAC":
			mimeType = "audio/flac"
		case "OGG":
			mimeType = "audio/ogg"
	}
	if format == "MP4" {
		mimeType = "video/mp4"
	}
	v := &Video{
		ID:          id,
		Title:       title,
		Album:       m.Album(),
		Description: m.Comment(),
		Modified:    modified,
		Size:        size,
		Path:        pth,
		Timestamp:   timestamp,
		FileName:    name,
		Format:      format,
		FileType:    filetype,
		MIMEType:    mimeType,
	}
	// Add thumbnail (if exists)
	pic := m.Picture()
	if pic != nil {
		v.Thumb = pic.Data
		v.ThumbType = pic.MIMEType
	}
	return v, nil
}
