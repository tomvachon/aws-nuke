package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type SNSEndpoint struct {
	svc *sns.SNS
	ARN *string
}

func init() {
	register("SNSEndpoint", ListSNSEndpoints)
}

func ListSNSEndpoints(sess *session.Session) ([]Resource, error) {
	svc := sns.New(sess)
	resources := []Resource{}
	applicationPlatforms := []*sns.PlatformApplication{}

	appParams := &sns.ListPlatformApplicationsInput{}

	for {
		resp, err := svc.ListPlatformApplications(appParams)
		if err != nil {
			return nil, err
		}

		for _, platformApplication := range resp.PlatformApplications {
			applicationPlatforms = append(applicationPlatforms, platformApplication)
		}
		if resp.NextToken == nil {
			break
		}

		appParams.NextToken = resp.NextToken
	}

	params := &sns.ListEndpointsByPlatformApplicationInput{}

	for _, platformApplication := range applicationPlatforms {
		params.PlatformApplicationArn = platformApplication.PlatformApplicationArn
		resp, err := svc.ListEndpointsByPlatformApplication(params)
		if err != nil {
			return nil, err
		}

		for _, endpoint := range resp.Endpoints {
			resources = append(resources, &SNSEndpoint{
				svc: svc,
				ARN: endpoint.EndpointArn,
			})
		}
		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}
	return resources, nil
}

func (f *SNSEndpoint) Remove() error {

	_, err := f.svc.DeleteEndpoint(&sns.DeleteEndpointInput{
		EndpointArn: f.ARN,
	})

	return err
}

func (f *SNSEndpoint) String() string {
	return *f.ARN
}
