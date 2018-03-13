package disks

import (
	"fmt"
	"google.golang.org/api/compute/v1"
	"strings"
	"time"
	"log"
	"gitlab.platformserviceaccount.com/lush-soa/dev-ops/gcp-disk-snapshot/service/globalconfig"
)

const (
	SnapshotDateFormat = `20060102-1504`
	snapshotDescription = `created by the go snapshotter`

	//Snapshot name can be not more that 62 characters, so allow for date and hyphen
	maxDiskNameLength  = 62 - len(SnapshotDateFormat) - 1
)

func SnapshotDisks(client *compute.Service, projectId string) error {

	disksToSnapshot, err := getDisksToSnapshot(client, projectId)
	if err != nil {
		return err
	}

	for _, disk := range disksToSnapshot {
		err := snapshotDisk(client, disk, projectId)
		if err != nil {
			log.Printf("failed to snapshot disk. %v", err)
		}
	}
	return nil
}

func getDisksToSnapshot(client *compute.Service, projectId string) ([]*compute.Disk, error) {
	var disksToSnapshot []*compute.Disk

	allDisks, err := client.Disks.AggregatedList(projectId).Filter("labels.snapshot eq true").Do()
	if err != nil {
		return nil, err
	}

	for _, scope := range allDisks.Items {
		if scope.Warning == nil {
			disksToSnapshot = append(disksToSnapshot, scope.Disks...)
		}
	}

	return disksToSnapshot, nil
}

func createSnapshotName(disk *compute.Disk, timestamp time.Time) string {
	dateString := timestamp.Format(SnapshotDateFormat)
	diskName := ""

	if disk.Labels["name"] != "" {
		diskName = disk.Labels["name"]
	} else {
		diskName = disk.Name
	}

	if len(diskName) > 48 {
		truncatedNameSlice := strings.Split(diskName, "")[:maxDiskNameLength]
		diskName = strings.Join(truncatedNameSlice, "")
	}

	return fmt.Sprintf("%v-%v", diskName, dateString)
}

func snapshotDisk(client *compute.Service, disk *compute.Disk, projectId string) error {
	snapshotName := createSnapshotName(disk, time.Now())
	snapshot := &compute.Snapshot{}
	snapshot.Name = snapshotName
	snapshot.Description = snapshotDescription
	snapshot.Labels = map[string]string{"auto-snapshot":"true"}

	// disk.Zone contains the full url of the zone, not just the name
	splitZoneUrl := strings.Split(disk.Zone, "/")
	zoneName := splitZoneUrl[len(splitZoneUrl)-1:][0]

	createSnapshotReq := client.Disks.CreateSnapshot(projectId, zoneName, disk.Name, snapshot)

	if globalconfig.DryRun {
		log.Printf("DRYRUN snapshot %v", snapshotName)
	} else {
		_, err := createSnapshotReq.Do()
		if err != nil {
			return err
		}

		log.Printf("triggered snapshot of %v\n", snapshotName)
	}
	return nil
}

