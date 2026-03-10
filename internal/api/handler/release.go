package handler

import (
	"nebula/internal/api/response"
	"nebula/internal/release"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReleaseHandler struct {
	service *release.ReleaseService
}

func NewReleaseHandler(service *release.ReleaseService) *ReleaseHandler {
	return &ReleaseHandler{service: service}
}

// List 获取所有版本列表
// GET /api/releases
func (h *ReleaseHandler) List(c *gin.Context) {
	releases, err := h.service.List()
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.Ok(c, releases)
}

// ListByApp 获取指定应用的所有版本
// GET /api/apps/:id/releases
func (h *ReleaseHandler) ListByApp(c *gin.Context) {
	appID := c.Param("id")
	channel := c.Query("channel") // 可选的渠道过滤

	var releases []release.Release
	var err error

	if channel != "" {
		releases, err = h.service.ListByAppAndChannel(appID, channel)
	} else {
		releases, err = h.service.ListByApp(appID)
	}

	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.Ok(c, releases)
}

// Get 获取单个版本详情
// GET /api/releases/:id
func (h *ReleaseHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailBadRequest(c, "invalid release id")
		return
	}

	rel, err := h.service.Get(uint(id))
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.Ok(c, rel)
}

// GetLatest 获取应用的最新版本
// GET /api/apps/:id/releases/latest
func (h *ReleaseHandler) GetLatest(c *gin.Context) {
	appID := c.Param("id")
	channel := c.Query("channel")

	rel, err := h.service.GetLatest(appID, channel)
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.Ok(c, rel)
}

// Create 创建新版本
// POST /api/releases
func (h *ReleaseHandler) Create(c *gin.Context) {
	var rel release.Release
	if err := c.ShouldBindJSON(&rel); err != nil {
		response.FailBadRequest(c, err.Error())
		return
	}

	// 基本验证
	if rel.AppID == "" {
		response.FailBadRequest(c, "app_id is required")
		return
	}
	if rel.Version == "" {
		response.FailBadRequest(c, "version is required")
		return
	}

	err := h.service.Create(rel)
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.OkMsg(c, "created")
}

// Update 更新版本信息
// PUT /api/releases/:id
func (h *ReleaseHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailBadRequest(c, "invalid release id")
		return
	}

	var rel release.Release
	if err := c.ShouldBindJSON(&rel); err != nil {
		response.FailBadRequest(c, err.Error())
		return
	}

	data := map[string]any{
		"version":  rel.Version,
		"notes":    rel.Notes,
		"channel":  rel.Channel,
		"pub_date": rel.PubDate,
	}

	err = h.service.Update(uint(id), data)
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.OkMsg(c, "updated")
}

// Delete 删除版本
// DELETE /api/releases/:id
func (h *ReleaseHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailBadRequest(c, "invalid release id")
		return
	}

	err = h.service.Delete(uint(id))
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.OkMsg(c, "deleted")
}
