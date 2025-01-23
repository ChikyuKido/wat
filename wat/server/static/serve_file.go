package wat

import (
	middleware "github.com/ChikyuKido/wat/wat/server/middleware"
	wat "github.com/ChikyuKido/wat/wat/util"
	"github.com/gin-gonic/gin"
	"io/fs"
	"mime"
	"path/filepath"
	"strconv"
	"strings"
)

type DataLoader func(c *gin.Context) any

func ServeFile(diskPath string, dataLoader DataLoader, cacheArena string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data any = nil
		// if the data is already cached don't load the data again because it isn't used
		if dataLoader != nil && GetLoadingState(cacheArena, c.Request.URL.String()) != Finished {
			data = dataLoader(c)
		}
		content := LoadFile(c.Request.URL.String(), diskPath, data, cacheArena)
		contentType := mime.TypeByExtension(filepath.Ext(diskPath))

		if !wat.Config.Debug && (filepath.Ext(diskPath) == ".css" || filepath.Ext(diskPath) == ".js" || filepath.Ext(diskPath) == ".webp") {
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		}

		if !strings.Contains(contentType, "image") {
			c.Header("Content-Encoding", "br")
		}
		c.Header("Content-Type", contentType)
		c.Header("X-Transfer-Size", strconv.Itoa(len(content)))
		c.Data(200, contentType, content)
	}
}

func ServeFolder(rootPath, dir string, dataLoader DataLoader, cacheArena string, r *gin.RouterGroup, permission string) {
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath, _ := filepath.Rel(dir, path)
			urlPath := rootPath + relativePath
			if permission != "" {
				r.GET(urlPath, middleware.RequiredPermission(permission, true), ServeFile(path, dataLoader, cacheArena))
			} else {
				r.GET(urlPath, ServeFile(path, dataLoader, cacheArena))
			}
		}
		return nil
	})
}
