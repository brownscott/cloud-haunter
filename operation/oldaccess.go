package operation

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/hortonworks/cloud-cost-reducer/context"
	"github.com/hortonworks/cloud-cost-reducer/types"
)

var defaultAvailablePeriod = 120 * 24 * time.Hour

type oldAccess struct {
	availablePeriod time.Duration
}

func init() {
	availableEnv := os.Getenv("ACCESS_AVAILABLE_PERIOD")
	var availablePeriod time.Duration
	if len(availableEnv) > 0 {
		if duration, err := time.ParseDuration(availableEnv); err != nil {
			log.Errorf("[OLDACCESS] err: %s", err)
			return
		} else {
			availablePeriod = duration
		}
	} else {
		availablePeriod = defaultAvailablePeriod
	}
	log.Infof("[OLDACCESS] running period set to: %s", availablePeriod)
	context.Operations[types.OLDACCESS] = oldAccess{availablePeriod}
}

func (o oldAccess) Execute(clouds []types.CloudType) []types.CloudItem {
	if context.DryRun {
		log.Debugf("Collecting old accesses on: [%s]", clouds)
	}
	accessChan, errChan := o.collect(clouds)
	items := wait(accessChan, errChan, "[OLDACCESS] Failed to collect old accesses")
	return o.filter(items)
}

func (o oldAccess) filter(items []types.CloudItem) []types.CloudItem {
	if context.DryRun {
		log.Debugf("Filtering accesses (%d): [%s]", len(items), items)
	}
	return filter(items, func(item types.CloudItem) bool {
		match := item.GetCreated().Add(o.availablePeriod).Before(time.Now())
		if context.DryRun {
			log.Debugf("Access: %s match: %b", item.GetName(), match)
		}
		return match
	})
}

func (o oldAccess) collect(clouds []types.CloudType) (chan []types.CloudItem, chan error) {
	return collect(clouds, func(provider types.CloudProvider) ([]types.CloudItem, error) {
		accesses, err := provider.GetAccesses()
		if err != nil {
			return nil, err
		}
		return o.convertToCloudItems(accesses), nil
	})
}

func (o oldAccess) convertToCloudItems(accesses []*types.Access) []types.CloudItem {
	items := []types.CloudItem{}
	for _, access := range accesses {
		items = append(items, access)
	}
	return items
}
