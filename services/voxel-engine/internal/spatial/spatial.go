package spatial

import "github.com/uber/h3-go/v4"

// Let's use a default resolution of 8, which is a good balance between
// accuracy and performance for most applications. This resolution corresponds
// to hexagons with an average edge length of about 1.2 km, which is suitable
// for many spatial analyses and visualizations.
const DefaultResolution = 8

// PointToH3Cell implements Eq.: H3_cell = h3.latlng_to_cell(lat, lon, res).
func PointToH3Cell(lat, lon float64, resolution int) (h3.Cell, error) {
	return h3.LatLngToCell(h3.NewLatLng(lat, lon), resolution)
}

// NeighborCells returns the k-ring: the cell itself plus its neighbours out
// to grid distance k. k=1 always returns exactly 7 cells (itself + 6
// neighbours) for a non-pentagon cell.
func NeighborCells(cell h3.Cell, k int) ([]h3.Cell, error) {
	return h3.GridDisk(cell, k)
}
