package media

import (
	"os"
	"path"
	"time"
	"strings"
	"github.com/dhowden/tag"
)

type MediaType string

const (
	AUDIO MediaType = "AUDIO"
	VIDEO MediaType = "VIDEO"
)

// Video represents metadata for a single video.
type Video struct {
	ID          string
	Ref         string
	Title       string
	Album       string
	Artist      string
	Track       int
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
	MediaType   MediaType
	BaseName    string
}

var replacer = strings.NewReplacer("#", "%23")

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
	ref := replacer.Replace(id)
	idx := strings.LastIndex(name, ".") + 1
	var ext string
	if idx == 0 {
		ext = name
	} else {
		ext = strings.ToUpper(name[idx:])
	}
	var track int = 0
	var title, mimeType, album, comment string
	var format tag.Format
	var filetype tag.FileType
	var pic *tag.Picture
	var artist string = ""
	var mediaType MediaType
	m, err := tag.ReadFrom(f)
	if err != nil {
		if ext == "WAV" {
			title = name
			format = "WAVE"
			filetype = "WAV"
			mimeType = "audio/wav"
			album = ""
			comment = ""
			mediaType = AUDIO
		} else if ext == "WEBA" {
			title = name
			format = ""
			filetype = "WEBA"
			mimeType = "audio/webm"
			album = ""
			comment = ""
			mediaType = AUDIO
		} else if ext == "WEBM" {
			title = name
			format = ""
			filetype = "WEBM"
			mimeType = "video/webm"
			album = ""
			comment = ""
			mediaType = VIDEO
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
		mediaType = AUDIO
		if format == "MP4" {
			mimeType = "video/mp4"
			mediaType = VIDEO
		}
		album = m.Album()
		artist = m.Artist()
		comment = m.Comment()
		pic = m.Picture()
		disc, _ := m.Disc()
		track, _ = m.Track()
		track += disc * 1000
	}
	v := &Video{
		ID:          id,
		Ref:         ref,
		Title:       title,
		Album:       album,
		Artist:      artist,
		Track:       track,
		Description: comment,
		Modified:    modified,
		Size:        size,
		Path:        pth,
		Timestamp:   timestamp,
		FileName:    name,
		Format:      format,
		FileType:    filetype,
		MIMEType:    mimeType,
		MediaType:   mediaType,
		BaseName:	 name[:idx-1],
	}
	// Add thumbnail (if exists)
	if pic != nil {
		v.Thumb = pic.Data
		v.ThumbType = pic.MIMEType
	}
	return v, nil
}
