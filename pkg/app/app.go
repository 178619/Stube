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

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	"github.com/178619/tube/pkg/media"
	"github.com/178619/tube/pkg/onionkey"
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
		a.Watcher.Add(p.Path)
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
		a.Templates.ExecuteTemplate(w, "player.html", &struct {
			Playing  *media.Video
			Playlist media.Playlist
			Captions media.CaptionList
		}{
			Playing:  &media.Video{ID: ""},
			Playlist: a.Library.Playlist(),
			Captions: a.Library.CaptionList(),
		})
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	a.Templates.ExecuteTemplate(w, "player.html", &struct {
		Playing  *media.Video
		Playlist media.Playlist
		Captions media.CaptionList
	}{
		Playing:  playing,
		Playlist: a.Library.Playlist(),
		Captions: a.Library.CaptionList(),
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
		a.Templates.ExecuteTemplate(w, "player_music.html", &struct {
			Playing  *media.Video
			Playlist media.Playlist
			Captions media.CaptionList
		}{
			Playing:  &media.Video{ID: ""},
			Playlist: a.Library.Playlist(),
			Captions: a.Library.CaptionList(),
		})
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	a.Templates.ExecuteTemplate(w, "player_music.html", &struct {
		Playing  *media.Video
		Playlist media.Playlist
		Captions media.CaptionList
	}{
		Playing:  playing,
		Playlist: a.Library.Playlist(),
		Captions: a.Library.CaptionList(),
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
			Playlist media.Playlist
			Captions media.CaptionList
		}{
			Playing:  &media.Video{ID: ""},
			Playlist: a.Library.Playlist(),
			Captions: a.Library.CaptionList(),
		})
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	a.Templates.ExecuteTemplate(w, "embed.html", &struct {
		Playing  *media.Video
		Playlist media.Playlist
		Captions media.CaptionList
	}{
		Playing:  playing,
		Playlist: a.Library.Playlist(),
		Captions: a.Library.CaptionList(),
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
	m, ok := a.Library.Videos[id]
	if !ok {
		m, ok := a.Library.Captions[id]
		if !ok {
			return
		}
		name := m.FileName
		disposition := "attachment; filename=\"" + name + "\""
		w.Header().Set("Content-Disposition", disposition)
		w.Header().Set("Content-Type", "text/vtt")
		http.ServeFile(w, r, m.Path)
		return
	}
	name := m.FileName
	disposition := "attachment; filename=\"" + name + "\""
	mimeType := m.MIMEType
	w.Header().Set("Content-Disposition", disposition)
	w.Header().Set("Content-Type", mimeType)
	http.ServeFile(w, r, m.Path)
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

// HTTP handler for /feed.xml
func (a *App) rssHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	w.Header().Set("Content-Type", "text/xml")
	w.Write(a.Feed)
}
