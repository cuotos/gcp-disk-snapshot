package snapshot

import (
	"google.golang.org/api/compute/v1"
	"testing"
	"time"
)

const (
	refDateString = `20180102-1354`
)

var (
	refDate time.Time
)

func TestFilterOldSnapshots(t *testing.T) {
	tcs := []struct {
		Name                    string
		AllSnapshots            []*compute.Snapshot
		ExpectedNumberOfResults int
	}{
		{
			"All old snapshots",
			createTestSnapshots(0, 5),
			5,
		},
		{
			"Some snapshots older than threshold",
			createTestSnapshots(2, 3),
			3,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			refDateString, _ := time.Parse(DateFormat, refDateString)
			filteredSnapshots := filterOldSnapshots(tc.AllSnapshots, refDateString)

			if len(filteredSnapshots) != tc.ExpectedNumberOfResults {
				t.Errorf("filtered list contains %v entries but should only have %v", len(filteredSnapshots), tc.ExpectedNumberOfResults)
			}
		})
	}
}

func createTestSnapshots(currentSnapshots int, oldSnapshots int) []*compute.Snapshot {
	var snapshots []*compute.Snapshot

	validDate, _ := time.Parse(DateFormat, refDateString)
	oldDate := validDate.AddDate(0, 0, -10)

	for i := 0; i < currentSnapshots; i++ {
		snapshots = append(snapshots, &compute.Snapshot{CreationTimestamp: validDate.Format(time.RFC3339Nano)})
	}

	for i := 0; i < oldSnapshots; i++ {
		snapshots = append(snapshots, &compute.Snapshot{CreationTimestamp: oldDate.Format(time.RFC3339Nano)})
	}

	return snapshots
}
