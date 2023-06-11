package plant

import "errors"

// ErrPlantGroupNotExisting is returned when a plant group does not exist.
var ErrPlantGroupNotExisting = errors.New("plant group does not exist")

// ErrPlantGroupStillInUse is returned when a plant group is about to be deleted but some plants still use it.
var ErrPlantGroupStillInUse = errors.New("cannot delete plant group because it is still in use")
