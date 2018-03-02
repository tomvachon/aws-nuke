package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SqsQueue struct {
	svc          *sqs.SQS
	queueurl		 *string
}

func init() {
	register("SqsQueue", ListSqsQueues)
}

func ListSqsQueues(sess *session.Session) ([]Resource, error) {
	svc := sqs.New(sess)

	params := &sqs.ListQueuesInput{}
	resp, err := svc.ListQueues(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, queue := range resp.QueueUrls {
		resources = append(resources, &SqsQueue{
			svc:          svc,
			queueurl: queue,
		})
	}

	return resources, nil
}

func (f *SqsQueue) Remove() error {

	_, err := f.svc.DeleteQueue(&sqs.DeleteQueueInput{
		QueueUrl: f.queueurl,
	})

	return err
}

func (f *SqsQueue) String() string {
	return *f.queueurl
}
