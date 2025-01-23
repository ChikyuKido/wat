package wat

import (
	"bytes"
	"fmt"
	util "github.com/ChikyuKido/wat/wat/util"
	"github.com/andybalholm/brotli"
	"github.com/sirupsen/logrus"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

type LoadingState int

const (
	NotLoaded LoadingState = iota
	Loading
	Finished
	Invalidated
)

type Cache struct {
	content      []byte
	loadingState LoadingState
}

var (
	caches    = make(map[string]map[string]*Cache)
	templates []string
)

func LoadFile(path string, diskPath string, data any, cacheArena string) []byte {
	if _, exists := caches[cacheArena]; !exists {
		caches[cacheArena] = make(map[string]*Cache)
	}
	cache, exists := caches[cacheArena][path]
	if !exists {
		cache = &Cache{
			content:      make([]byte, 0),
			loadingState: NotLoaded,
		}
		caches[cacheArena][path] = cache
	}
	if cache.loadingState == Loading {
		logrus.Debug("cache is still loading. Returning best speed compression")
		return loadFile(path, data, brotli.BestSpeed)
	}
	if cache.loadingState == Finished {
		logrus.Debug("cache is finished. Returning it")
		if content := cache.content; len(content) > 0 {
			return content
		}
	}
	go func() {
		cache.loadingState = Loading
		var content = loadFile(diskPath, data, brotli.BestCompression)
		if cache.loadingState != Invalidated { // if the cache was invalided during an update
			cache.content = content
			cache.loadingState = Finished
			logrus.Debug("Best compression run finished.")
		} else {
			logrus.Debug("Cache was invalidated during a best compression run. Throw away the result")
		}
	}()
	logrus.Debug("Cache is not loaded yet. Load it and return best speed cache")
	return loadFile(diskPath, data, brotli.BestSpeed)
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

func InvalidateArena(arena string) {
	logrus.Debug("Invalidate arena: " + arena)
	cache, exists := caches[arena]
	if exists {
		for _, c := range cache {
			c.content = make([]byte, 0)
			c.loadingState = Invalidated
		}
	}
}
func GetLoadingState(cacheArena, path string) LoadingState {
	if _, exists := caches[cacheArena]; !exists {
		caches[cacheArena] = make(map[string]*Cache)
	}
	cache, exists := caches[cacheArena][path]
	if !exists {
		cache = &Cache{
			content:      make([]byte, 0),
			loadingState: NotLoaded,
		}
		caches[cacheArena][path] = cache
	}
	return cache.loadingState
}
func loadFile(path string, data any, compression int) []byte {
	var content []byte
	if strings.HasSuffix(path, ".html") || strings.HasSuffix(path, ".js") {
		byteBuffer := bytes.NewBuffer(make([]byte, 0))
		t, err := template.ParseFiles(append([]string{path}, templates...)...)
		if err != nil {
			fmt.Printf("Failed to parse template %v", err)
			return nil
		}
		err = t.Execute(byteBuffer, data)
		if err != nil {
			fmt.Printf("Failed to execute template %v", err)
			return nil
		}
		content = byteBuffer.Bytes()
		if strings.HasSuffix(path, ".js") {
			content = unescapeJavaScript(content)
		}
	} else {
		content, _ = os.ReadFile(path)
	}
	if strings.Contains(path, "html") && !util.Config.Debug {
		bodyStr := string(content)
		fmt.Println(bodyStr)
		bodyStr = strings.ReplaceAll(bodyStr, "{rep}", util.Config.ResourceVersion)
		content = []byte(bodyStr)
	}
	if !strings.Contains(path, "imgs") {
		var compressedContent bytes.Buffer
		writer := brotli.NewWriterLevel(&compressedContent, compression)
		_, err := writer.Write(content)
		if err != nil {
			return nil
		}
		writer.Close()
		compressedData := compressedContent.Bytes()
		return compressedData
	} else {
		return content
	}
}

func unescapeJavaScript(content []byte) []byte {
	strContent := string(content)
	strContent = strings.ReplaceAll(strContent, "&lt;", "<")
	strContent = strings.ReplaceAll(strContent, "&gt;", ">")
	strContent = strings.ReplaceAll(strContent, "&amp;", "&")
	return []byte(strContent)
}
