package go_m3u8

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type StreamInf struct {
	Codecs           []string
	Bandwidth        string
	AverageBandwidth string
	Resolution       string
	FrameRate        string
	URI              string
}

type MasterPlaylist struct {
	Variants []StreamInf
	Version  string
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

type MediaPlaylist struct {
	DateRanges            []DateRange
	ProgramsDateTime      []ProgramDateTime
	Segments              []Segment
	Version               string
	MediaSequence         string
	DiscontinuitySequence string
	TargetDuration        string
}

type Parser interface {
	Parse(reader io.Reader) error
}

func ParseMediaPlaylist(filename string) (*MediaPlaylist, error) {
	m := new(MediaPlaylist)
	playlist, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening playlist: %v", err)
	}

	scanner := bufio.NewScanner(playlist)
	var lines []string
	uriIndex := 0
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "#EXT-X-VERSION:"):
			version := m.parseVersion(line)
			m.Version = version

		case strings.HasPrefix(line, "#EXT-X-TARGETDURATION:"):
			targetDuration := m.parseTargetDuration(line)
			m.TargetDuration = targetDuration

		case strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE:"):
			mediaSequence := m.parseMediaSequence(line)
			m.MediaSequence = mediaSequence

		case strings.HasPrefix(line, "#EXT-X-PROGRAM-DATE-TIME"):
			lines = append(lines, scanner.Text())
			for range lines {
				programDateTime := m.parseProgramDateTime(line)
				m.ProgramsDateTime = append(m.ProgramsDateTime, *programDateTime)
				lines = []string{}
			}

		case strings.HasPrefix(line, "#EXT-X-DATERANGE:"):
			lines = append(lines, scanner.Text())
			for range lines {
				dateRange := m.parseDateRange(line)
				m.DateRanges = append(m.DateRanges, *dateRange)
				lines = []string{}
			}

		case strings.HasPrefix(line, "#EXTINF:"):
			lines = append(lines, scanner.Text())
			for range lines {
				duration := m.parseSegmentDuration(line)
				m.Segments = append(m.Segments, *duration)
				lines = []string{}
			}

		case strings.Contains(line, ".ts"):
			m.Segments[uriIndex].URI = line
			uriIndex++
		}

	}
	return m, nil
}

func ParseMasterPlaylist(filename string) (*MasterPlaylist, error) {
	m := new(MasterPlaylist)
	variant := new(StreamInf)
	playlist, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening playlist: %v", err)
	}

	scanner := bufio.NewScanner(playlist)
	var lines []string
	uriIndex := 0

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "#EXT-X-VERSION:"):
			version := m.parseVersion(line)
			m.Version = version

		case strings.HasPrefix(line, "#EXT-X-STREAM-INF:"):
			lines = append(lines, scanner.Text())
			for range lines {
				variant = m.parseStreamInf(line)
				m.Variants = append(m.Variants, *variant)
				lines = []string{}
			}
		case strings.Contains(line, ".m3u8"):
			m.Variants[uriIndex].URI = line
			uriIndex++
		}
	}

	return m, nil
}

func (m *MediaPlaylist) parseSegmentDuration(line string) *Segment {
	segment := new(Segment)
	duration := strings.Split(line, ":")[1]
	duration = strings.Split(duration, ",")[0]
	segment.Duration = duration
	return segment
}

func (m *MediaPlaylist) parseProgramDateTime(line string) *ProgramDateTime {
	programDateTime := new(ProgramDateTime)
	time := strings.Split(line, ":")[1:4]
	dateTime := strings.Join(time, ":")
	programDateTime.DateTime = dateTime
	return programDateTime
}

func (m *MediaPlaylist) parseDateRange(line string) *DateRange {
	params := splitParams(line)
	dateRange := new(DateRange)
	for param, value := range params {
		switch {
		case param == "ID":
			dateRange.Id = value
		case param == "START-DATE":
			dateRange.StartDate = value
		case param == "END-DATE":
			dateRange.EndDate = value
		case param == "PLANNED-DURATION":
			dateRange.PlannedDuration = value
		case param == "DURATION":
			dateRange.Duration = value
		case strings.HasPrefix(param, "SCTE"):
			dateRange.Scte35Mark = m.parseSCTE(param, value)
		}
	}
	return dateRange
}

func (m *MediaPlaylist) parseSCTE(param, value string) map[string]string {
	scte := map[string]string{}
	if strings.HasSuffix(param, "IN") {
		scte = map[string]string{
			"IN": value,
		}
	} else if strings.HasSuffix(param, "OUT") {
		scte = map[string]string{
			"OUT": value,
		}
	}
	return scte
}

func (m *MediaPlaylist) parseMediaSequence(line string) string {
	mediaSequence := strings.Split(line, ":")[1]
	m.MediaSequence = mediaSequence
	return mediaSequence
}

func (m *MediaPlaylist) parseTargetDuration(line string) string {
	duration := strings.Split(line, ":")[1]
	m.TargetDuration = duration
	return duration
}

func (m *MediaPlaylist) parseVersion(line string) string {
	version := strings.Split(line, ":")[1]
	m.Version = version
	return version
}

func (m *MasterPlaylist) parseVersion(line string) string {
	version := strings.Split(line, ":")[1]
	m.Version = version
	return version
}

func (m *MasterPlaylist) parseStreamInf(line string) *StreamInf {
	params := splitParams(line)
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
func splitParams(line string) map[string]string {
	re := regexp.MustCompile(`([a-zA-Z\d_-]+)=("[^"]+"|[^",]+)`)
	m := make(map[string]string)
	for _, kv := range re.FindAllStringSubmatch(line, -1) {
		k, v := kv[1], kv[2]
		m[strings.ToUpper(k)] = strings.Trim(v, "\"")
	}
	return m
}
