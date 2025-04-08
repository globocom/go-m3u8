package go_m3u8

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/globocom/go-m3u8/internal"
)

var (
	ErrParseLine = errors.New("failed to parse tag")
	paramRegex   = regexp.MustCompile(`([a-zA-Z\d_-]+)=("[^"]+"|[^",]+)`)
)

type (
	m3u8IdentifierParser      struct{}
	versionParser             struct{}
	targetDurationParser      struct{}
	mediaSequenceParser       struct{}
	programDateTimeParser     struct{}
	dateRangeParser           struct{}
	extInfParser              struct{}
	streamInfParser           struct{}
	independentSegmentsParser struct{}
	discontinuityParser       struct{}
	uspTimestampMapParser     struct{}
	cueOutParser              struct{}
	cueInParser               struct{}
)

func (p m3u8IdentifierParser) Parse(tag string, playlist *Playlist) error {
	playlist.Insert(&internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "M3u8Identifier",
			Attrs: map[string]string{
				m3u8IdentifierTag: "",
			},
		},
	})
	return nil
}

func (p versionParser) Parse(tag string, playlist *Playlist) error {
	parts := strings.Split(tag, ":")
	if len(parts) > 1 && parts[1] != "" {
		playlist.Insert(&internal.Node{
			HLSElement: &internal.HLSElement{
				Name:  "Version",
				Attrs: map[string]string{versionTag: strings.TrimSpace(parts[1])},
			},
		})
		return nil
	}
	return fmt.Errorf("%w: invalid version tag", ErrParseLine)
}

func (p targetDurationParser) Parse(tag string, playlist *Playlist) error {
	parts := strings.Split(tag, ":")
	if len(parts) > 1 && parts[1] != "" {
		playlist.Insert(&internal.Node{
			HLSElement: &internal.HLSElement{
				Name:  "TargetDuration",
				Attrs: map[string]string{targetDurationTag: strings.TrimSpace(parts[1])},
			},
		})
		return nil
	}
	return fmt.Errorf("%w: invalid target duration tag", ErrParseLine)
}

func (p mediaSequenceParser) Parse(tag string, playlist *Playlist) error {
	var err error
	parts := strings.Split(tag, ":")
	if len(parts) > 1 && parts[1] != "" {
		mediaSequence := strings.TrimSpace(parts[1])
		playlist.Insert(&internal.Node{
			HLSElement: &internal.HLSElement{
				Name:  "MediaSequence",
				Attrs: map[string]string{mediaSequenceTag: mediaSequence},
			},
		})

		playlist.MediaSequence, err = strconv.Atoi(mediaSequence)
		if err != nil {
			return fmt.Errorf("%w: invalid media sequence number", ErrParseLine)
		}

		return nil
	}
	return fmt.Errorf("%w: invalid media sequence tag", ErrParseLine)
}

func (p programDateTimeParser) Parse(tag string, playlist *Playlist) error {
	parts := strings.SplitN(tag, ":", 2)

	if len(parts) <= 1 {
		return fmt.Errorf("%w: invalid program date time tag", ErrParseLine)
	}

	dateTimeValue := strings.TrimSpace(parts[1])
	playlist.Insert(&internal.Node{
		HLSElement: &internal.HLSElement{
			Name:  "ProgramDateTime",
			Attrs: map[string]string{programDateTimeTag: dateTimeValue},
		},
	})

	if playlist.ProgramDateTime.Format(time.DateOnly) == "0001-01-01" {
		parsedTime, err := time.Parse(time.RFC3339Nano, dateTimeValue)

		if err != nil {
			return fmt.Errorf("%w: invalid program date time tag", ErrParseLine)
		}

		playlist.ProgramDateTime = parsedTime
	}

	return nil
}

func (p dateRangeParser) Parse(tag string, playlist *Playlist) error {
	params := tagsToMap(tag)
	if len(params) < 1 {
		return fmt.Errorf("%w: invalid date range tag", ErrParseLine)
	}

	dateRangeNode := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name:  "DateRange",
			Attrs: params,
			Details: map[string]string{
				"MediaSequence": fmt.Sprintf("%d", playlist.MediaSequence+playlist.SegmentsCounter),
			},
		},
	}

	if dateRangeNode.HLSElement.Attrs["SCTE35-OUT"] != "" {
		playlist.CurrentDateRange = dateRangeNode
	}

	playlist.Insert(dateRangeNode)
	return nil
}

func (p streamInfParser) Parse(tag string, playlist *Playlist) error {
	params := tagsToMap(tag)
	playlist.CurrentStreamInf = &internal.StreamInf{
		Bandwidth:        params["BANDWIDTH"],
		AverageBandwidth: params["AVERAGE-BANDWIDTH"],
		Codecs:           strings.Split(params["CODECS"], ","),
		Resolution:       params["RESOLUTION"],
		FrameRate:        params["FRAME-RATE"],
	}
	return nil
}

