package utils

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"go.uber.org/zap"
)

const region = "us-west-1"
const namespace = "SomeNameSpace"

func NewClient(region string) (*cloudwatch.CloudWatch, error) {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(region),
		},
	)

	if err != nil {
		klog.LOGGER.Debug(
			"Create CloudWatch Client Failed",
			zap.String("TimeStamp", klog.TimeStamp()),
			zap.String("Region", region),
			zap.Error(err),
		)

		return nil, errors.New("Create CloudWatch Client Failed")
	}

	return cloudwatch.New(sess, aws.NewConfig().WithRegion(region)), nil
}

func PublishSimple(c *cloudwatch.CloudWatch, name string, value float64) error {
	m := newMetric(name, value)

	ml := []*cloudwatch.MetricDatum{m}

	_, err := publishMetric(c, ml)

	if err != nil {
		klog.LOGGER.Debug(
			"PublishMetricFailed",
			zap.String("TimeStamp", klog.TimeStamp()),
			zap.String("Region", region),
			zap.Error(err),
		)

		return errors.New(
			fmt.Sprintf("Publish %s metric failed", name),
		)
	}

	return nil
}
