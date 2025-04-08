package internal

import "time"

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
