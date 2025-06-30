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
	CurrentSegment        *ExtInfData
	CurrentStreamInf      *StreamInfData
	ProgramDateTime       time.Time
	MediaSequence         int
	DiscontinuitySequence int
	SegmentsCounter       int
	DVR                   float64
}

func NewPlaylist() *Playlist {
	return &Playlist{
		DoublyLinkedList:      new(internal.DoublyLinkedList),
		CurrentSegment:        nil,
		CurrentStreamInf:      nil,
		ProgramDateTime:       *new(time.Time),
		MediaSequence:         0,
		DiscontinuitySequence: 0,
		SegmentsCounter:       0,
		DVR:                   0,
	}
}

func (p *Playlist) Print() {
	if p.Head == nil || p.Tail == nil {
		log.Warn().Msg("playlist is empty")
		return
	}

	current := p.Head
	for current != nil {
		fmt.Printf(">>>>>>>>> Node: %+v\n", current)
		if current.HLSElement != nil {
			fmt.Printf("HLSElement: %+v\n", current.HLSElement)
		}
		current = current.Next
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

func (p *Playlist) DiscontinuitySequenceTag() (*internal.Node, bool) {
	return p.Find("DiscontinuitySequence")
}

func (p *Playlist) DiscontinuitySequenceValue() string {
	node, found := p.Find("DiscontinuitySequence")
	if !found {
		return ""
	}
	return node.HLSElement.Attrs["#EXT-X-DISCONTINUITY-SEQUENCE"]
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

// Returns true if media segment is inside ad break and false otherwise.
// When true, method also returns de DateRange object for the segment Ad Break.
//
// For entering the Ad Break, we always have DateRange tag with SCTE-OUT and CueOutEvent tag.
// However, for exiting the Ad Break, we have three possible manifests:
//
//   - DateRange SCTE-IN is ALWAYS present.
//   - No DateRange SCTE-IN. Exit is ONLY marked by CueInEvent tag instead.
//   - SOMETIMES DateRange SCTE-IN is present, alongside the CueInEvent tag.
func (p *Playlist) FindSegmentAdBreak(segment *internal.Node) (*internal.Node, bool) {
	current := segment
	for current != nil {
		// segment is inside Ad Break if it is preceeded by a DateRange tag with attribute SCTE35-OUT
		if (current.HLSElement.Name == "DateRange") && (current.HLSElement.Attrs["SCTE35-OUT"] != "") {
			return current, true
		}

		// segment is outside Ad Break if it is preceeded by a CueIn tag or a DateRange tag with attribute SCTE35-IN
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

//// METHODS FOR DECODING MULTI-LINE TAGS

// StreamInfData holds data for StreamInf HLS Element, whose format in manifest is multi-line:
//
// #EXT-X-STREAM-INF:<attribute-list>
//
// <URI>
type StreamInfData struct {
	Codecs           []string
	Bandwidth        string
	AverageBandwidth string
	Resolution       string
	FrameRate        string
	URI              string
}

// ExtInfData holds data for ExtInf HLS element, whose format in manifest is multi-line:
//
// #EXTINF:<duration>,[<title>]
//
// <URI>
type ExtInfData struct {
	Duration        float64
	ProgramDateTime time.Time
	MediaSequence   int
	URI             string
	Title           string
}

// Internal parser returns new StreamInfData object
func GetStreamInfData(mappedAttr map[string]string) *StreamInfData {
	return &StreamInfData{
		Bandwidth:        mappedAttr["BANDWIDTH"],
		AverageBandwidth: mappedAttr["AVERAGE-BANDWIDTH"],
		Codecs:           strings.Split(mappedAttr["CODECS"], ","),
		Resolution:       mappedAttr["RESOLUTION"],
		FrameRate:        mappedAttr["FRAME-RATE"],
	}
}

// Internal parser returns new ExtInfData object
func GetExtInfData(duration, title string, playlistMediaSequence, playlistSegmentsCounter int, playlistDVR float64, playlistPDT time.Time) *ExtInfData {
	floatDuration, err := strconv.ParseFloat(duration, 64)
	if err != nil {
		log.Error().Err(err).Msgf("failed to parse duration for segment: %s", duration)
		return &ExtInfData{}
	}

	currentDVRInNanoseconds := int(playlistDVR * float64(time.Second))
	segmentProgramDateTime := playlistPDT.Add(time.Duration(currentDVRInNanoseconds))

	return &ExtInfData{
		Duration:        floatDuration,
		Title:           title,
		MediaSequence:   playlistMediaSequence + playlistSegmentsCounter,
		ProgramDateTime: segmentProgramDateTime,
	}
}

// Handles HLS Elements whose format in manifest are multi-line: tag + uri.
// The URI line that follows the EXT-X-STREAM-INF and EXTINF tags is REQUIRED.
func HandleMultiLineHLSElements(line string, p *Playlist) error {
	switch {
	// Handle EXTINF
	case p.CurrentSegment != nil:
		p.Insert(&internal.Node{
			HLSElement: &internal.HLSElement{
				Name: "ExtInf",
				URI:  line,
				Attrs: map[string]string{
					"Duration": strconv.FormatFloat(p.CurrentSegment.Duration, 'f', -1, 64),
					"Title":    p.CurrentSegment.Title,
				},
				Details: map[string]string{
					"MediaSequence":   fmt.Sprintf("%d", p.CurrentSegment.MediaSequence),
					"ProgramDateTime": p.CurrentSegment.ProgramDateTime.Format(time.RFC3339Nano),
				},
			},
		})
		p.CurrentSegment = nil
		return nil

	// Handle EXT-X-STREAM-INF
	case p.CurrentStreamInf != nil:
		p.Insert(&internal.Node{
			HLSElement: &internal.HLSElement{
				Name: "StreamInf",
				URI:  line,
				Attrs: map[string]string{
					"BANDWIDTH":         p.CurrentStreamInf.Bandwidth,
					"AVERAGE-BANDWIDTH": p.CurrentStreamInf.AverageBandwidth,
					"CODECS":            strings.Join(p.CurrentStreamInf.Codecs, ","),
					"RESOLUTION":        p.CurrentStreamInf.Resolution,
					"FRAME-RATE":        p.CurrentStreamInf.FrameRate,
				},
			},
		})
		p.CurrentStreamInf = nil
		return nil
	default:
		return nil
	}
}

//// AUXILIARY METHODS FOR DECODING

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

//// AUXILIARY METHODS FOR ENCODING

// Encodes tag with attributes into string object
func EncodeTagWithAttributes(builder *strings.Builder, tag string, attrs map[string]string, order []string, shouldQuote map[string]bool) error {
	if len(attrs) == 0 {
		_, err := builder.WriteString(tag + "\n")
		return err
	}

	var formattedAttrs []string
	processed := make(map[string]bool)

	for _, key := range order {
		if value, exists := attrs[key]; exists {
			formattedAttrs = append(formattedAttrs, FormatAttribute(key, value, shouldQuote))
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
		formattedAttrs = append(formattedAttrs, FormatAttribute(key, attrs[key], shouldQuote))
	}

	attributes := fmt.Sprintf("%s:%s\n", tag, strings.Join(formattedAttrs, ","))

	_, err := builder.WriteString(attributes)
	return err
}

// Encodes tag without attributes into string object
func EncodeSimpleTag(node *internal.Node, builder *strings.Builder, tag, attrKey string) error {
	if value, exists := node.HLSElement.Attrs[attrKey]; exists {
		attr := fmt.Sprintf("%s:%s\n", tag, value)
		_, err := builder.WriteString(attr)
		return err
	}
	return fmt.Errorf("attribute %s not found for tag %s", attrKey, tag)
}

func FormatAttribute(key, value string, shouldQuote map[string]bool) string {
	shouldQuoteValue, exists := shouldQuote[key]

	if !exists {
		shouldQuoteValue = true // default to quoting if not specified
	}

	if shouldQuoteValue {
		return fmt.Sprintf(`%s="%s"`, key, value)
	}

	return fmt.Sprintf(`%s=%s`, key, value)
}
