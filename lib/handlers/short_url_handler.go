package handlers

import (
	"github.com/azach/shorturl/lib/pool"
	"github.com/azach/shorturl/lib/storage"
)

type ShortURLHandler struct {
	urlStorage storage.Storage
	urlPool    *pool.Pool
}

func NewShortUrlHandler(urlStorage storage.Storage, urlPool *pool.Pool) *ShortURLHandler {
	return &ShortURLHandler{
		urlStorage: urlStorage,
		urlPool:    urlPool,
	}
}

func (h *ShortURLHandler) GeneratePool() {
	for {
		h.urlPool.Generate()
	}
}
