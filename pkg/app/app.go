// Package app manages main application server.
package app

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"

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
	r.HandleFunc("/n/{id}", a.infoHandler).Methods("GET")
	r.HandleFunc("/n/{prefix:.*}/{id}", a.infoHandler).Methods("GET")
	r.HandleFunc("/i", a.imageHandler).Methods("GET")
	r.HandleFunc("/i/{prefix:.*}", a.imageHandler).Methods("GET")
	r.HandleFunc("/feed.xml", a.rssHandler).Methods("GET")
	// Static file handler
	r.PathPrefix("/static/").Handler(http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("./static/")),
	)).Methods("GET")
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
