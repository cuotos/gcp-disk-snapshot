package snapshot

import (
	"github.com/cuotos/gcp-disk-snapshot/globalconfig"
	"google.golang.org/api/compute/v1"
	"log"
	"time"
)

const (
	DateFormat = `20060102-1504`
)

func CleanUpSnapshots(client *compute.Service, gcpProjectId string, monthsToKeep int) error {

	snapshots, err := GetAllAutoGeneratedSnapshots(client, gcpProjectId)
	if err != nil {
		log.Fatalf("failed to get list of snapshots from gcp. %v", err)
	}

	oldSnapshots := filterOldSnapshots(snapshots, time.Now().AddDate(0, -monthsToKeep, 0))

	for _, ss := range oldSnapshots {
		err := deleteSnapshot(client, gcpProjectId, ss)
		if err != nil {
			log.Print(err)
		}
	}

	return nil
}

// FilterOldSnapshots accepts a slice of Snapshot Pointers and a time.Time and will deleted snapshots older than this
func filterOldSnapshots(allSnapshots []*compute.Snapshot, threshold time.Time) []*compute.Snapshot {
	var snapshotsToDelete []*compute.Snapshot

	for _, s := range allSnapshots {
		createdDate, _ := time.Parse(time.RFC3339Nano, s.CreationTimestamp)
		if createdDate.Before(threshold) {
			snapshotsToDelete = append(snapshotsToDelete, s)
		}
	}

	return snapshotsToDelete
}

func deleteSnapshot(client *compute.Service, gcpProjectId string, snapshot *compute.Snapshot) error {

	if globalconfig.DryRun {
		log.Printf("DRYRUN deleted snapshot %v", snapshot.Name)
	} else {
		_, err := client.Snapshots.Delete(gcpProjectId, snapshot.Name).Do()
		if err != nil {
			return err
		}
		log.Printf("deleted old snapshot %v\n", snapshot.Name)
	}

	return nil
}

func GetAllAutoGeneratedSnapshots(client *compute.Service, projectName string) ([]*compute.Snapshot, error) {
	snapshotList, err := client.Snapshots.List(projectName).Filter("labels.auto-snapshot eq true").Do()
	if err != nil {
		return nil, err
	}

	return snapshotList.Items, nil
}