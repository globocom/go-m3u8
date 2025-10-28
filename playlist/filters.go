package playlist

import (
	"strconv"
	"strings"
	"time"

	"github.com/globocom/go-m3u8/internal"
)

// Removes all variant streams (#EXT-X-STREAM-INF) from the playlist that exceed the given maxHeight.
func (p *Playlist) FilterByMaxHeight(maxHeight int) {
	nodes := p.Variants()
	for _, node := range nodes {
		resolution := node.HLSElement.Attrs["RESOLUTION"]
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
			p.Remove(node)
		}
	}
}

// Validates START-DATE attribute of the given ad break
func invalidStartDate(adBreak *internal.Node) bool {
	_, err := time.Parse(time.RFC3339Nano, adBreak.HLSElement.Attrs["START-DATE"])
	if err != nil {
		return true
	}
	return false
}

// Validates StartMediaSequence of the given ad break
func invalidMediaSequence(adBreak *internal.Node) bool {
	_, err := strconv.Atoi(adBreak.HLSElement.Details["StartMediaSequence"])
	if err != nil {
		return true
	}
	return false
}

// Validates PLANNED-DURATION attribute of the given ad break
func invalidPlannedDuration(adBreak *internal.Node) bool {
	plannedDurationStr := adBreak.HLSElement.Attrs["PLANNED-DURATION"]
	plannedDuration, err := strconv.ParseFloat(plannedDurationStr, 64)
	const maxPlannedDuration = 600 // 10 minutes

	if plannedDuration == 0 || err != nil || plannedDuration > maxPlannedDuration {
		return true
	}
	return false
}

// Checks if the given ad break is a duplicate of the previous ad break
func duplicatedBreak(adBreak, previousAdBreak *internal.Node) bool {
	sameStartDate := adBreak.HLSElement.Attrs["START-DATE"] == previousAdBreak.HLSElement.Attrs["START-DATE"]
	sameDuration := adBreak.HLSElement.Attrs["PLANNED-DURATION"] == previousAdBreak.HLSElement.Attrs["PLANNED-DURATION"]
	if sameStartDate && sameDuration {
		return true
	}
	return false
}
