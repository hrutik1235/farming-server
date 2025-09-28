package types

type CreateCrop struct {
	Name            string  `json:"name" validate:"required" name:"name"`
	BasePrice       float64 `json:"base_price" validate:"required" name:"base_price"`
	GrowthTimeHours int     `json:"growth_time_hours" validate:"required" name:"growth_time_hours"`
	YieldPerUnit    int     `json:"yield_per_unit" validate:"required" name:"yield_per_unit"`
	CostPerUnit     float64 `json:"cost_per_unit" validate:"required" name:"cost_per_unit"`
	Description     string  `json:"description" validate:"required" name:"description"`
}

type PlantCrop struct {
	LandUnits int `json:"land_units" validate:"required" name:"land_units"`
}
