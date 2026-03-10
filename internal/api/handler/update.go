package handler

import (
	"nebula/internal/api/response"
	"nebula/internal/updater"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CheckUpdate(c *gin.Context) {
	req := updater.CheckRequest{
		App:      c.Query("app"),
		Version:  c.Query("version"),
		Platform: c.Query("platform"),
		Arch:     c.Query("arch"),
	}

	resp, err := updater.CheckUpdate(h.db, req)

	if err != nil {
		response.FailServer(c, err.Error())
		return
	}

	response.Ok(c, resp)
}
