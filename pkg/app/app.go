// Package app manages main application server.
package app

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"path"
	"strings"

	"github.com/178619/stube/pkg/media"
	"github.com/178619/stube/pkg/onionkey"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
)

// App represents main application.
type App struct {
	Config    *Config
	Library   *media.Library
	Watcher   *fsnotify.Watcher
	Templates *template.Template
	Feed      []byte
	Tor       *tor
	Listener  net.Listener
	Router    *mux.Router
}

// NewApp returns a new instance of App from Config.
func NewApp(cfg *Config) (*App, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	a := &App{
		Config: cfg,
	}
	// Setup Library
	a.Library = media.NewLibrary()
	// Setup Watcher
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	a.Watcher = w
	// Setup Listener
	ln, err := newListener(cfg.Server)
	if err != nil {
		return nil, err
	}
	a.Listener = ln
	// Setup Templates
	a.Templates = template.Must(template.ParseGlob("templates/*"))
	// Setup Tor
	if cfg.Tor.Enable {
		t, err := newTor(cfg.Tor)
		if err != nil {
			return nil, err
		}
		a.Tor = t
	}
	// Setup Router
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", a.indexHandler).Methods("GET")
	r.HandleFunc("/v", a.homeHandler).Methods("GET")
	r.HandleFunc("/f/{id}", a.fileHandler).Methods("GET")
	r.HandleFunc("/f/{prefix:.*}/{id}", a.fileHandler).Methods("GET")
	r.HandleFunc("/t/{id}", a.thumbHandler).Methods("GET")
	r.HandleFunc("/t/{prefix:.*}/{id}", a.thumbHandler).Methods("GET")
	r.HandleFunc("/v/{id}", a.pageHandler).Methods("GET")
	r.HandleFunc("/v/{prefix:.*}/{id}", a.pageHandler).Methods("GET")
	r.HandleFunc("/m/{id}", a.musicHandler).Methods("GET")
	r.HandleFunc("/m/{prefix:.*}/{id}", a.musicHandler).Methods("GET")
	r.HandleFunc("/e/{id}", a.embedHandler).Methods("GET")
	r.HandleFunc("/e/{prefix:.*}/{id}", a.embedHandler).Methods("GET")
	r.HandleFunc("/i", a.imageHandler).Methods("GET")
	r.HandleFunc("/i/{prefix:.*}", a.imageHandler).Methods("GET")
	r.HandleFunc("/feed.xml", a.rssHandler).Methods("GET")
	// Static file handler
	fsHandler := http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("./static/")),
	)
	r.PathPrefix("/static/").Handler(fsHandler).Methods("GET")
	a.Router = r
	return a, nil
}

// Run imports the library and starts server.
func (a *App) Run() error {
	if a.Tor != nil {
		var err error
		cs := a.Config.Server
		key := a.Tor.OnionKey
		if key == nil {
			key, err = onionkey.GenerateKey()
			if err != nil {
				return err
			}
			a.Tor.OnionKey = key
		}
		onion, err := key.Onion()
		if err != nil {
			return err
		}
		onion.Ports[80] = fmt.Sprintf("%s:%d", cs.Host, cs.Port)
		err = a.Tor.Controller.AddOnion(onion)
		if err != nil {
			return errors.New("unable to start Tor onion service")
		}
		log.Printf("Onion service: http://%s.onion", onion.ServiceID)
	}
	for _, pc := range a.Config.Library {
		p := &media.Path{
			Path:   pc.Path,
			Prefix: pc.Prefix,
		}
		err := a.Library.AddPath(p)
		if err != nil {
			return err
		}
		err = a.Library.Import(p)
		if err != nil {
			return err
		}
	}
	for path := range a.Library.Paths {
		a.Watcher.Add(path)
	}
	buildFeed(a)
	go startWatcher(a)
	return http.Serve(a.Listener, a.Router)
}

// HTTP handler for /
func (a *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("/")
	http.Redirect(w, r, "/v", 302)
}

// HTTP handler for /v
func (a *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("/v")
	a.Templates.ExecuteTemplate(w, "home.html", &struct {
		Playlist media.Playlist
	}{
		Playlist: a.Library.Playlist(),
	})
}

// HTTP handler for /v/id
func (a *App) pageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	prefix, ok := vars["prefix"]
	if ok {
		id = path.Join(prefix, id)
	}
	log.Printf("/v/%s", id)
	playing, ok := a.Library.Videos[id]
	if !ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		a.Templates.ExecuteTemplate(w, "video.html", &struct {
			Playing  *media.Video
			Playlist media.Playlist
			Captions map[string]*media.Caption
		}{
			Playing:  &media.Video{ID: ""},
			Playlist: a.Library.Playlist(),
			Captions: a.Library.Captions,
		})
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	a.Templates.ExecuteTemplate(w, "video.html", &struct {
		Playing  *media.Video
		Playlist media.Playlist
		Captions map[string]*media.Caption
	}{
		Playing:  playing,
		Playlist: a.Library.Playlist(),
		Captions: a.Library.Captions,
	})
}

