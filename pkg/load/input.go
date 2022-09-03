package load

import (
	"bufio"
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
	AFS           = &afero.Afero{Fs: cacheOnReadFS}
)

func PathExists(path string) (bool, error) {
	return AFS.Exists(path)
}

type Input struct {
	path          string
	info          os.FileInfo
	detectedTypes []DetectedType
}

func NewInput(path string, info os.FileInfo) *Input {
	p := new(Input)

	p.path = path

	if info == nil {
		info, _ = AFS.Stat(path)
	}
	p.info = info

	return p
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
	return AFS.ReadFile(i.Path())
}

func (i *Input) Open() (afero.File, error) {
	if i.IsDir() {
		return nil, fmt.Errorf("It is not a file: %s", i.Path())
	}
	return AFS.Open(i.Path())
}

func (i *Input) ReadLines(from, to int) ([]string, error) {
	reader, err := i.Open()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var result []string
	scanner := bufio.NewScanner(reader)
	n := 0
	for scanner.Scan() {
		n++
		if n < from {
			continue
		}
		if n > to {
			break
		}

		result = append(result, scanner.Text())
	}

	return result, nil
}

func (i *Input) DetectedTypes() []DetectedType {
	return i.detectedTypes
}

func FlatPathsToInputs(paths []string) ([]Input, error) {
	var inputs []Input
	for _, path := range paths {
		err := AFS.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
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

			inputs = append(inputs, *NewInput(absPath, info))
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	uniqInputs := lo.UniqBy[Input, string](inputs, func(f Input) string { return f.path })

	for idx, input := range uniqInputs {
		uniqInputs[idx].detectedTypes = detectTypesOfInput(input)
	}

	return uniqInputs, nil
}
