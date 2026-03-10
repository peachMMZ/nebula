package handler

import (
	"nebula/internal/api/response"
	"nebula/internal/asset"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AssetHandler struct {
	service *asset.AssetService
}

func NewAssetHandler(service *asset.AssetService) *AssetHandler {
	return &AssetHandler{service: service}
}

// List 获取所有资源列表
// GET /api/assets
func (h *AssetHandler) List(c *gin.Context) {
	assets, err := h.service.List()
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.Ok(c, assets)
}

// ListByRelease 获取指定发布版本的所有资源
// GET /api/releases/:id/assets
func (h *AssetHandler) ListByRelease(c *gin.Context) {
	releaseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailBadRequest(c, "invalid release id")
		return
	}

	assets, err := h.service.ListByRelease(uint(releaseID))
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.Ok(c, assets)
}

// Get 获取单个资源详情
// GET /api/assets/:id
func (h *AssetHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailBadRequest(c, "invalid asset id")
		return
	}

	ast, err := h.service.Get(uint(id))
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.Ok(c, ast)
}

// Upload 上传文件并创建资源
// POST /api/releases/:id/assets/upload
func (h *AssetHandler) Upload(c *gin.Context) {
	releaseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailBadRequest(c, "invalid release id")
		return
	}

	// 获取表单参数
	platform := c.PostForm("platform")
	arch := c.PostForm("arch")

	if platform == "" || arch == "" {
		response.FailBadRequest(c, "platform and arch are required")
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.FailBadRequest(c, "file is required")
		return
	}

	// 验证文件扩展名（可选）
	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{
		".exe": true, ".dmg": true, ".pkg": true, ".deb": true,
		".rpm": true, ".appimage": true, ".zip": true, ".tar.gz": true,
		".msi": true, ".app": true,
	}
	if !allowedExts[ext] {
		response.FailBadRequest(c, "unsupported file type: "+ext)
		return
	}

	// 上传文件
	ast, err := h.service.Upload(uint(releaseID), platform, arch, file)
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}

	response.Ok(c, ast)
}

// Create 创建资源记录（用于外部 URL）
// POST /api/assets
func (h *AssetHandler) Create(c *gin.Context) {
	var ast asset.Asset
	if err := c.ShouldBindJSON(&ast); err != nil {
		response.FailBadRequest(c, err.Error())
		return
	}

	// 基本验证
	if ast.ReleaseID == 0 {
		response.FailBadRequest(c, "releaseId is required")
		return
	}
	if ast.Platform == "" {
		response.FailBadRequest(c, "platform is required")
		return
	}
	if ast.Arch == "" {
		response.FailBadRequest(c, "arch is required")
		return
	}
	if ast.URL == "" {
		response.FailBadRequest(c, "url is required")
		return
	}

	err := h.service.Create(ast)
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.OkMsg(c, "created")
}

// Update 更新资源信息
// PUT /api/assets/:id
func (h *AssetHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailBadRequest(c, "invalid asset id")
		return
	}

	var ast asset.Asset
	if err := c.ShouldBindJSON(&ast); err != nil {
		response.FailBadRequest(c, err.Error())
		return
	}

	data := map[string]any{
		"platform":  ast.Platform,
		"arch":      ast.Arch,
		"signature": ast.Signature,
		"checksum":  ast.Checksum,
	}

	err = h.service.Update(uint(id), data)
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.OkMsg(c, "updated")
}

// Delete 删除资源
// DELETE /api/assets/:id
func (h *AssetHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailBadRequest(c, "invalid asset id")
		return
	}

	err = h.service.Delete(uint(id))
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.OkMsg(c, "deleted")
}

// Download 下载资源文件
// GET /api/assets/:id/download
func (h *AssetHandler) Download(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailBadRequest(c, "invalid asset id")
		return
	}

	// 获取存储路径
	storagePath, err := h.service.GetStoragePath(uint(id))
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}

	// 返回文件（Gin 会自动处理文件下载）
	c.File(storagePath)
}
