### [ Added | Changed | Removed ]


## [unreleased]
### Changed
- reduced memory limit to 50Mb

## [1.6] - 2018-03-29
### Added
- adding the label "force-delete" to a snapshot will include it in the next delete run

## [1.7] - 2018-11-29
### Changed
- moved to using golang modules
- added depency for libc6-compat in the runtime docker image. I believe it is require by the net packages 
