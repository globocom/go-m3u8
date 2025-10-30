package playlist

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/globocom/go-m3u8/internal"
)

var (
	ErrMissingPlannedDuration      = errors.New("missing planned duration")
	ErrPlannedDurationExceedsLimit = errors.New("planned duration exceeds 10 minutes")
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
func ValidateStartDate(adBreak *internal.Node) (time.Time, error) {
	startDate, err := time.Parse(time.RFC3339Nano, adBreak.HLSElement.Attrs["START-DATE"])
	return startDate, err
}

// Validates StartMediaSequence of the given ad break
func ValidateMediaSequence(adBreak *internal.Node) (int, error) {
	startMediaSequence, err := strconv.Atoi(adBreak.HLSElement.Details["StartMediaSequence"])
	return startMediaSequence, err
}

// Validates PLANNED-DURATION attribute of the given ad break
func ValidatePlannedDuration(adBreak *internal.Node) error {
	plannedDurationStr := adBreak.HLSElement.Attrs["PLANNED-DURATION"]
	plannedDuration, err := strconv.ParseFloat(plannedDurationStr, 64)
	const maxPlannedDuration = 600 // 10 minutes

	if err != nil {
		return err
	}

	if plannedDuration == 0 {
		return ErrMissingPlannedDuration
	}

	if plannedDuration > maxPlannedDuration {
		return ErrPlannedDurationExceedsLimit
	}

	return nil
}

// Checks if the given ad break is a duplicate of the previous ad break
func IsDuplicatedBreak(adBreak, previousAdBreak *internal.Node) bool {
	sameStartDate := adBreak.HLSElement.Attrs["START-DATE"] == previousAdBreak.HLSElement.Attrs["START-DATE"]
	sameDuration := adBreak.HLSElement.Attrs["PLANNED-DURATION"] == previousAdBreak.HLSElement.Attrs["PLANNED-DURATION"]
	return sameStartDate && sameDuration
}
