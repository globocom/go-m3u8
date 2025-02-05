package go_m3u8

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var (
	ErrOpenPlaylist = errors.New("failed to open playlist")
	ErrParseLine    = errors.New("failed to parse line")
)

type Source interface {
	io.ReadCloser
}
type Playlist struct {
	Head, Tail *Node
	Elements   []Node
}

type Node struct {
	Tags       any
	Prev, Next *Node
}

type StreamInf struct {
	Codecs           []string
	Bandwidth        string
	AverageBandwidth string
	Resolution       string
	FrameRate        string
	URI              string
}

type DateRange struct {
	Scte35Mark      map[string]string
	Id              string
	StartDate       string
	EndDate         string
	PlannedDuration string
	Duration        string
}

type Segment struct {
	Duration string
	URI      string
}

type ProgramDateTime struct {
	DateTime string
}

func ParsePlaylist(src Source) (*Playlist, error) {
	var segment *Segment
	var streamInf *StreamInf
	playlist := new(Playlist)

	scanner := bufio.NewScanner(src)
	defer src.Close()
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case strings.HasPrefix(line, "#EXT-X-VERSION:"):
			version, err := playlist.parseVersion(line)
			if err != nil {
				return nil, fmt.Errorf("failed to parse version: %w", err)
			}
			playlist.Insert(&Node{Tags: version})

		case strings.HasPrefix(line, "#EXT-X-TARGETDURATION:"):
			targetDuration, err := playlist.parseTargetDuration(line)
			if err != nil {
				return nil, fmt.Errorf("failed to parse target duration: %w", err)
			}
			playlist.Insert(&Node{Tags: targetDuration})

		case strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE:"):
			mediaSequence, err := playlist.parseMediaSequence(line)
			if err != nil {
				return nil, fmt.Errorf("failed to parse media sequence: %w", err)
			}
			playlist.Insert(&Node{Tags: mediaSequence})

		case strings.HasPrefix(line, "#EXT-X-PROGRAM-DATE-TIME"):
			programDateTime := playlist.parseProgramDateTime(line)
			playlist.Insert(&Node{Tags: programDateTime})

		case strings.HasPrefix(line, "#EXT-X-DATERANGE:"):
			dateRange := playlist.parseDateRange(line)
			playlist.Insert(&Node{Tags: dateRange})

		case strings.HasPrefix(line, "#EXTINF:"):
			segment = playlist.parseSegmentDuration(line)

		case strings.HasSuffix(line, ".ts"):
			if segment != nil {
				segment.URI = line
				playlist.Insert(&Node{Tags: segment})
				segment = nil
			}
		case strings.HasPrefix(line, "#EXT-X-STREAM-INF:"):
			streamInf = playlist.parseStreamInf(line)

		case strings.HasSuffix(line, ".m3u8"):
			if streamInf != nil {
				streamInf.URI = line
				playlist.Insert(&Node{Tags: streamInf})
				streamInf = nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, ErrOpenPlaylist
	}

	return playlist, nil
}

func (p *Playlist) parseSegmentDuration(line string) *Segment {
	segment := new(Segment)
	parts := strings.Split(line, ":")
	if len(parts) > 1 {
		duration := strings.Split(parts[1], ",")[0]
		segment.Duration = strings.TrimSpace(duration)
	}

	return segment
}

func (p *Playlist) parseProgramDateTime(line string) *ProgramDateTime {
	programDateTime := new(ProgramDateTime)
	parts := strings.SplitN(line, ":", 2)
	if len(parts) > 1 {
		programDateTime.DateTime = strings.TrimSpace(parts[1])
	}

	return programDateTime
}

func (p *Playlist) parseDateRange(line string) *DateRange {
	params := tagsToMap(line)
	dateRange := new(DateRange)
	dateRange.Scte35Mark = make(map[string]string)
	for param, value := range params {
		switch strings.ToUpper(param) {
		case "ID":
			dateRange.Id = value
		case "START-DATE":
			dateRange.StartDate = value
		case "END-DATE":
			dateRange.EndDate = value
		case "PLANNED-DURATION":
			dateRange.PlannedDuration = value
		case "DURATION":
			dateRange.Duration = value
		default:
			if strings.HasPrefix(param, "SCTE") {
				scte := p.parseSCTE(param, value)
				for k, v := range scte {
					dateRange.Scte35Mark[k] = v
				}
			}
		}
	}

	return dateRange
}

func (p *Playlist) parseSCTE(param, value string) map[string]string {
	scte := map[string]string{}
	if strings.HasSuffix(param, "IN") {
		scte["IN"] = value
	} else if strings.HasSuffix(param, "OUT") {
		scte["OUT"] = value
	}

	return scte
}

func (p *Playlist) parseMediaSequence(line string) (string, error) {
	return parseLine(line)
}

func (p *Playlist) parseTargetDuration(line string) (string, error) {
	return parseLine(line)
}

func (p *Playlist) parseVersion(line string) (string, error) {
	return parseLine(line)
}

func parseLine(line string) (string, error) {
	parts := strings.Split(line, ":")
	if len(parts) > 1 && parts[1] != "" {
		return strings.TrimSpace(parts[1]), nil
	}

	return "", ErrParseLine
}

func (p *Playlist) parseStreamInf(line string) *StreamInf {
	params := tagsToMap(line)
	variant := new(StreamInf)

	for param, value := range params {
		switch param {
		case "BANDWIDTH":
			variant.Bandwidth = value
		case "AVERAGE-BANDWIDTH":
			variant.AverageBandwidth = value
		case "RESOLUTION":
			variant.Resolution = value
		case "FRAME-RATE":
			variant.FrameRate = value
		case "CODECS":
			variant.Codecs = strings.Split(value, ",")
		}
	}
	return variant
}

// https://regex101.com/r/0A2ulC/1
func tagsToMap(line string) map[string]string {
	re := regexp.MustCompile(`([a-zA-Z\d_-]+)=("[^"]+"|[^",]+)`)
	m := make(map[string]string)
	for _, kv := range re.FindAllStringSubmatch(line, -1) {
		k, v := kv[1], kv[2]
		m[strings.ToUpper(k)] = strings.Trim(v, "\"")
	}

	return m
}

func (p *Playlist) Insert(node *Node) {
	if p.Head == nil {
		p.Head = node
		p.Tail = node
	} else {
		node.Prev = p.Tail
		p.Tail.Next = node
		p.Tail = node
	}

	p.Elements = append(p.Elements, *node)
}

func (p *Playlist) Find(tag any) (*Node, bool) {
	current := p.Head
	for current != nil {
		if current.Tags == tag {
			return current, true
		}
		current = current.Next
	}

	return nil, false
}
