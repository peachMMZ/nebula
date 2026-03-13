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

// ListByAppAndVersion 获取指定应用和版本的所有资源 (GitHub API 风格)
// GET /api/:name/releases/:version/assets
func (h *AssetHandler) ListByAppAndVersion(c *gin.Context) {
	appName := c.Param("name")
	version := c.Param("version")

	assets, err := h.service.ListByAppAndVersion(appName, version)
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}
	response.Ok(c, assets)
}

// DownloadByAppAndVersion 下载指定应用和版本的资源文件 (GitHub API 风格)
// GET /api/:name/releases/:version/assets/:assetId
func (h *AssetHandler) DownloadByAppAndVersion(c *gin.Context) {
	assetID, err := strconv.ParseUint(c.Param("assetId"), 10, 32)
	if err != nil {
		response.FailBadRequest(c, "invalid asset id")
		return
	}

	// 获取存储路径
	storagePath, err := h.service.GetStoragePath(uint(assetID))
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}

	// 返回文件（Gin 会自动处理文件下载）
	c.File(storagePath)
}

// DownloadByTag 通过 app_name、tag、platform-arch 和 filename 下载资源 (完全符合 GitHub API)
// GET /api/:name/releases/download/:tag/:platformArch/:filename
func (h *AssetHandler) DownloadByTag(c *gin.Context) {
	appName := c.Param("name")
	tag := c.Param("tag")
	platformArch := c.Param("platformArch")
	filename := c.Param("filename")

	// 通过 app_name、tag、platformArch 和 filename 查找资源
	storagePath, err := h.service.GetStoragePathByTag(appName, tag, platformArch, filename)
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}

	// 返回文件（Gin 会自动处理文件下载）
	c.File(storagePath)
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

// Create 创建资源（上传文件并创建记录，原子操作）
// POST /api/:name/releases/:tag/assets
func (h *AssetHandler) Create(c *gin.Context) {
	appName := c.Param("name")
	tag := c.Param("tag")

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

	// 上传文件并创建记录（通过 app_name 和 tag）
	ast, err := h.service.CreateByTag(appName, tag, platform, arch, file)
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}

	response.Ok(c, ast)
}
