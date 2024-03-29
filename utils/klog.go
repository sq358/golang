package klog

import (
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// For some users, the presets offered by the NewProduction, NewDevelopment,
// and NewExample constructors won't be appropriate. For most of those
// users, the bundled Config struct offers the right balance of flexibility
// and convenience. (For more complex needs, see the AdvancedConfiguration
// example.)
//
// See the documentation for Config and zapcore.EncoderConfig for all the
// available options.

var Debug bool
var LOGGER *zap.Logger
var TimeFormat = "2006-01-02T15:04:05-0700"

var err error

func Init(logfile string) {
	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logfile,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.DebugLevel,
	)

	LOGGER = zap.New(core)
}

func TimeStamp() string {
	return time.Now().UTC().Format(TimeFormat)
}

func publishMetric(c *cloudwatch.CloudWatch, ml []*cloudwatch.MetricDatum) (*cloudwatch.PutMetricDataOutput, error) {
	i := cloudwatch.PutMetricDataInput{}

	i.SetMetricData(ml)
	i.SetNamespace(namespace)

	out, err := c.PutMetricData(&i)

	if err != nil {
		klog.LOGGER.Debug(
			"Publish Metric Data Failed",
			zap.String("TimeStamp", klog.TimeStamp()),
			zap.String(
				"Metric Input",
				fmt.Sprintf("%+v", i),
			),
			zap.Error(err),
		)

		return nil, errors.New("Publish Metric to CloudWatch Failed")
	}

	klog.LOGGER.Info(
		"Publish Metric Data Succeeded",
		zap.String("TimeStamp", klog.TimeStamp()),
		zap.String(
			"Output",
			fmt.Sprintf("%+v", *out),
		),
		zap.Error(err),
	)

	return out, nil
}
