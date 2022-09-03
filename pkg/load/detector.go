package load

type Detector interface {
	Detect(input Input) (DetectedType, error)
}

var allDetectors []Detector

func detectTypesOfInput(input Input) []DetectedType {
	var detectedTypes []DetectedType

	for _, detector := range allDetectors {
		detectedType, _ := detector.Detect(input)
		if detectedType != DetectedUnknown {
			detectedTypes = append(detectedTypes, detectedType)
		}
	}

	return detectedTypes
}
