package plugins

import (
	"embed"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

//go:embed node/runtime.js python/runtime.py
var runtimeFiles embed.FS

var (
	extractOnce sync.Once
	extractedDir string
)

// RuntimeDir returns the directory containing the extracted plugin runtime
// files (runtime.js, runtime.py). On first call, it extracts the embedded
// files to a temporary directory. Subsequent calls return the cached path.
//
// This solves the "Module not found plugins/node/runtime.js" error that
// occurs when serving from a project directory — the binary no longer
// depends on the working directory or the engine source tree.
func RuntimeDir() string {
	extractOnce.Do(func() {
		dir, err := os.MkdirTemp("", "bitcode-plugins-*")
		if err != nil {
			return
		}
		for _, f := range []string{"node/runtime.js", "python/runtime.py"} {
			data, err := runtimeFiles.ReadFile(f)
			if err != nil {
				continue
			}
			full := filepath.Join(dir, filepath.FromSlash(f))
			os.MkdirAll(filepath.Dir(full), 0755)
			os.WriteFile(full, data, 0644)
		}
		extractedDir = dir
	})
	if extractedDir == "" {
		// Fallback: use source-relative path (dev mode)
		_, file, _, _ := runtime.Caller(0)
		return filepath.Join(filepath.Dir(file))
	}
	return extractedDir
}
