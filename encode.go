package go_m3u8

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/globocom/go-m3u8/internal"
)

type PlaylistEncoder interface {
	Encode(node *internal.Node, builder *strings.Builder) error
}

var encoders = map[string]PlaylistEncoder{
	"M3u8Identifier":      m3u8IdentifierEncoder{},
	"Version":             versionEncoder{},
	"TargetDuration":      targetDurationEncoder{},
	"MediaSequence":       mediaSequenceEncoder{},
	"ProgramDateTime":     programDateTimeEncoder{},
	"DateRange":           dateRangeEncoder{},
	"StreamInf":           streamInfEncoder{},
	"ExtInf":              extInfEncoder{},
	"IndependentSegments": independentSegmentsEncoder{},
	"Discontinuity":       discontinuityEncoder{},
	"UspTimestampMap":     uspTimestampMapEncoder{},
	"CueOut":              cueOutEncoder{},
	"CueIn":               cueInEncoder{},
	"Comment":             commentEncoder{},
}

type (
	m3u8IdentifierEncoder      struct{}
	versionEncoder             struct{}
	targetDurationEncoder      struct{}
	mediaSequenceEncoder       struct{}
	programDateTimeEncoder     struct{}
	extInfEncoder              struct{}
	streamInfEncoder           struct{}
	commentEncoder             struct{}
	dateRangeEncoder           struct{}
	independentSegmentsEncoder struct{}
	discontinuityEncoder       struct{}
	uspTimestampMapEncoder     struct{}
	cueOutEncoder              struct{}
	cueInEncoder               struct{}
)

func EncodePlaylist(playlist *Playlist) (string, error) {
	if playlist == nil || playlist.Head == nil {
		return "", fmt.Errorf("playlist is empty")
	}

	var builder strings.Builder
	current := playlist.Head
	for current != nil {
		encoder, exists := encoders[current.HLSElement.Name]
		if !exists {
			return "", fmt.Errorf("unknown tag: %s", current.HLSElement.Name)
		}
		if err := encoder.Encode(current, &builder); err != nil {
			return "", fmt.Errorf("error encoding tag %s: %w", current.HLSElement.Name, err)
		}
		current = current.Next
	}
	return builder.String(), nil
}

func (e m3u8IdentifierEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	_, err := builder.WriteString(m3u8IdentifierTag + "\n")
	return err
}

func (e versionEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	return encodeSimpleTag(node, builder, versionTag, versionTag)
}

func (e targetDurationEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	return encodeSimpleTag(node, builder, targetDurationTag, targetDurationTag)
}

func (e mediaSequenceEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	return encodeSimpleTag(node, builder, mediaSequenceTag, mediaSequenceTag)
}

func (e programDateTimeEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	return encodeSimpleTag(node, builder, programDateTimeTag, programDateTimeTag)
}
func (e extInfEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	_, err := builder.WriteString(fmt.Sprintf("%s:%s\n%s\n", extInfTag, node.HLSElement.Attrs["Duration"], node.HLSElement.URI))
	return err
}

func (e streamInfEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	order := []string{"BANDWIDTH", "AVERAGE-BANDWIDTH", "CODECS", "RESOLUTION", "FRAME-RATE"}
	if err := encodeTagWithAttributes(builder, streamInfTag, node.HLSElement.Attrs, order); err != nil {
		return err
	}
	if node.HLSElement.URI != "" {
		_, err := builder.WriteString(node.HLSElement.URI + "\n")
		return err
	}
	return nil
}

func (e commentEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	_, err := builder.WriteString(fmt.Sprintf("%s\n", node.HLSElement.Attrs["Comment"]))
	return err
}

func (e dateRangeEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	order := []string{"ID", "START-DATE", "PLANNED-DURATION", "END-DATE", "DURATION", "SCTE35-OUT", "SCTE35-IN"}
	return encodeTagWithAttributes(builder, dateRangeTag, node.HLSElement.Attrs, order)
}

func (e independentSegmentsEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	_, err := builder.WriteString(independentSegmentTag + "\n")
	return err
}

func (e discontinuityEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	_, err := builder.WriteString(discontinuityTag + "\n")
	return err
}

func (e uspTimestampMapEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	order := []string{"MPEGTS", "LOCAL"}
	return encodeTagWithAttributes(builder, uspTimestampMapTag, node.HLSElement.Attrs, order)
}

func (e cueOutEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	return encodeSimpleTag(node, builder, cueOutTag, cueOutTag)
}

func (e cueInEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	_, err := builder.WriteString(cueInTag + "\n")
	return err
}

func encodeTagWithAttributes(builder *strings.Builder, tag string, attrs map[string]string, order []string) error {
	if len(attrs) == 0 {
		_, err := builder.WriteString(tag + "\n")
		return err
	}

	var formattedAttrs []string
	processed := make(map[string]bool)

	for _, key := range order {
		if value, exists := attrs[key]; exists {
			formattedAttrs = append(formattedAttrs, formatAttribute(key, value))
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
		formattedAttrs = append(formattedAttrs, formatAttribute(key, attrs[key]))
	}

	_, err := builder.WriteString(fmt.Sprintf("%s:%s\n", tag, strings.Join(formattedAttrs, ",")))
	return err
}

func encodeSimpleTag(node *internal.Node, builder *strings.Builder, tag, attrKey string) error {
	if value, exists := node.HLSElement.Attrs[attrKey]; exists {
		_, err := builder.WriteString(fmt.Sprintf("%s:%s\n", tag, value))
		return err
	}
	return fmt.Errorf("attribute %s not found for tag %s", attrKey, tag)
}

func formatAttribute(key, value string) string {
	if shouldQuote(key, value) {
		return fmt.Sprintf(`%s="%s"`, key, value)
	}
	return fmt.Sprintf(`%s=%s`, key, value)
}

func shouldQuote(key, value string) bool {
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
