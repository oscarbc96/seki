package aws_cloudformation_authentication

import (
	"fmt"
	"github.com/awslabs/goformation/v6"
	"github.com/oscarbc96/seki/pkg/analyser"
	"github.com/oscarbc96/seki/pkg/result"
)

func init() {
	analyser.Analysers["a"] = Run
}
func Run() (*result.RuleResult, error) {
	template, err := goformation.Open("test.json")
	if err != nil {
		return nil, err
	}

	for _, resource := range template.GetAllEC2InstanceResources() {
		if resource.AWSCloudFormationMetadata != nil {
			if authentication, ok := resource.AWSCloudFormationMetadata["AWS::CloudFormation::Authentication"]; ok {
				for _, authenticationBlock := range authentication.(map[string]interface{}) {
					for key, value := range authenticationBlock.(map[string]interface{}) {
						fmt.Println("Key:", key, "=>", "Element:", value)
					}
					//for auth in self.Metadata["AWS::CloudFormation::Authentication"].values():
					//	if not all(
					//		[
					//			auth.get("accessKeyId", Parameter.NO_ECHO_NO_DEFAULT) == Parameter.NO_ECHO_NO_DEFAULT,
					//			auth.get("password", Parameter.NO_ECHO_NO_DEFAULT) == Parameter.NO_ECHO_NO_DEFAULT,
					//			auth.get("secretKey", Parameter.NO_ECHO_NO_DEFAULT) == Parameter.NO_ECHO_NO_DEFAULT,
					//		]
					//	):
					//		return True
					return &result.RuleResult{Result: result.FAIL, Severity: result.Medium}, nil
				}
			}

		}
	}

	return &result.RuleResult{Result: result.PASS, Severity: result.Medium}, nil
}
