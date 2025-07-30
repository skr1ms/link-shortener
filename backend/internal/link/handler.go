package link

import (
	"linkshortener/config"
	"linkshortener/pkg/di"
	"linkshortener/pkg/event"
	"linkshortener/pkg/middleware"
	"linkshortener/pkg/req"
	"linkshortener/pkg/res"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type LinkHandlerDeps struct {
	Config         *config.Config
	LinkRepository *LinkRepository
	EventBus       di.IEventBus
}

type LinkHandler struct {
	deps *LinkHandlerDeps
}

func NewLinkHandler(router *http.ServeMux, deps *LinkHandlerDeps) {
	linkHandler := &LinkHandler{
		deps: deps,
	}
	router.HandleFunc("GET /link/{hash}", linkHandler.GoTo())
	router.Handle("POST /link", middleware.IsAuthenticated(linkHandler.Create(), deps.Config))
	router.Handle("PATCH /link/{id}", middleware.IsAuthenticated(linkHandler.Update(), deps.Config))
	router.Handle("DELETE /link/{id}", middleware.IsAuthenticated(linkHandler.Delete(), deps.Config))
	router.Handle("GET /link", middleware.IsAuthenticated(linkHandler.GetLinks(), deps.Config))
}

func (handler *LinkHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		link, err := handler.deps.LinkRepository.GetByHash(hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		go handler.deps.EventBus.Publish(event.Event{
			Type: event.LinkClicked,
			Data: link.ID,
		})
		http.Redirect(w, r, link.OriginalURL, http.StatusTemporaryRedirect)
	}
}

func (handler *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[CreateLinkRequest](&w, r)
		if err != nil {
			return
		}
		link := NewLink(body.URL)
		createdLink, err := handler.deps.LinkRepository.Create(link)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Response(w, 201, createdLink)
	}
}

func (handler *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := req.HandleBody[UpdateLinkRequest](&w, r)
		if err != nil {
			return
		}

		id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		link, err := handler.deps.LinkRepository.Update(&Link{
			Model: gorm.Model{
				ID: uint(id),
			},
			OriginalURL: body.URL,
			Hash:        body.Hash,
		})

		res.Response(w, 200, link)
	}
}

func (handler *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = handler.deps.LinkRepository.Delete(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Response(w, 200, nil)
	}
}

func (handler *LinkHandler) GetLinks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		links, err := handler.deps.LinkRepository.GetLinks(uint(limit), uint(offset))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		count, err := handler.deps.LinkRepository.GetLinksCount()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Response(w, 200, GetLinksResponse{
			Links: links,
			Count: count,
		})
	}
}
