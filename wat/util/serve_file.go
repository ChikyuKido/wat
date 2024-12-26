package wat

import (
	"bytes"
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"html/template"
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
	templates  []string
)

func unescapeJavaScript(content []byte) []byte {
	strContent := string(content)
	strContent = strings.ReplaceAll(strContent, "&lt;", "<")
	strContent = strings.ReplaceAll(strContent, "&gt;", ">")
	strContent = strings.ReplaceAll(strContent, "&amp;", "&")
	return []byte(strContent)
}

func getCachedContent(path string, filepath string) []byte {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	if content, found := cache[path]; found && !Config.Debug {
		return content
	}
	var content []byte
	if strings.HasSuffix(filepath, ".html") || strings.HasSuffix(filepath, ".js") {
		byteBuffer := bytes.NewBuffer(make([]byte, 0))
		t, err := template.ParseFiles(append([]string{filepath}, templates...)...)
		if err != nil {
			fmt.Printf("Failed to parse template %v", err)
			return nil
		}
		err = t.Execute(byteBuffer, nil)
		if err != nil {
			fmt.Printf("Failed to execute template %v", err)
			return nil
		}
		content = byteBuffer.Bytes()
		if strings.HasSuffix(filepath, ".js") {
			content = unescapeJavaScript(content)
		}
	} else {
		content, _ = os.ReadFile(filepath)
	}
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
func LoadTemplates(directory string) {
	var files []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".html" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		logrus.Fatalf("failed to load templates")
	}
	templates = files
}

func ServeFile(relPath, path string, r *gin.RouterGroup) {
	r.GET(relPath, func(c *gin.Context) {
		content := getCachedContent(relPath, path)
		contentType := mime.TypeByExtension(filepath.Ext(path))

		if !strings.Contains(contentType, "image") {
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
