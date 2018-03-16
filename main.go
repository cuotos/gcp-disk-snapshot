package main

import (
	"flag"
	"github.com/cuotos/gcp-disk-snapshot/client"
	"github.com/cuotos/gcp-disk-snapshot/disks"
	"github.com/cuotos/gcp-disk-snapshot/globalconfig"
	"github.com/cuotos/gcp-disk-snapshot/snapshot"
	"google.golang.org/api/compute/v1"
	"log"
	"time"
)

var (
	gcpProjectId                 string
	monthsWorthOfSnapshotsToKeep int
	interval                     time.Duration
)

func init() {

	flag.StringVar(&gcpProjectId, "gcpprojectid", "", "the gcp project id")
	flag.IntVar(&monthsWorthOfSnapshotsToKeep, "months", 6, "number of months worth of snapshots to keep")
	dryrun := flag.Bool("dryrun", false, "dryrun")
	flag.DurationVar(&interval, "interval", time.Duration(time.Hour*24), "run interval")

	flag.Parse()

	if gcpProjectId == "" {
		log.Println("gcpprojectid not provided, attempt to establish from gcp metadata")
		retrievedProjectId, err := client.EstablishProjectId()
		if err != nil {
			log.Fatalf("unable to esablish project id. %v", err)
		}
		gcpProjectId = retrievedProjectId
	}

	globalconfig.DryRun = *dryrun
}

func main() {
	client, err := client.NewComputeClient()
	if err != nil {
		log.Fatal(err)
	}

	//Ticker doesn't run the first one
	doit(client)

	for range time.Tick(interval) {
		doit(client)
	}
}

func doit(client *compute.Service) {
	err := snapshot.CleanUpSnapshots(client, gcpProjectId, monthsWorthOfSnapshotsToKeep)
	if err != nil {
		log.Println(err)
	}

	err = disks.SnapshotDisks(client, gcpProjectId)
	if err != nil {
		log.Println(err)
	}
}
