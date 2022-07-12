package load

import (
	"github.com/awslabs/goformation/v6"
	"github.com/awslabs/goformation/v6/cloudformation"
	"github.com/awslabs/goformation/v6/intrinsics"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func init() {
	allDetectors = append(
		allDetectors,
		new(DetectorAwsCloudformation),
	)
}

type DetectorAwsCloudformation struct{}

func (DetectorAwsCloudformation) Detect(input Input) (DetectedType, error) {
	log.Debug().Str("path", input.Path()).Msg("Detecting AWS Cloudformation")

	if input.IsDir() {
		return DetectedUnknown, nil
	}

	if input.Ext() == ".yml" || input.Ext() == ".yaml" || input.Ext() == ".json" {
		data, err := input.Contents()
		if err != nil {
			return DetectedUnknown, err
		}

		var template *cloudformation.Template
		options := intrinsics.ProcessorOptions{NoProcess: false, EvaluateConditions: true}

		if input.Ext() == ".json" {
			template, err = goformation.ParseJSONWithOptions(data, &options)
		} else {
			template, err = goformation.ParseYAMLWithOptions(data, &options)
		}
		if err != nil {
			return DetectedUnknown, err
		}
		if template.Resources == nil {
			return DetectedUnknown, errors.New("Does not contain resources")
		}

		return DetectedAwsCloudformation, nil
	}

	return DetectedUnknown, nil
}
