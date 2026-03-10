package handler

import (
	"nebula/internal/api/response"
	"nebula/internal/app"

	"github.com/gin-gonic/gin"
)

type AppHandler struct {
	service *app.AppService
}

func NewAppHandler(service *app.AppService) *AppHandler {
	return &AppHandler{service: service}
}

func (h *AppHandler) List(c *gin.Context) {
	params := make(map[string]any)

	// 从查询参数中获取过滤条件
	if name := c.Query("name"); name != "" {
		params["name"] = name
	}
	if description := c.Query("description"); description != "" {
		params["description"] = description
	}

	apps, err := h.service.List(params)
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.Ok(c, apps)
}

func (h *AppHandler) Get(c *gin.Context) {
	id := c.Param("id")
	app, err := h.service.Get(id)
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.Ok(c, app)
}

func (h *AppHandler) Create(c *gin.Context) {
	var app app.App
	if err := c.ShouldBindJSON(&app); err != nil {
		response.FailBadRequest(c, err.Error())
		return
	}
	err := h.service.Create(app)
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.OkMsg(c, "created")
}

func (h *AppHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var app app.App
	if err := c.ShouldBindJSON(&app); err != nil {
		response.FailBadRequest(c, err.Error())
		return
	}
	err := h.service.Update(id, map[string]any{
		"name":        app.Name,
		"description": app.Description,
	})
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.OkMsg(c, "updated")
}

func (h *AppHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.service.Delete(id)
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.OkMsg(c, "deleted")
}
