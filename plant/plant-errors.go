package plant

import "errors"

var ErrPlantNameRequired = errors.New("plant name is required")
var ErrPlantGroupRequired = errors.New("plant group is required")
var ErrPlantGroupNameRequired = errors.New("plant group name is required")
