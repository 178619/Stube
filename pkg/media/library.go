package media

import (
	"errors"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// Library manages importing and retrieving video data.
type Library struct {
	mu       sync.RWMutex
	Paths    map[string]*Path
	Videos   map[string]*Video
	Captions map[string]*Caption
}

// NewLibrary returns new instance of Library.
func NewLibrary() *Library {
	lib := &Library{
		Paths:    make(map[string]*Path),
		Videos:   make(map[string]*Video),
		Captions: make(map[string]*Caption),
	}
	return lib
}

// AddPath adds a media path to the library.
func (lib *Library) AddPath(p *Path) error {
	lib.mu.Lock()
	defer lib.mu.Unlock()
	// make sure new path doesn't collide with existing ones
	for _, p2 := range lib.Paths {
		if p.Path == p2.Path {
			return errors.New("media: duplicate library path")
		}
		if p.Prefix == p2.Prefix {
			return errors.New("media: duplicate library prefix")
		}
	}
	lib.Paths[p.Path] = p
	return nil
}

// Import adds all valid videos from a given path.
func (lib *Library) Import(p *Path) error {
	files, err := ioutil.ReadDir(p.Path)
	if err != nil {
		return err
	}
	for _, info := range files {
		if info.IsDir() {
			n := &Path{
				Path:   path.Join(p.Path, info.Name()),
				Prefix: path.Join(p.Prefix, info.Name()),
			}
			lib.AddPath(n)
			lib.Import(n)
			continue
		}
		err = lib.Add(path.Join(p.Path, info.Name()))
		if err != nil {
			// Ignore files that can't be parsed
			log.Println("Ignored:", path.Join(p.Path, info.Name()))
			continue
		}
	}
	return nil
}

// Add adds a single video from a given file path.
func (lib *Library) Add(fp string) error {
	lib.mu.Lock()
	defer lib.mu.Unlock()
	fp = filepath.ToSlash(fp)
	d := path.Dir(fp)
	p, ok := lib.Paths[d]
	if !ok {
		return errors.New("media: path not found")
	}
	n := path.Base(fp)
	if strings.HasSuffix(strings.ToUpper(n), ".VTT") {
		v, err := ParseCaption(p, n)
		if err != nil {
			return err
		}
		lib.Captions[v.ID] = v
		log.Println("Added:", v.Path)
		return nil
	}
	v, err := ParseVideo(p, n)
	if err != nil {
		return err
	}
	lib.Videos[v.ID] = v
	log.Println("Added:", v.Path)
	return nil
}

// Remove removes a single video from a given file path.
func (lib *Library) Remove(fp string) {
	lib.mu.Lock()
	defer lib.mu.Unlock()
	fp = filepath.ToSlash(fp)
	d := path.Dir(fp)
	p, ok := lib.Paths[d]
	if !ok {
		return
	}
	n := path.Base(fp)
	id := n
	if len(p.Prefix) > 0 {
		id = path.Join(p.Prefix, id)
	}
	if strings.HasSuffix(strings.ToUpper(n), ".VTT") {
		v, ok := lib.Captions[id]
		if ok {
			delete(lib.Captions, id)
			log.Println("Removed:", v.Path)
		}
	}
	v, ok := lib.Videos[id]
	if ok {
		delete(lib.Videos, id)
		log.Println("Removed:", v.Path)
	}
}

// Playlist returns a sorted Playlist of all videos.
func (lib *Library) Playlist() Playlist {
	lib.mu.RLock()
	defer lib.mu.RUnlock()
	pl := make(Playlist, len(lib.Videos))
	i := 0
	for _, v := range lib.Videos {
		pl[i] = v
		i++
	}
	sort.Sort(pl)
	return pl
}

// Playlist returns a sorted Playlist of all videos.
func (lib *Library) CaptionList() CaptionList {
	lib.mu.RLock()
	defer lib.mu.RUnlock()
	pl := make(CaptionList, len(lib.Captions))
	i := 0
	for _, v := range lib.Captions {
		pl[i] = v
		i++
	}
	return pl
}