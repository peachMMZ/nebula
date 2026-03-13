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

// List 获取指定应用的所有版本 (通过 app_name)
// GET /api/releases?app_name=xxx&channel=yyy
func (h *ReleaseHandler) List(c *gin.Context) {
	appName := c.Query("app_name")
	if appName == "" {
		response.FailBadRequest(c, "app_name is required")
		return
	}

	channel := c.Query("channel") // 可选的渠道过滤

	releases, err := h.service.List(appName, channel)
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.Ok(c, releases)
}

// ListByApp 获取指定应用的所有版本 (通过 app_name)
// GET /api/apps/:name/releases
func (h *ReleaseHandler) ListByApp(c *gin.Context) {
	appName := c.Param("name")
	channel := c.Query("channel") // 可选的渠道过滤

	var releases []release.Release
	var err error

	if channel != "" {
		releases, err = h.service.ListByAppAndChannel(appName, channel)
	} else {
		releases, err = h.service.ListByApp(appName)
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

// GetLatest 获取应用的最新版本 (通过 app_name)
// GET /api/apps/:name/releases/latest
func (h *ReleaseHandler) GetLatest(c *gin.Context) {
	appName := c.Param("name")
	channel := c.Query("channel")

	rel, err := h.service.GetLatest(appName, channel)
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.Ok(c, rel)
}

// Create 创建新版本
// POST /api/:name/releases
func (h *ReleaseHandler) Create(c *gin.Context) {
	appName := c.Param("name")

	var rel release.Release
	if err := c.ShouldBindJSON(&rel); err != nil {
		response.FailBadRequest(c, err.Error())
		return
	}

	// 基本验证
	if rel.Tag == "" {
		response.FailBadRequest(c, "tag is required")
		return
	}
	if rel.Version == "" {
		response.FailBadRequest(c, "version is required")
		return
	}

	// 通过 app_name 查找 app_id
	err := h.service.CreateByAppName(appName, rel)
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.OkMsg(c, "created")
}

// Update 更新版本信息
// PUT /api/:name/releases/:tag
func (h *ReleaseHandler) Update(c *gin.Context) {
	appName := c.Param("name")
	tag := c.Param("tag")

	var rel release.Release
	if err := c.ShouldBindJSON(&rel); err != nil {
		response.FailBadRequest(c, err.Error())
		return
	}

	err := h.service.UpdateByTag(appName, tag, rel)
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.OkMsg(c, "updated")
}

// Delete 删除版本
// DELETE /api/:name/releases/:tag
func (h *ReleaseHandler) Delete(c *gin.Context) {
	appName := c.Param("name")
	tag := c.Param("tag")

	err := h.service.DeleteByTag(appName, tag)
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.OkMsg(c, "deleted")
}
