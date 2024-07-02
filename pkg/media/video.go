package media

import (
	"os"
	"path"
	"time"
	"strings"
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
	var title, mimeType, album, comment string
	var format tag.Format
	var filetype tag.FileType
	var pic *tag.Picture
	m, err := tag.ReadFrom(f)
	if err != nil {
		idx := strings.LastIndex(name, ".") + 1
		if idx == 0 {
			return nil, err
		}
		ext := strings.ToUpper(name[idx:])
		if ext == "WAV" {
			title = name
			format = "WAVE"
			filetype = "WAV"
			mimeType = "audio/wav"
			album = ""
			comment = ""
		} else if ext == "WEBA" {
			title = name
			format = ""
			filetype = "WEBA"
			mimeType = "audio/webm"
			album = ""
			comment = ""
		} else if ext == "WEBM" {
			title = name
			format = ""
			filetype = "WEBM"
			mimeType = "video/webm"
			album = ""
			comment = ""
		} else {
			return nil, err
		}
	} else {
		title = m.Title()
		// Default title is filename
		if title == "" {
			title = name
		}
		format = m.Format()
		filetype = m.FileType()
		switch filetype {
			case "MP3":
				mimeType = "audio/mpeg"
			case "M4A", "M4B", "M4P":
				mimeType = "audio/aac"
			case "FLAC":
				mimeType = "audio/flac"
			case "OGG":
				mimeType = "audio/ogg"
			default:
				mimeType = ""
		}
		if format == "MP4" {
			mimeType = "video/mp4"
		}
		album = m.Album()
		comment = m.Comment()
		pic = m.Picture()
	}
	v := &Video{
		ID:          id,
		Title:       title,
		Album:       album,
		Description: comment,
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
	if pic != nil {
		v.Thumb = pic.Data
		v.ThumbType = pic.MIMEType
	}
	return v, nil
}
