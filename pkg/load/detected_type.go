package load

import "encoding/json"

type DetectedType int

const (
	DetectedUnknown DetectedType = iota
	DetectedAwsCloudformation
	DetectedContainerDockerfile
	DetectedCICDGitHubActions
)

var detectedTypeToString = map[DetectedType]string{
	DetectedUnknown:             "UNKNOWN",
	DetectedAwsCloudformation:   "cloudformation",
	DetectedContainerDockerfile: "dockerfile",
	DetectedCICDGitHubActions:   "github_actions",
}

func (d DetectedType) String() string {
	return detectedTypeToString[d]
}

func (d DetectedType) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}
