package playlist

import (
	"strconv"
	"strings"
)

// Removes all variant streams (#EXT-X-STREAM-INF) from the playlist that exceed the given maxHeight.
func (p *Playlist) FilterByMaxHeight(maxHeight int) {
	nodes := p.Variants()
	for i := range nodes {
		resolution := nodes[i].HLSElement.Attrs["RESOLUTION"]
		if resolution == "" {
			continue
		}
		parts := strings.Split(resolution, "x")
		if len(parts) != 2 {
			continue
		}
		heightStr := parts[1]
		height, err := strconv.Atoi(heightStr)
		if err != nil {
			continue
		}
		if height > maxHeight {
			p.Remove(nodes[i])
		}
	}
}
