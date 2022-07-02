package load

import (
	"github.com/awslabs/goformation/v6"
	"github.com/awslabs/goformation/v6/cloudformation"
	"github.com/awslabs/goformation/v6/intrinsics"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func DetectCloudformation(f InputFile) (bool, error) {
	log.Debug().Str("path", f.path).Msg("Detecting Cloudformation")

	if f.IsDir() {
		return false, nil
	}

	if f.Ext() == ".yml" || f.Ext() == ".yaml" || f.Ext() == ".json" {
		data, err := f.Contents()
		if err != nil {
			return false, err
		}

		var template *cloudformation.Template
		options := intrinsics.ProcessorOptions{NoProcess: false, EvaluateConditions: true}

		if f.Ext() == ".json" {
			template, err = goformation.ParseJSONWithOptions(data, &options)
		} else {
			template, err = goformation.ParseYAMLWithOptions(data, &options)
		}
		if err != nil {
			return false, err
		}
		if template.Resources == nil {
			return false, errors.New("Does not contain resources")
		}

		return true, nil
	}

	return false, nil
}
