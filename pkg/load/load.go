package load

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
)

var (
	baseOSFS      = afero.NewOsFs()
	layerMemFS    = afero.NewMemMapFs()
	roBaseFS      = afero.NewReadOnlyFs(baseOSFS)
	cacheOnReadFS = afero.NewCacheOnReadFs(roBaseFS, layerMemFS, 0) // 0 means unlimited layerMemFS time
	afs           = &afero.Afero{Fs: cacheOnReadFS}
)

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type InputFile struct {
	path string
	info os.FileInfo
}

func (f *InputFile) Path() string {
	return f.path
}

func (f *InputFile) Name() string {
	return f.info.Name()
}

func (f *InputFile) IsDir() bool {
	return f.info.IsDir()
}

func (f *InputFile) IsFile() bool {
	return !f.IsDir()
}

func (f *InputFile) Ext() string {
	return filepath.Ext(f.Name())
}

func (f *InputFile) Contents() ([]byte, error) {
	if f.IsDir() {
		return nil, fmt.Errorf("It is not a file: %s", f.path)
	}
	return afs.ReadFile(f.path)
}

func (f *InputFile) Open() (afero.File, error) {
	return afs.Open(f.path)
}

func ListFiles(paths []string) ([]InputFile, error) {
	result := []InputFile{}
	for _, path := range paths {
		err := afs.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			absPath, err := filepath.Abs(walkPath)
			if err != nil {
				return err
			}

			result = append(result, InputFile{path: absPath, info: info})
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	uniq := lo.UniqBy[InputFile, string](result, func(f InputFile) string { return f.path })
	return uniq, nil
}
