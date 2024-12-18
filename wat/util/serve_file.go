package wat

import (
	"bytes"
	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
	"io/fs"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	cache      = make(map[string][]byte)
	cacheMutex = &sync.Mutex{}
)

func getCachedContent(path string, filepath string) []byte {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	if content, found := cache[path]; found && !Config.Debug {
		return content
	}
	var content []byte
	content, _ = os.ReadFile(filepath)
	if !strings.Contains(filepath, "imgs") {
		var compressedContent bytes.Buffer
		compressionLevel := brotli.BestSpeed
		if !Config.Debug {
			compressionLevel = brotli.BestCompression
		}
		writer := brotli.NewWriterLevel(&compressedContent, compressionLevel)
		_, err := writer.Write(content)
		if err != nil {
			return nil
		}
		writer.Close()
		compressedData := compressedContent.Bytes()
		cache[path] = compressedData
	} else {
		cache[path] = content
	}
	return cache[path]
}

func ServeFile(relPath, path string, r *gin.RouterGroup) {
	r.GET(relPath, func(c *gin.Context) {
		content := getCachedContent(relPath, path)
		contentType := mime.TypeByExtension(filepath.Ext(path))

		if !strings.Contains(path, "imgs") {
			c.Header("Content-Encoding", "br")
		}
		c.Header("Content-Type", contentType)
		c.Data(200, contentType, content)
	})
}

func ServeFolder(rootPath, dir string, r *gin.RouterGroup) {
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath, _ := filepath.Rel(dir, path)
			urlPath := rootPath + relativePath

			ServeFile(urlPath, path, r)
		}
		return nil
	})
}
