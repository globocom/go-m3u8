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
)

var (
	ErrParseLine = errors.New("failed to parse tag")
	ParamRegex   = regexp.MustCompile(`([a-zA-Z\d_-]+)=("[^"]+"|[^",]+)`)
)

type Playlist struct {
	*internal.DoublyLinkedList
	CurrentNode      *internal.Node
	CurrentDateRange *internal.DateRange
	CurrentSegment   *internal.Segment
	CurrentStreamInf *internal.StreamInf
	ProgramDateTime  time.Time
	MediaSequence    int
	SegmentsCounter  int
	DVR              float64
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

//// PARSER/DECODE METHODS

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

	_, err := builder.WriteString(fmt.Sprintf("%s:%s\n", tag, strings.Join(formattedAttrs, ",")))
	return err
}

func EncodeSimpleTag(node *internal.Node, builder *strings.Builder, tag, attrKey string) error {
	if value, exists := node.HLSElement.Attrs[attrKey]; exists {
		_, err := builder.WriteString(fmt.Sprintf("%s:%s\n", tag, value))
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
