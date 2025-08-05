package playlist

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/globocom/go-m3u8/internal"
	"github.com/rs/zerolog/log"
)

// METHODS FOR DECODING MULTI-LINE TAGS

// StreamInfData holds data for StreamInf HLS Element, whose format in manifest is multi-line:
//
//	#EXT-X-STREAM-INF:<attribute-list>
//	<URI>
type StreamInfData struct {
	Codecs           []string
	Bandwidth        string
	AverageBandwidth string
	Resolution       string
	FrameRate        string
	VideoRange       string
	Audio            string
	Video            string
	Subtitles        string
	ClosedCaptions   string
}

// ExtInfData holds data for ExtInf HLS element, whose format in manifest is multi-line:
//
//	#EXTINF:<duration>,[<title>]
//	<URI>
type ExtInfData struct {
	Duration        float64
	ProgramDateTime time.Time
	MediaSequence   int
	URI             string
	Title           string
}

// Parser function that returns new StreamInfData object.
func GetStreamInfData(mappedAttr map[string]string) *StreamInfData {
	return &StreamInfData{
		Bandwidth:        mappedAttr["BANDWIDTH"],
		AverageBandwidth: mappedAttr["AVERAGE-BANDWIDTH"],
		Codecs:           strings.Split(mappedAttr["CODECS"], ","),
		Resolution:       mappedAttr["RESOLUTION"],
		FrameRate:        mappedAttr["FRAME-RATE"],
		VideoRange:       mappedAttr["VIDEO-RANGE"],
		Audio:            mappedAttr["AUDIO"],
		Video:            mappedAttr["VIDEO"],
		Subtitles:        mappedAttr["SUBTITLES"],
		ClosedCaptions:   mappedAttr["CLOSED-CAPTIONS"],
	}
}

// Parser function that returns new ExtInfData object.
func GetExtInfData(duration, title string, playlistMediaSequence, playlistSegmentsCounter int, playlistDVR float64, playlistPDT time.Time) *ExtInfData {
	floatDuration, err := strconv.ParseFloat(duration, 64)
	if err != nil {
		log.Error().Str("service", "go-m3u8/playlist/helpers.go").Err(err).Msgf("failed to parse duration for segment: %s", duration)
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
	// handle EXTINF
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

	// handle EXT-X-STREAM-INF
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
					"VIDEO-RANGE":       p.CurrentStreamInf.VideoRange,
					"AUDIO":             p.CurrentStreamInf.Audio,
					"VIDEO":             p.CurrentStreamInf.Video,
					"SUBTITLES":         p.CurrentStreamInf.Subtitles,
					"CLOSED-CAPTIONS":   p.CurrentStreamInf.ClosedCaptions,
				},
			},
		})
		p.CurrentStreamInf = nil
		return nil
	default:
		return nil
	}
}

// AUXILIARY METHODS FOR DECODING

// Converts given tag's (line) attributes into a map of key-value pairs.
// https://regex101.com/r/0A2ulC/1
func TagsToMap(line string) map[string]string {
	m := make(map[string]string)
	for _, kv := range ParamRegex.FindAllStringSubmatch(line, -1) {
		k, v := kv[1], kv[2]
		m[strings.ToUpper(k)] = strings.Trim(v, "\"")
	}

	return m
}

// Rounds a float64 value to a specified precision.
func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// AUXILIARY METHODS FOR ENCODING

// Encodes a tag with key-value attributes into a string.
func EncodeTagWithAttributes(builder *strings.Builder, tag string, attrs map[string]string, order []string, shouldQuote map[string]bool) error {
	if len(attrs) == 0 {
		_, err := builder.WriteString(tag + "\n")
		return err
	}

	formattedAttrs := make([]string, 0, len(attrs))
	processed := make(map[string]bool)

	for _, key := range order {
		if value, exists := attrs[key]; exists && value != "" {
			formattedAttrs = append(formattedAttrs, FormatAttribute(key, value, shouldQuote))
			processed[key] = true
		}
	}

	unorderedKeys := make([]string, 0, len(attrs))
	for key := range attrs {
		if !processed[key] && attrs[key] != "" {
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

// Encodes a tag without attributes into a string.
func EncodeSimpleTag(node *internal.Node, builder *strings.Builder, tag, attrKey string) error {
	if value, exists := node.HLSElement.Attrs[attrKey]; exists {
		attr := fmt.Sprintf("%s:%s\n", tag, value)
		_, err := builder.WriteString(attr)
		return err
	}
	return fmt.Errorf("attribute %s not found for tag %s", attrKey, tag)
}

// Formats a key-value tag attribute into a string, optionally quoting the value based on the shouldQuote map.
func FormatAttribute(key, value string, shouldQuote map[string]bool) string {
	shouldQuoteValue, exists := shouldQuote[key]

	if !exists {
		shouldQuoteValue = true // default to quoting if not specified
	}

	if shouldQuoteValue {
		return fmt.Sprintf(`%s=%q`, key, value)
	}

	return fmt.Sprintf(`%s=%s`, key, value)
}
