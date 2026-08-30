package adaptive

import "github.com/uber/h3-go/v4"

type Config struct {
	BaseResolution   int
	FineResolution   int
	BaseTimeBinS     int
	FineTimeBinS     int
	DensityThreshold int // forecast occupancy count that triggers "zoom in"
}

func DefaultConfig() Config {
	return Config{BaseResolution: 8, FineResolution: 9, BaseTimeBinS: 10, FineTimeBinS: 5, DensityThreshold: 6}
}

type Engine struct{ cfg Config }

func New(cfg Config) *Engine { return &Engine{cfg: cfg} }

func (e *Engine) ResolutionFor(forecastOccupancy int) int {
	if forecastOccupancy >= e.cfg.DensityThreshold {
		return e.cfg.FineResolution
	}
	return e.cfg.BaseResolution
}

func (e *Engine) TimeBinFor(forecastOccupancy int) int {
	if forecastOccupancy >= e.cfg.DensityThreshold {
		return e.cfg.FineTimeBinS
	}
	return e.cfg.BaseTimeBinS
}

// RefineCell expands one base-resolution cell into its finer-resolution
// children, so the voxelizer can re-bin points inside a hot cell
// more precisely without changing resolution anywhere else in the city.
func (e *Engine) RefineCell(parent h3.Cell) ([]h3.Cell, error) {
    return parent.Children(e.cfg.FineResolution)
}
