package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/serverlessapplicationrepository"
)

type SarApplication struct {
	svc *serverlessapplicationrepository.ServerlessApplicationRepository
	id  string
}

func init() {
	register("SarApplication", ListSarApplications)
}

func ListSarApplications(sess *session.Session) ([]Resource, error) {
	svc := serverlessapplicationrepository.New(sess)

	params := &serverlessapplicationrepository.ListApplicationsInput{}
	resp, err := svc.ListApplications(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.Applications {
		resources = append(resources, &SarApplication{
			svc: svc,
			id:  *out.ApplicationId,
		})
	}

	return resources, nil
}

func (e *SarApplication) Remove() error {
	_, err := e.svc.DeleteApplication(&serverlessapplicationrepository.DeleteApplicationInput{
		ApplicationId: &e.id,
	})
	return err
}

func (e *SarApplication) String() string {
	return e.id
}
