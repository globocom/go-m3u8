package go_m3u8

import (
	"errors"
	"fmt"
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
		Name: "M3u8Identifier",
		Attrs: map[string]string{
			m3u8IdentifierTag: "",
		},
	})
	return nil
}

func (p versionParser) Parse(tag string, playlist *Playlist) error {
	parts := strings.Split(tag, ":")
	if len(parts) > 1 && parts[1] != "" {
		playlist.Insert(&internal.Node{
			Name:  "Version",
			Attrs: map[string]string{versionTag: strings.TrimSpace(parts[1])},
		})
		return nil
	}
	return fmt.Errorf("%w: invalid version tag", ErrParseLine)
}

func (p targetDurationParser) Parse(tag string, playlist *Playlist) error {
	parts := strings.Split(tag, ":")
	if len(parts) > 1 && parts[1] != "" {
		playlist.Insert(&internal.Node{
			Name:  "TargetDuration",
			Attrs: map[string]string{targetDurationTag: strings.TrimSpace(parts[1])},
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
			Name:  "MediaSequence",
			Attrs: map[string]string{mediaSequenceTag: mediaSequence},
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
		Name:  "ProgramDateTime",
		Attrs: map[string]string{programDateTimeTag: dateTimeValue},
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
	if len(params) > 0 {
		playlist.Insert(&internal.Node{
			Name:  "DateRange",
			Attrs: params,
		})
		return nil
	}
	return fmt.Errorf("%w: invalid date range tag", ErrParseLine)
}

func (p streamInfParser) Parse(tag string, playlist *Playlist) error {
	params := tagsToMap(tag)
	playlist.CurrentStreamInf = &StreamInf{
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
		floatDuration, _ := strconv.ParseFloat(duration, 64)
		dvrInNanoseconds := playlist.DVR * float64(time.Second)

		playlist.CurrentSegment = &Segment{
			Duration:        floatDuration,
			MediaSequence:   playlist.MediaSequence + playlist.SegmentsCounter,
			ProgramDateTime: playlist.ProgramDateTime.Add(time.Duration(dvrInNanoseconds)),
		}

		playlist.DVR += floatDuration
		playlist.SegmentsCounter += 1

		return nil
	}
	return fmt.Errorf("%w: invalid extension tag", ErrParseLine)
}

func (p independentSegmentsParser) Parse(tag string, playlist *Playlist) error {
	playlist.Insert(&internal.Node{
		Name: "IndependentSegments",
		Attrs: map[string]string{
			independentSegmentTag: "",
		},
	})
	return nil
}

func (p discontinuityParser) Parse(tag string, playlist *Playlist) error {
	playlist.Insert(&internal.Node{
		Name: "Discontinuity",
		Attrs: map[string]string{
			discontinuityTag: "",
		},
	})
	return nil
}

func (p uspTimestampMapParser) Parse(tag string, playlist *Playlist) error {
	parts := strings.SplitN(tag, ":", 2)
	if len(parts) > 0 {
		params := tagsToMap(parts[1])
		playlist.Insert(&internal.Node{
			Name:  "UspTimestampMap",
			Attrs: params,
		})
		return nil
	}
	return fmt.Errorf("%w: invalid usp timestamp map tag", ErrParseLine)
}

func (p cueOutParser) Parse(tag string, playlist *Playlist) error {
	parts := strings.SplitN(tag, ":", 2)
	if len(parts) > 1 {
		playlist.Insert(&internal.Node{
			Name:  "CueOut",
			Attrs: map[string]string{cueOutTag: strings.TrimSpace(parts[1])},
		})
		return nil
	}
	return fmt.Errorf("%w: invalid cue out tag", ErrParseLine)
}
func (p cueInParser) Parse(tag string, playlist *Playlist) error {
	playlist.Insert(&internal.Node{
		Name: "CueIn",
		Attrs: map[string]string{
			cueInTag: "",
		},
	})
	return nil
}

func HandleNonTags(line string, playlist *Playlist) error {
	switch {
	// Handle HLS segment lines
	case playlist.CurrentSegment != nil && strings.HasSuffix(line, ".ts"):
		playlist.CurrentSegment.URI = line
		playlist.Insert(&internal.Node{
			Name:   "ExtInf",
			Attrs:  map[string]string{"Duration": strconv.FormatFloat(playlist.CurrentSegment.Duration, 'f', -1, 64)},
			URI:    line,
			Object: playlist.CurrentSegment,
		},
		)
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
			Name:   "StreamInf",
			Attrs:  attrs,
			URI:    line,
			Object: playlist.CurrentStreamInf,
		})
		playlist.CurrentStreamInf = nil
		return nil
	// Comment tags
	default:
		attrs := map[string]string{
			"Comment": line,
		}
		playlist.Insert(&internal.Node{
			Name:  "Comment",
			Attrs: attrs,
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
