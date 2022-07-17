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

func PathExists(path string) (bool, error) {
	return afs.Exists(path)
}

type Input struct {
	path           string
	info           os.FileInfo
	detected_types []DetectedType
}

func (i *Input) Path() string {
	return i.path
}

func (i *Input) Name() string {
	return i.info.Name()
}

func (i *Input) IsDir() bool {
	return i.info.IsDir()
}

func (i *Input) IsFile() bool {
	return !i.IsDir()
}

func (i *Input) Ext() string {
	return filepath.Ext(i.Name())
}

func (i *Input) Contents() ([]byte, error) {
	if i.IsDir() {
		return nil, fmt.Errorf("It is not a file: %s", i.Path())
	}
	return afs.ReadFile(i.Path())
}

func (i *Input) Open() (afero.File, error) {
	if i.IsDir() {
		return nil, fmt.Errorf("It is not a file: %s", i.Path())
	}
	return afs.Open(i.Path())
}

func (i *Input) DetectedTypes() []DetectedType {
	return i.detected_types
}

func FlatPathsToInputs(paths []string) ([]Input, error) {
	var inputs []Input
	for _, path := range paths {
		err := afs.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// TODO move the following logic to parameters, use runner options?
			if info.Name() == ".git" {
				return filepath.SkipDir
			}

			absPath, err := filepath.Abs(walkPath)
			if err != nil {
				return err
			}

			inputs = append(inputs, Input{path: absPath, info: info})
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	uniqInputs := lo.UniqBy[Input, string](inputs, func(f Input) string { return f.path })

	for idx, input := range uniqInputs {
		uniqInputs[idx].detected_types = detectTypesOfInput(input)
	}

	return uniqInputs, nil
}
