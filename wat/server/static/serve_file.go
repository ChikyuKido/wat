package static

import (
	"github.com/gin-gonic/gin"
	"io/fs"
	"mime"
	"path/filepath"
	"strconv"
	"strings"
)

func ServeFile(diskPath string, data any, cacheArena string) gin.HandlerFunc {
	return func(c *gin.Context) {
		content := LoadFile(diskPath, data, cacheArena)
		contentType := mime.TypeByExtension(filepath.Ext(diskPath))

		if !strings.Contains(contentType, "image") {
			c.Header("Content-Encoding", "br")
		}
		c.Header("Content-Type", contentType)
		c.Header("X-Transfer-Size", strconv.Itoa(len(content)))
		c.Data(200, contentType, content)
	}
}

func ServeFolder(rootPath, dir string, data any, cacheArena string, r *gin.RouterGroup) {
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath, _ := filepath.Rel(dir, path)
			urlPath := rootPath + relativePath

			r.GET(urlPath, ServeFile(path, data, cacheArena))
		}
		return nil
	})
}