func (p extInfParser) Parse(tag string, playlist *Playlist) error {
	parts := strings.Split(tag, ":")
	if len(parts) > 1 {
		duration := strings.TrimSpace(strings.Split(parts[1], ",")[0])
		floatDuration, err := strconv.ParseFloat(duration, 64)
		if err != nil {
			return fmt.Errorf("%w: invalid duration tag", ErrParseLine)
		}

		currentDVRInNanoseconds := int(playlist.DVR * float64(time.Second))
		segmentProgramDateTime := playlist.ProgramDateTime.Add(time.Duration(currentDVRInNanoseconds))

		playlist.CurrentSegment = &internal.Segment{
			Duration:        floatDuration,
			MediaSequence:   playlist.MediaSequence + playlist.SegmentsCounter,
			ProgramDateTime: segmentProgramDateTime,
		}

		playlist.DVR = roundFloat(playlist.DVR+floatDuration, 4)
		playlist.SegmentsCounter += 1

		return nil
	}
	return fmt.Errorf("%w: invalid extension tag", ErrParseLine)
}

func (p independentSegmentsParser) Parse(tag string, playlist *Playlist) error {
	playlist.Insert(&internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "IndependentSegments",
			Attrs: map[string]string{
				independentSegmentTag: "",
			},
		},
	})
	return nil
}

func (p discontinuityParser) Parse(tag string, playlist *Playlist) error {
	playlist.Insert(&internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "Discontinuity",
			Attrs: map[string]string{
				discontinuityTag: "",
			},
		},
	})
	return nil
}

func (p uspTimestampMapParser) Parse(tag string, playlist *Playlist) error {
	parts := strings.SplitN(tag, ":", 2)
	if len(parts) > 0 {
		params := tagsToMap(parts[1])
		playlist.Insert(&internal.Node{
			HLSElement: &internal.HLSElement{
				Name:  "UspTimestampMap",
				Attrs: params,
			},
		})
		return nil
	}
	return fmt.Errorf("%w: invalid usp timestamp map tag", ErrParseLine)
}

func (p cueOutParser) Parse(tag string, playlist *Playlist) error {
	parts := strings.SplitN(tag, ":", 2)
	if len(parts) > 1 {
		playlist.Insert(&internal.Node{
			HLSElement: &internal.HLSElement{
				Name:  "CueOut",
				Attrs: map[string]string{cueOutTag: strings.TrimSpace(parts[1])},
			},
		})
		return nil
	}
	return fmt.Errorf("%w: invalid cue out tag", ErrParseLine)
}
func (p cueInParser) Parse(tag string, playlist *Playlist) error {
	playlist.Insert(&internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "CueIn",
			Attrs: map[string]string{
				cueInTag: "",
			},
		},
	})

	// When break ends, should reset current daterange
	playlist.CurrentDateRange = &internal.Node{}

	return nil
}

func HandleNonTags(line string, playlist *Playlist) error {
	switch {
	// Handle HLS segment lines
	case playlist.CurrentSegment != nil && strings.HasSuffix(line, ".ts"):
		playlist.CurrentSegment.URI = line
		playlist.Insert(&internal.Node{
			HLSElement: &internal.HLSElement{
				Name: "ExtInf",
				URI:  line,
				Attrs: map[string]string{
					"Duration": strconv.FormatFloat(playlist.CurrentSegment.Duration, 'f', -1, 64),
				},
				Details: map[string]string{
					"MediaSequence":   fmt.Sprintf("%d", playlist.CurrentSegment.MediaSequence),
					"ProgramDateTime": playlist.CurrentSegment.ProgramDateTime.Format(time.RFC3339Nano),
				},
			},
		})
		playlist.CurrentSegment = nil
		return nil

	// Handle HLS media playlist lines
	case playlist.CurrentStreamInf != nil && strings.HasSuffix(line, ".m3u8"):
		playlist.CurrentStreamInf.URI = line
		attrs := map[string]string{
			"BANDWIDTH":         playlist.CurrentStreamInf.Bandwidth,
			"AVERAGE-BANDWIDTH": playlist.CurrentStreamInf.AverageBandwidth,
			"CODECS":            strings.Join(playlist.CurrentStreamInf.Codecs, ","),
			"RESOLUTION":        playlist.CurrentStreamInf.Resolution,
			"FRAME-RATE":        playlist.CurrentStreamInf.FrameRate,
		}
		playlist.Insert(&internal.Node{
			HLSElement: &internal.HLSElement{
				Name:  "StreamInf",
				Attrs: attrs,
				URI:   line,
			},
		})
		playlist.CurrentStreamInf = nil
		return nil
	// Handle Comments
	default:
		attrs := map[string]string{
			"Comment": line,
		}
		playlist.Insert(&internal.Node{
			HLSElement: &internal.HLSElement{
				Name:  "Comment",
				Attrs: attrs,
			},
		})
		return nil
	}
}

// https://regex101.com/r/0A2ulC/1
func tagsToMap(line string) map[string]string {
	m := make(map[string]string)
	for _, kv := range paramRegex.FindAllStringSubmatch(line, -1) {
		k, v := kv[1], kv[2]
		m[strings.ToUpper(k)] = strings.Trim(v, "\"")
	}

	return m
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
