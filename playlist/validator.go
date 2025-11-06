package playlist

import (
	"errors"
	"strconv"
	"time"

	"github.com/globocom/go-m3u8/internal"
)

var (
	ErrMissingPlannedDuration      = errors.New("missing planned duration")
	ErrPlannedDurationExceedsLimit = errors.New("planned duration exceeds 10 minutes")
	offset                         = 1 * time.Second
)

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
	sameDuration := adBreak.HLSElement.Attrs["PLANNED-DURATION"] == previousAdBreak.HLSElement.Attrs["PLANNED-DURATION"]

	// We consider the start date equal if it is within one second of difference.
	startDate, err1 := time.Parse(time.RFC3339Nano, adBreak.HLSElement.Attrs["START-DATE"])
	startDatePrevious, err2 := time.Parse(time.RFC3339Nano, previousAdBreak.HLSElement.Attrs["START-DATE"])
	if err1 != nil || err2 != nil {
		return false
	}

	sameStartDate := startDate.Sub(startDatePrevious) < offset

	return sameStartDate && sameDuration
}

func (p *Playlist) IsCueOutFromBreak(breakMediaSequence int, currentElement *internal.Node) bool {
	nextSegment := p.FindNextSegment(currentElement)

	if nextSegment != nil {
		currentMediaSequence, _ := strconv.Atoi(nextSegment.HLSElement.Details["MediaSequence"])
		if breakMediaSequence == currentMediaSequence {
			return true
		}
	}
	return false
}

func (p *Playlist) IsCueInFromBreak(breakMediaSequence int, currentElement *internal.Node) bool {
	previousSegment := p.FindPreviousSegment(currentElement)

	if previousSegment != nil {
		node, found := p.FindNodeInsideAdBreak(previousSegment)
		if found {
			nodeMediaSequence, _ := strconv.Atoi(node.HLSElement.Details["StartMediaSequence"])
			if nodeMediaSequence == breakMediaSequence {
				return true
			}
		}
	}
	return false
}
