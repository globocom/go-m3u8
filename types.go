package go_m3u8

import (
	"time"
)

const (
	m3u8IdentifierTag     = "#EXTM3U"
	versionTag            = "#EXT-X-VERSION"
	targetDurationTag     = "#EXT-X-TARGETDURATION"
	mediaSequenceTag      = "#EXT-X-MEDIA-SEQUENCE"
	programDateTimeTag    = "#EXT-X-PROGRAM-DATE-TIME"
	streamInfTag          = "#EXT-X-STREAM-INF"
	extInfTag             = "#EXTINF"
	dateRangeTag          = "#EXT-X-DATERANGE"
	independentSegmentTag = "#EXT-X-INDEPENDENT-SEGMENTS"
	discontinuityTag      = "#EXT-X-DISCONTINUITY"
	uspTimestampMapTag    = "#USP-X-TIMESTAMP-MAP"
	cueOutTag             = "#EXT-X-CUE-OUT"
	cueInTag              = "#EXT-X-CUE-IN"
	scteOutAttribute      = "SCTE35-OUT"
)

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
