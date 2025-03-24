package go_m3u8

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
	"unicode"

	"github.com/globocom/go-m3u8/internal"
)

type Playlist struct {
	*internal.DoublyLinkedList
	CurrentSegment   *Segment
	CurrentStreamInf *StreamInf
	ProgramDateTime  time.Time
	MediaSequence    int
	SegmentsCounter  int
	DVR              float64
}

type TagParser interface {
	Parse(tag string, playlist *Playlist) error
}

var parsers = map[string]TagParser{
	m3u8IdentifierTag:     m3u8IdentifierParser{},
	versionTag:            versionParser{},
	targetDurationTag:     targetDurationParser{},
	mediaSequenceTag:      mediaSequenceParser{},
	programDateTimeTag:    programDateTimeParser{},
	dateRangeTag:          dateRangeParser{},
	extInfTag:             extInfParser{},
	streamInfTag:          streamInfParser{},
	independentSegmentTag: independentSegmentsParser{},
	discontinuityTag:      discontinuityParser{},
	uspTimestampMapTag:    uspTimestampMapParser{},
	cueOutTag:             cueOutParser{},
	cueInTag:              cueInParser{},
}

type Source interface {
	io.ReadCloser
}

func ParsePlaylist(src Source) (*Playlist, error) {
	playlist := &Playlist{
		DoublyLinkedList: new(internal.DoublyLinkedList),
		CurrentSegment:   new(Segment),
		CurrentStreamInf: new(StreamInf),
		ProgramDateTime:  *new(time.Time),
		MediaSequence:    0,
		SegmentsCounter:  0,
		DVR:              0,
	}

	scanner := bufio.NewScanner(src)
	defer src.Close()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		linePrefix := extractPrefix(line)
		parser, exists := parsers[linePrefix]
		if exists {
			if err := parser.Parse(line, playlist); err != nil {
				return nil, fmt.Errorf("error parsing tag %s: %w", linePrefix, err)
			}
		} else {
			if err := HandleNonTags(line, playlist); err != nil {
				return nil, fmt.Errorf("error handling non-tag line %q: %w", line, err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse playlist at line: %q, error: %w", scanner.Text(), err)
	}

	return playlist, nil
}

func extractPrefix(line string) string {
	if line == "" {
		return ""
	}

	for i, r := range line {
		if r == ':' || unicode.IsSpace(r) {
			return line[:i]
		}
	}
	return line
}
