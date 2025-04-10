package internal

import (
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// The HLSElement data type holds the following attributes:
//   - Name: The name of the Element (e.g. tag name).
//   - URI: The Uniform Resource Identifier of the Element (if applicable).
//   - Attrs: In-manifest Element attributes, in key-value format.
//   - Details: Not-in-manifest Element attributes, in key-value format.
type HLSElement struct {
	Name    string
	URI     string
	Attrs   map[string]string
	Details map[string]string
}

type StreamInf struct {
	Codecs           []string
	Bandwidth        string
	AverageBandwidth string
	Resolution       string
	FrameRate        string
	URI              string
}

type Segment struct {
	Duration        float64
	ProgramDateTime time.Time
	MediaSequence   int
	URI             string
}

type DateRange struct {
	ID              string
	Class           string
	StartDate       time.Time
	EndDate         time.Time
	PlannedDuration float64
	Scte35Out       string
	Scte35In        string
	MediaSequence   int
}

func (e *HLSElement) ToDateRangeType(playlistMediaSequence, playlistSegmentsCounter int) *DateRange {
	startDate, _ := time.Parse(time.RFC3339Nano, e.Attrs["START-DATE"])
	endDate, _ := time.Parse(time.RFC3339Nano, e.Attrs["END-DATE"])
	duration, _ := strconv.ParseFloat(e.Attrs["DURATION"], 64)

	return &DateRange{
		ID:              e.Attrs["ID"],
		Class:           e.Attrs["CLASS"],
		StartDate:       startDate,
		EndDate:         endDate,
		PlannedDuration: duration,
		Scte35Out:       e.Attrs["SCTE35-OUT"],
		Scte35In:        e.Attrs["SCTE35-IN"],
		MediaSequence:   playlistMediaSequence + playlistSegmentsCounter,
	}
}

func ToStreamInfType(mappedAttr map[string]string) *StreamInf {
	return &StreamInf{
		Bandwidth:        mappedAttr["BANDWIDTH"],
		AverageBandwidth: mappedAttr["AVERAGE-BANDWIDTH"],
		Codecs:           strings.Split(mappedAttr["CODECS"], ","),
		Resolution:       mappedAttr["RESOLUTION"],
		FrameRate:        mappedAttr["FRAME-RATE"],
	}
}

func ToSegmentType(duration string, playlistMediaSequence, playlistSegmentsCounter int, playlistDVR float64, playlistPDT time.Time) *Segment {
	floatDuration, err := strconv.ParseFloat(duration, 64)
	if err != nil {
		log.Error().Err(err).Msgf("failed to parse duration for segment: %s", duration)
		return &Segment{}
	}

	currentDVRInNanoseconds := int(playlistDVR * float64(time.Second))
	segmentProgramDateTime := playlistPDT.Add(time.Duration(currentDVRInNanoseconds))

	return &Segment{
		Duration:        floatDuration,
		MediaSequence:   playlistMediaSequence + playlistSegmentsCounter,
		ProgramDateTime: segmentProgramDateTime,
	}
}
