package app

import (
	"net/http"
	"log"
	"path"
	"strings"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/178619/stube/pkg/media"
)

type InfoResponse struct {
    Status      string           `json:"status"`
    ID          string           `json:"id,omitempty"`
    Ref         string           `json:"ref,omitempty"`
    Title       string           `json:"title,omitempty"`
    Album       string           `json:"album"`
    Disc        int              `json:"disc"`
    Track       int              `json:"track"`
    Description string           `json:"description"`
    Modified    string           `json:"modified"`
    Size        int64            `json:"size"`
    MIMEType    string           `json:"mimeType"`
    Captions    []CaptionSummary `json:"captions"`
}

type CaptionSummary struct {
    Ref     string `json:"ref"`
    SrcLang string `json:"srcLang"`
}

type ErrorResponse struct {
    Status string `json:"status"`
    Error  struct {
        Code    int    `json:"code"`
        Message string `json:"message"`
    } `json:"error"`
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
		w.WriteHeader(404)
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
		http.ServeFile(w, r, "static/default.png")
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

// HTTP handler for /n/id
func (a *App) infoHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    if prefix, ok := vars["prefix"]; ok {
        id = path.Join(prefix, id)
    }
    log.Printf("/n/%s", id)
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    playing, ok := a.Library.Videos[id]
    if !ok {
        w.WriteHeader(404)
        json.NewEncoder(w).Encode(ErrorResponse{
            Status: "error",
            Error: struct {
                Code    int    `json:"code"`
                Message string `json:"message"`
            }{Code: 404,Message: "File Not Found"},
        })
        return
    }

    matchedCaptions := []CaptionSummary{}
    for _, m := range a.Library.Captions {
        if m.Origin == playing.FileName || m.Origin == playing.BaseName {
            matchedCaptions = append(matchedCaptions, CaptionSummary{
                Ref:     m.Ref,
                SrcLang: m.SrcLang,
            })
        }
    }

    json.NewEncoder(w).Encode(InfoResponse{
        Status:      "success",
        ID:          playing.ID,
        Ref:         playing.Ref,
        Title:       playing.Title,
        Album:       playing.Album,
        Disc:        playing.Disc,
        Track:       playing.Track,
        Description: playing.Description,
        Modified:    playing.Modified,
        Size:        playing.Size,
        MIMEType:    playing.MIMEType,
        Captions:    matchedCaptions,
    })
}

// HTTP handler for /feed.xml
func (a *App) rssHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("/feed.xml")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	w.Header().Set("Content-Type", "text/xml")
	w.Write(a.Feed)
}