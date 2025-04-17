package playlist

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/globocom/go-m3u8/internal"
	"github.com/rs/zerolog/log"
)

var (
	ErrParseLine = errors.New("failed to parse tag")
	ParamRegex   = regexp.MustCompile(`([a-zA-Z\d_-]+)=("[^"]+"|[^",]+)`)
)

type Playlist struct {
	*internal.DoublyLinkedList
	CurrentNode      *internal.Node
	CurrentSegment   *ExtInfData
	CurrentStreamInf *StreamInfData
	ProgramDateTime  time.Time
	MediaSequence    int
	SegmentsCounter  int
	DVR              float64
}

func NewPlaylist() *Playlist {
	return &Playlist{
		DoublyLinkedList: new(internal.DoublyLinkedList),
		CurrentNode:      new(internal.Node),
		CurrentSegment:   new(ExtInfData),
		CurrentStreamInf: new(StreamInfData),
		ProgramDateTime:  *new(time.Time),
		MediaSequence:    0,
		SegmentsCounter:  0,
		DVR:              0,
	}
}

func (p *Playlist) VersionValue() string {
	node, found := p.Find("Version")
	if !found {
		return ""
	}
	return node.HLSElement.Attrs["#EXT-X-VERSION"]
}

func (p *Playlist) Version() (*internal.Node, bool) {
	return p.Find("Version")
}

func (p *Playlist) MediaSequenceValue() string {
	node, found := p.Find("MediaSequence")
	if !found {
		return ""
	}
	return node.HLSElement.Attrs["#EXT-X-MEDIA-SEQUENCE"]
}

func (p *Playlist) MediaSequenceTag() (*internal.Node, bool) {
	return p.Find("MediaSequence")
}

func (p *Playlist) Variants() []*internal.Node {
	return p.FindAll("StreamInf")
}

func (p *Playlist) Segments() []*internal.Node {
	return p.FindAll("ExtInf")
}

func (p *Playlist) Breaks() []*internal.Node {
	result := make([]*internal.Node, 0)
	nodes := p.FindAll("DateRange")
	for _, node := range nodes {
		if node.HLSElement.Attrs["SCTE35-OUT"] != "" {
			result = append(result, node)
		}
	}
	return result
}

func (p *Playlist) FindSegment(segmentName, segmentURI string) (*internal.Node, bool) {
	current := p.Head
	for current != nil {
		if (current.HLSElement.Name == segmentName) && (current.HLSElement.URI == segmentURI) {
			return current, true
		}
		current = current.Next
	}

	return nil, false
}

func (p *Playlist) FindSegmentAdBreak(segmentName, segmentURI string) (*internal.Node, bool) {
	segment, ok := p.FindSegment(segmentName, segmentURI)
	if !ok {
		return nil, false
	}

	current := segment
	for current != nil {
		if (current.HLSElement.Name == "DateRange") && (current.HLSElement.Attrs["SCTE35-OUT"] != "") {
			return current, true
		}

		if (current.HLSElement.Name == "CueIn") || (current.HLSElement.Name == "DateRange" && current.HLSElement.Attrs["SCTE35-IN"] != "") {
			return nil, false
		}

		current = current.Prev
	}

	return nil, false
}

func (p *Playlist) ReplaceBreaksURI(transform func(string) string) error {
	startCondition := func(node *internal.Node) bool {
		return node.HLSElement.Name == "DateRange" && node.HLSElement.Attrs["SCTE35-OUT"] != ""
	}
	endCondition := func(node *internal.Node) bool {
		return node.HLSElement.Name == "DateRange" && node.HLSElement.Attrs["SCTE35-IN"] != ""
	}
	transformFunc := func(node *internal.Node) {
		if node.HLSElement.Name == "ExtInf" && node.HLSElement.URI != "" {
			node.HLSElement.URI = transform(node.HLSElement.URI)
		}
	}
	return p.ModifyNodesBetween(startCondition, endCondition, transformFunc)
}

// / METHODS FOR MULTI-LINE TAGS (StreamInf and ExtInf)
type StreamInfData struct {
	Codecs           []string
	Bandwidth        string
	AverageBandwidth string
	Resolution       string
	FrameRate        string
	URI              string
}

type ExtInfData struct {
	Duration        float64
	ProgramDateTime time.Time
	MediaSequence   int
	URI             string
}

func GetStreamInfData(mappedAttr map[string]string) *StreamInfData {
	return &StreamInfData{
		Bandwidth:        mappedAttr["BANDWIDTH"],
		AverageBandwidth: mappedAttr["AVERAGE-BANDWIDTH"],
		Codecs:           strings.Split(mappedAttr["CODECS"], ","),
		Resolution:       mappedAttr["RESOLUTION"],
		FrameRate:        mappedAttr["FRAME-RATE"],
	}
}