// HTTP handler for /m/id
func (a *App) musicHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	prefix, ok := vars["prefix"]
	if ok {
		id = path.Join(prefix, id)
	}
	log.Printf("/m/%s", id)
	playing, ok := a.Library.Videos[id]
	if !ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		a.Templates.ExecuteTemplate(w, "music.html", &struct {
			Playing  *media.Video
			Playlist media.Playlist
			Captions map[string]*media.Caption
		}{
			Playing:  &media.Video{ID: ""},
			Playlist: a.Library.Playlist(),
			Captions: a.Library.Captions,
		})
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	a.Templates.ExecuteTemplate(w, "music.html", &struct {
		Playing  *media.Video
		Playlist media.Playlist
		Captions map[string]*media.Caption
	}{
		Playing:  playing,
		Playlist: a.Library.Playlist(),
		Captions: a.Library.Captions,
	})
}

// HTTP handler for /e/id
func (a *App) embedHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	prefix, ok := vars["prefix"]
	if ok {
		id = path.Join(prefix, id)
	}
	log.Printf("/e/%s", id)
	playing, ok := a.Library.Videos[id]
	if !ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		a.Templates.ExecuteTemplate(w, "embed.html", &struct {
			Playing  *media.Video
			Captions map[string]*media.Caption
		}{
			Playing:  &media.Video{ID: ""},
			Captions: a.Library.Captions,
		})
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	a.Templates.ExecuteTemplate(w, "embed.html", &struct {
		Playing  *media.Video
		Captions map[string]*media.Caption
	}{
		Playing:  playing,
		Captions: a.Library.Captions,
	})
}

// HTTP handler for /f/id
func (a *App) fileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	prefix, ok := vars["prefix"]
	if ok {
		id = path.Join(prefix, id)
	}
	log.Printf("/f/%s", id)
	switch strings.ToUpper(id)[strings.LastIndex(id, ".")+1:] {
	case "VTT":
		m, ok := a.Library.Captions[id]
		if !ok {
			return
		}
		name := m.FileName
		disposition := "attachment; filename=\"" + name + "\""
		w.Header().Set("Content-Disposition", disposition)
		w.Header().Set("Content-Type", "text/vtt")
		http.ServeFile(w, r, m.Path)
	case "PNG", "BMP", "GIF", "WEBP", "JPG", "JPEG":
		m, ok := a.Library.Images[id]
		if !ok {
			return
		}
		name := m.FileName
		disposition := "attachment; filename=\"" + name + "\""
		w.Header().Set("Content-Disposition", disposition)
		w.Header().Set("Content-Type", m.MIMEType)
		http.ServeFile(w, r, m.Path)
	default:
		m, ok := a.Library.Videos[id]
		if !ok {
			return
		}
		name := m.FileName
		disposition := "attachment; filename=\"" + name + "\""
		w.Header().Set("Content-Disposition", disposition)
		w.Header().Set("Content-Type", m.MIMEType)
		http.ServeFile(w, r, m.Path)
	}
}

// HTTP handler for /t/id
func (a *App) thumbHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	prefix, ok := vars["prefix"]
	if ok {
		id = path.Join(prefix, id)
	}
	log.Printf("/t/%s", id)
	m, ok := a.Library.Videos[id]
	if !ok {
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	if m.ThumbType == "" {
		w.Header().Set("Content-Type", "image/png")
		http.ServeFile(w, r, "static/defaulticon.png")
	} else {
		w.Header().Set("Content-Type", m.ThumbType)
		w.Write(m.Thumb)
	}
}

// HTTP handler for /i/id
func (a *App) imageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	prefix, ok := vars["prefix"]
	if !ok {
		prefix = ""
	}
	log.Printf("/i/%s", prefix)
	paths := make(map[string]*media.Path, 0)
	for k, v := range a.Library.Paths {
		if strings.HasPrefix(v.Prefix, prefix) {
			paths[k] = v
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	a.Templates.ExecuteTemplate(w, "image.html", &struct {
		ImageList media.ImageList
		Paths     map[string]*media.Path
		Prefix    string
	}{
		ImageList: a.Library.ImageList(),
		Paths:     paths,
		Prefix:    prefix,
	})
}

// HTTP handler for /feed.xml
func (a *App) rssHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("/feed.xml")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	w.Header().Set("Content-Type", "text/xml")
	w.Write(a.Feed)
}
