package disks

import (
	"testing"
	"time"

	"google.golang.org/api/compute/v1"
)

const (
	refDateString = `20180102-1354`
)

func TestCreateSnapshotName(t *testing.T) {
	tcs := []struct {
		Name     string
		Disk     *compute.Disk
		Expected string
	}{
		{
			"valid date and name",
			&compute.Disk{Name: "test-disk"},
			"test-disk-20180102-1354",
		},
		{
			"really long disk name",
			&compute.Disk{Name: "disk-name-is-really-long-and-should-get-truncated-to-something-shorter"},
			"disk-name-is-really-long-and-should-get-truncate-20180102-1354",
		},
		{
			"short disk name should be fine",
			&compute.Disk{Name: "the-disk"},
			"the-disk-20180102-1354",
		},
		{"disk with a short name label",
			&compute.Disk{Name: "dontcare", Labels: map[string]string{"name": "name-label"}},
			"name-label-20180102-1354",
		},
		{"disk with a really long name label",
			&compute.Disk{Name: "dontcare", Labels: map[string]string{"name": "disk-name-label-is-really-long-and-should-get-truncated-to-something-shorter"}},
			"disk-name-label-is-really-long-and-should-get-tr-20180102-1354",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			diskSnapshotDate, err := time.Parse(SnapshotDateFormat, refDateString)
			if err != nil {
				t.Fatal(err)
			}
			actual := createSnapshotName(tc.Disk, diskSnapshotDate)
			if actual != tc.Expected {
				t.Errorf("got %v but expected %v", actual, tc.Expected)
			}
		})
	}
}
