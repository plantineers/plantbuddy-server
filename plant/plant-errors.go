package plant

import "errors"

var ErrPlantGroupNotExisting = errors.New("plant group does not exist")
var ErrPlantGroupStillInUse = errors.New("cannot delete plant group because it is still in use")