func GetExtInfData(duration string, playlistMediaSequence, playlistSegmentsCounter int, playlistDVR float64, playlistPDT time.Time) *ExtInfData {
	floatDuration, err := strconv.ParseFloat(duration, 64)
	if err != nil {
		log.Error().Err(err).Msgf("failed to parse duration for segment: %s", duration)
		return &ExtInfData{}
	}

	currentDVRInNanoseconds := int(playlistDVR * float64(time.Second))
	segmentProgramDateTime := playlistPDT.Add(time.Duration(currentDVRInNanoseconds))

	return &ExtInfData{
		Duration:        floatDuration,
		MediaSequence:   playlistMediaSequence + playlistSegmentsCounter,
		ProgramDateTime: segmentProgramDateTime,
	}
}

// DECODE METHODS
func HandleNonTags(line string, p *Playlist) error {
	switch {
	// Handle HLS segment lines
	case p.CurrentSegment != nil && strings.HasSuffix(line, ".ts"):
		p.CurrentSegment.URI = line
		p.Insert(&internal.Node{
			HLSElement: &internal.HLSElement{
				Name: "ExtInf",
				URI:  line,
				Attrs: map[string]string{
					"Duration": strconv.FormatFloat(p.CurrentSegment.Duration, 'f', -1, 64),
				},
				Details: map[string]string{
					"MediaSequence":   fmt.Sprintf("%d", p.CurrentSegment.MediaSequence),
					"ProgramDateTime": p.CurrentSegment.ProgramDateTime.Format(time.RFC3339Nano),
				},
			},
		})
		p.CurrentSegment = nil
		return nil

	// Handle HLS media playlist lines
	case p.CurrentStreamInf != nil && strings.HasSuffix(line, ".m3u8"):
		p.CurrentStreamInf.URI = line
		attrs := map[string]string{
			"BANDWIDTH":         p.CurrentStreamInf.Bandwidth,
			"AVERAGE-BANDWIDTH": p.CurrentStreamInf.AverageBandwidth,
			"CODECS":            strings.Join(p.CurrentStreamInf.Codecs, ","),
			"RESOLUTION":        p.CurrentStreamInf.Resolution,
			"FRAME-RATE":        p.CurrentStreamInf.FrameRate,
		}
		p.Insert(&internal.Node{
			HLSElement: &internal.HLSElement{
				Name:  "StreamInf",
				Attrs: attrs,
				URI:   line,
			},
		})
		p.CurrentStreamInf = nil
		return nil
	// Handle Comments
	default:
		attrs := map[string]string{
			"Comment": line,
		}
		p.Insert(&internal.Node{
			HLSElement: &internal.HLSElement{
				Name:  "Comment",
				Attrs: attrs,
			},
		})
		return nil
	}
}

// https://regex101.com/r/0A2ulC/1
func TagsToMap(line string) map[string]string {
	m := make(map[string]string)
	for _, kv := range ParamRegex.FindAllStringSubmatch(line, -1) {
		k, v := kv[1], kv[2]
		m[strings.ToUpper(k)] = strings.Trim(v, "\"")
	}

	return m
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// // ENCODE METHODS
func EncodeTagWithAttributes(builder *strings.Builder, tag string, attrs map[string]string, order []string) error {
	if len(attrs) == 0 {
		_, err := builder.WriteString(tag + "\n")
		return err
	}

	var formattedAttrs []string
	processed := make(map[string]bool)

	for _, key := range order {
		if value, exists := attrs[key]; exists {
			formattedAttrs = append(formattedAttrs, FormatAttribute(key, value))
			processed[key] = true
		}
	}

	unorderedKeys := make([]string, 0, len(attrs))
	for key := range attrs {
		if !processed[key] {
			unorderedKeys = append(unorderedKeys, key)
		}
	}
	sort.Strings(unorderedKeys)
	for _, key := range unorderedKeys {
		formattedAttrs = append(formattedAttrs, FormatAttribute(key, attrs[key]))
	}

	attributes := fmt.Sprintf("%s:%s\n", tag, strings.Join(formattedAttrs, ","))

	_, err := builder.WriteString(attributes)
	return err
}

func EncodeSimpleTag(node *internal.Node, builder *strings.Builder, tag, attrKey string) error {
	if value, exists := node.HLSElement.Attrs[attrKey]; exists {
		attr := fmt.Sprintf("%s:%s\n", tag, value)
		_, err := builder.WriteString(attr)
		return err
	}
	return fmt.Errorf("attribute %s not found for tag %s", attrKey, tag)
}

func FormatAttribute(key, value string) string {
	if ShouldQuote(key, value) {
		return fmt.Sprintf(`%s="%s"`, key, value)
	}
	return fmt.Sprintf(`%s=%s`, key, value)
}

func ShouldQuote(key, value string) bool {
	numericAttrs := map[string]bool{
		"BANDWIDTH":         true,
		"AVERAGE-BANDWIDTH": true,
		"FRAME-RATE":        true,
		"RESOLUTION":        true,
		"PLANNED-DURATION":  true,
		"DURATION":          true,
		"MPEGTS":            true,
	}

	hexAttrs := map[string]bool{
		"SCTE35-OUT": true,
		"SCTE35-IN":  true,
	}

	if numericAttrs[key] {
		if _, err := strconv.ParseFloat(value, 64); err == nil {
			return false
		}
	}

	if hexAttrs[key] && strings.HasPrefix(value, "0x") {
		return false
	}

	if key == "LOCAL" {
		return false
	}

	return true
}
