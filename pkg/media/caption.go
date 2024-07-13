package media

import (
	"os"
	"path"
	"time"
	"strings"
)

type Caption struct {
	ID          string
	Ref         string
	Modified    string
	Size        int64
	Path        string
	Timestamp   time.Time
	FileName	string
	FileType    string
	Origin		string
	SrcLang     string
}

func ParseCaption(p *Path, name string) (*Caption, error) {
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
	idx := strings.LastIndex(name, ".")
	idx = strings.LastIndex(name[:idx], ".")
	var ext string
	if idx == -1 {
		idx = strings.LastIndex(name, ".")
		ext = ""
	} else {
		ext = name[idx+1:strings.LastIndex(name, ".")]
	}
	v := &Caption{
		ID:         id,
		Ref:        ref,
		Modified:   modified,
		Size:       size,
		Path:       pth,
		Timestamp:  timestamp,
		FileName:   name,
		FileType:	"VTT",
		Origin:     name[:idx],
		SrcLang:    ext,
	}
	return v, nil
}