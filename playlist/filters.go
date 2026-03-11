package playlist

import (
	"strconv"
	"strings"
	"time"
)

// Parses a resolution string in the format "WIDTHxHEIGHT"
// returns the width, height, and a boolean indicating success.
func parseResolution(resolution string) (int, int, bool) {
	if resolution == "" {
		return 0, 0, false
	}

	parts := strings.Split(resolution, "x")
	if len(parts) != 2 {
		return 0, 0, false
	}

	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, false
	}

	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, false
	}

	return width, height, true
}

// Parses an integer attribute from a string
// returning the parsed value and a boolean indicating success.
func parseIntAttr(value string) (int, bool) {
	if value == "" {
		return 0, false
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, false
	}

	return parsed, true
}

// Parses a float attribute from a string
// returning the parsed value and a boolean indicating success.
func parseFloatAttr(value string) (float64, bool) {
	if value == "" {
		return 0, false
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, false
	}

	return parsed, true
}

// Master Manifest filters
// Removes all variant streams (#EXT-X-STREAM-INF) from the playlist that exceed the given maxHeight.
func (p *Playlist) FilterByMaxHeight(maxHeight int) {
	nodes := p.Variants()
	for i := range nodes {
		_, height, ok := parseResolution(nodes[i].HLSElement.Attrs["RESOLUTION"])
		if !ok {
			continue
		}
		if height > maxHeight {
			p.Remove(nodes[i])
		}
	}
}

// Removes all variant streams (#EXT-X-STREAM-INF) from the playlist that are below the given minHeight.
func (p *Playlist) FilterByMinHeight(minHeight int) {
	nodes := p.Variants()
	for i := range nodes {
		_, height, ok := parseResolution(nodes[i].HLSElement.Attrs["RESOLUTION"])
		if !ok {
			continue
		}
		if height < minHeight {
			p.Remove(nodes[i])
		}
	}
}

// Removes all variant streams (#EXT-X-STREAM-INF) from the playlist that exceed the given maxWidth.
func (p *Playlist) FilterByMaxWidth(maxWidth int) {
	nodes := p.Variants()
	for i := range nodes {
		width, _, ok := parseResolution(nodes[i].HLSElement.Attrs["RESOLUTION"])
		if !ok {
			continue
		}
		if width > maxWidth {
			p.Remove(nodes[i])
		}
	}
}

// Removes all variant streams (#EXT-X-STREAM-INF) from the playlist that are below the given minWidth.
func (p *Playlist) FilterByMinWidth(minWidth int) {
	nodes := p.Variants()
	for i := range nodes {
		width, _, ok := parseResolution(nodes[i].HLSElement.Attrs["RESOLUTION"])
		if !ok {
			continue
		}
		if width < minWidth {
			p.Remove(nodes[i])
		}
	}
}

// Removes all variant streams (#EXT-X-STREAM-INF) from the playlist that exceed the given maxBitrate.
func (p *Playlist) FilterByMaxBitrate(maxBitrate int) {
	nodes := p.Variants()
	for i := range nodes {
		bitrate, ok := parseIntAttr(nodes[i].HLSElement.Attrs["BANDWIDTH"])
		if !ok {
			continue
		}
		if bitrate > maxBitrate {
			p.Remove(nodes[i])
		}
	}
}

// Removes all variant streams (#EXT-X-STREAM-INF) from the playlist that are below the given minBitrate.
func (p *Playlist) FilterByMinBitrate(minBitrate int) {
	nodes := p.Variants()
	for i := range nodes {
		bitrate, ok := parseIntAttr(nodes[i].HLSElement.Attrs["BANDWIDTH"])
		if !ok {
			continue
		}
		if bitrate < minBitrate {
			p.Remove(nodes[i])
		}
	}
}

// Media Manifest filters
// Removes older media timeline entries to keep only the most recent dvrWindowSeconds of segments.
func (p *Playlist) FilterByDVRWindow(dvrWindowSeconds float64) {
	if dvrWindowSeconds <= 0 {
		return
	}

	segments := p.Segments()
	if len(segments) == 0 {
		return
	}

	keptStartIdx := len(segments) - 1
	cumulativeDuration := 0.0

	for i := len(segments) - 1; i >= 0; i-- {
		duration, ok := parseFloatAttr(segments[i].HLSElement.Attrs["Duration"])
		if !ok {
			continue
		}

		if cumulativeDuration+duration > dvrWindowSeconds {
			break
		}

		cumulativeDuration = RoundFloat(cumulativeDuration+duration, 4)
		keptStartIdx = i
	}

	if keptStartIdx == 0 {
		return
	}

	oldestSegment := segments[0]
	newOldestSegment := segments[keptStartIdx]

	current := oldestSegment
	for current != nil && current != newOldestSegment {
		next := current.Next
		p.Remove(current)
		current = next
	}

	if mediaSequence, ok := parseIntAttr(newOldestSegment.HLSElement.Details["MediaSequence"]); ok {
		if mediaSequenceTag, found := p.MediaSequenceTag(); found {
			mediaSequenceTag.HLSElement.Attrs["#EXT-X-MEDIA-SEQUENCE"] = strconv.Itoa(mediaSequence)
		}
		p.MediaSequence = mediaSequence
	}

	if segmentPDT, err := time.Parse(time.RFC3339Nano, newOldestSegment.HLSElement.Details["ProgramDateTime"]); err == nil {
		if pdtTag, found := p.Find("ProgramDateTime"); found {
			pdtTag.HLSElement.Attrs["#EXT-X-PROGRAM-DATE-TIME"] = segmentPDT.Format(time.RFC3339Nano)
		}
		p.ProgramDateTime = segmentPDT
	}

	updatedSegments := p.Segments()
	updatedDVR := 0.0
	for i := range updatedSegments {
		duration, ok := parseFloatAttr(updatedSegments[i].HLSElement.Attrs["Duration"])
		if !ok {
			continue
		}
		updatedDVR = RoundFloat(updatedDVR+duration, 4)
	}

	p.SegmentsCounter = len(updatedSegments)
	p.DVR = updatedDVR
}
