package scheduler

import (
	"context"
	"crab/cluster"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func HealthProbeStatus(name string) error {
	for {
		value, err := cluster.Client.Clientset.BatchV1().Jobs("island-system").Get(context.Background(),
			name, v1.GetOptions{})
		if err != nil {
			return errors.Wrap(err, "Error HealthProbeStatus")
		}
		if value.Status.Active > 0 {
			time.Sleep(time.Second * 2)
			continue
		} else {
			if value.Status.Succeeded > 0 {
				return nil
			}
			return errors.New("health probe fail")
		}
	}
}