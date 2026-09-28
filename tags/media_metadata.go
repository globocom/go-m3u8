//	Media Metadata Tags (Section 4.4.5 on RFC)
//
// Media Metadata tags provide information about the playlist that is
// not associated with specific Media Segments.  There MAY be more than
// one Media Metadata tag of each type in any Media Playlist.  The only
// exception to this rule is EXT-X-SKIP, which MUST NOT appear more than
// once.
package tags

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/globocom/go-m3u8/internal"
	pl "github.com/globocom/go-m3u8/playlist"
	"github.com/rs/zerolog/log"
)

const (
	BreakStatusLeavingDVR = "leavingDVRLimit"
	BreakStatusNotReady   = "segmentsNotReady"
	BreakStatusComplete   = "complete"
	DateRangeName         = "DateRange"
	DateRangeUPIDData     = "SegmentationUPIDData"
	breakNotReadyLimit    = 20 * time.Millisecond
)

var (
	DateRangeTag       = "#EXT-X-DATERANGE"
	SkipTag            = "#EXT-X-SKIP"             //todo: has attributes
	PreLoadHintTag     = "#EXT-X-PRELOAD-HINT"     //todo: has attributes
	RenditionReportTag = "#EXT-X-RENDITION-REPORT" //todo: has attributes

	// Attribute X-<client-attribute> is a client-specific attribute and new ones must be added manually below (e.g., X-ASSET-URI)
	dateRangeOrderAttr       = []string{"ID", "CLASS", "START-DATE", "END-DATE", "DURATION", "PLANNED-DURATION", "X-ASSET-URI", "SCTE35-OUT", "SCTE35-IN"}
	dateRangeShouldQuoteAttr = map[string]bool{
		"ID":               true,
		"CLASS":            true,
		"START-DATE":       true,
		"END-DATE":         true,
		"DURATION":         false,
		"PLANNED-DURATION": false,
		"X-ASSET-URI":      true,
		"SCTE35-OUT":       false,
		"SCTE35-IN":        false,
	}
)

type DateRangeParser struct{}

type DateRangeEncoder struct{}

func (p DateRangeParser) Parse(tag string, playlist *pl.Playlist) error {
	params := pl.TagsToMap(tag)
	if len(params) < 1 {
		return fmt.Errorf("invalid date range tag: %s", tag)
	}

	dateRangeNode := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name:  DateRangeName,
			Attrs: params,
		},
	}

	// An EXT-X-DATERANGE SCTE35-OUT tag signals the start of an Ad Break
	if dateRangeNode.HLSElement.Attrs["SCTE35-OUT"] != "" {
		mediaSequence, status := getAdBreakDetails(playlist, dateRangeNode)
		dateRangeNode.HLSElement.Details = map[string]string{
			"StartMediaSequence": mediaSequence,
			"Status":             status,
		}

		if upidData, ok := extractSegmentationUPIDData(dateRangeNode.HLSElement.Attrs["SCTE35-OUT"]); ok {
			dateRangeNode.HLSElement.Details[DateRangeUPIDData] = upidData
		}
	}

	playlist.Insert(dateRangeNode)
	return nil
}

func (e DateRangeEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	return pl.EncodeTagWithAttributes(builder, DateRangeTag, node.HLSElement.Attrs, dateRangeOrderAttr, dateRangeShouldQuoteAttr)
}

// Returns the Ad Break's media sequence (string) and status (string).
//   - The Break's media sequence will be the media sequence of the first segment inside the break (or zero if Break is incomplete).
//   - The Break's status will be: "complete" or incomplete ("leavingDVRLimit" or "segmentsNotReady").
func getAdBreakDetails(playlist *pl.Playlist, dateRangeNode *internal.Node) (value, status string) {
	currentMediaSequence := strconv.Itoa(playlist.MediaSequence + playlist.SegmentsCounter)
	breakStartDate, _ := time.Parse(time.RFC3339Nano, dateRangeNode.HLSElement.Attrs["START-DATE"])

	// when ad break segments are leaving DVR, we lose the break's first segment's media sequence
	if playlist.ProgramDateTime.IsZero() {
		// if the playlist's PDT tag was not parsed yet, we check if there are any media segments before the date range tag
		if len(playlist.Segments()) == 0 {
			log.Debug().Str("service", "go-m3u8/tags/media/metadata.go").Msg("ad break leaving dvr limit")
			return "0", BreakStatusLeavingDVR
		}
	} else {
		// if the playlist's PDT tag was already parsed, we check if the playlist PDT is equal or higher than the break's start date
		if playlist.ProgramDateTime.Equal(breakStartDate) || playlist.ProgramDateTime.After(breakStartDate) {
			log.Debug().Str("service", "go-m3u8/tags/media/metadata.go").Msg("ad break leaving dvr limit")
			return "0", BreakStatusLeavingDVR
		}
	}

	// when date range tag exists, but we don't know if we have the break's first media segment yet
	// we check if the break's start date comes later than the estimated next segment's PDT
	nextSegmentEstimatedPDT := playlist.ProgramDateTime.Add(time.Duration(playlist.DVR * float64(time.Second)))
	breakStartIsAfterNextSegment := breakStartDate.After(nextSegmentEstimatedPDT)

	// due to precision issues, we accept a small time difference of +/- 1ms
	// between the break's start date and the next segment's estimated PDT
	timeDifference := nextSegmentEstimatedPDT.Sub(breakStartDate)

	if breakStartIsAfterNextSegment && timeDifference.Abs() > breakNotReadyLimit {
		log.Debug().Str("service", "go-m3u8/tags/media/metadata.go").Msg("ad break not ready yet")
		return "0", BreakStatusNotReady
	}

	return currentMediaSequence, BreakStatusComplete
}

// Extracts segmentation_upid().Data from the SCTE-35 segmentation descriptor found.
// Returned value is raw string decoded from UPID bytes.
func extractSegmentationUPIDData(scte35 string) (string, bool) {
	payload := strings.TrimSpace(scte35)
	payload = strings.TrimPrefix(payload, "0x")
	payload = strings.TrimPrefix(payload, "0X")

	if len(payload) < 6 || len(payload)%2 != 0 {
		return "", false
	}

	binaryPayload, err := hex.DecodeString(payload)
	if err != nil || len(binaryPayload) < 17 {
		return "", false
	}

	idx := 0
	idx++ // table_id

	if idx+2 > len(binaryPayload) {
		return "", false
	}

	sectionLength := int(uint16(binaryPayload[idx]&0x0F)<<8 | uint16(binaryPayload[idx+1]))
	idx += 2

	sectionEnd := idx + sectionLength
	if sectionEnd > len(binaryPayload) {
		return "", false
	}

	if idx+5 > sectionEnd {
		return "", false
	}

	idx++    // protocol_version
	idx += 5 // encrypted_packet + encryption_algorithm + pts_adjustment

	if idx+6 > sectionEnd {
		return "", false
	}

	idx++ // cw_index

	// tier (12 bits) + splice_command_length (12 bits)
	if idx+2 >= len(binaryPayload) {
		return "", false
	}
	spliceCommandLength := int(uint16(binaryPayload[idx+1]&0x0F)<<8 | uint16(binaryPayload[idx+2]))
	idx += 3

	if idx >= sectionEnd {
		return "", false
	}

	idx++ // splice_command_type

	if idx+spliceCommandLength > sectionEnd {
		return "", false
	}
	idx += spliceCommandLength

	if idx+2 > sectionEnd {
		return "", false
	}

	descriptorLoopLength := int(uint16(binaryPayload[idx])<<8 | uint16(binaryPayload[idx+1]))
	idx += 2

	descriptorEnd := idx + descriptorLoopLength
	if descriptorEnd > sectionEnd {
		return "", false
	}

	for idx+2 <= descriptorEnd {
		descriptorTag := binaryPayload[idx]
		descriptorLength := int(binaryPayload[idx+1])
		idx += 2

		if idx+descriptorLength > descriptorEnd {
			return "", false
		}

		if descriptorTag == 0x02 {
			if upidData, found := parseSegmentationDescriptorUPID(binaryPayload[idx : idx+descriptorLength]); found {
				return upidData, true
			}
		}

		idx += descriptorLength
	}

	return "", false
}

func parseSegmentationDescriptorUPID(descriptor []byte) (string, bool) {
	if len(descriptor) < 11 {
		return "", false
	}

	idx := 0
	idx += 4 // identifier, e.g. CUEI
	idx += 4 // segmentation_event_id

	segmentationEventCancelIndicator := descriptor[idx]&0x80 != 0
	idx++
	if segmentationEventCancelIndicator {
		return "", false
	}

	if idx >= len(descriptor) {
		return "", false
	}

	flags := descriptor[idx]
	idx++

	programSegmentationFlag := flags&0x80 != 0
	segmentationDurationFlag := flags&0x40 != 0
	deliveryNotRestrictedFlag := flags&0x20 != 0

	if !deliveryNotRestrictedFlag {
		if idx >= len(descriptor) {
			return "", false
		}
		idx++
	}

	if !programSegmentationFlag {
		if idx >= len(descriptor) {
			return "", false
		}

		componentCount := int(descriptor[idx])
		idx++

		componentBytesLength := componentCount * 6
		if idx+componentBytesLength > len(descriptor) {
			return "", false
		}
		idx += componentBytesLength
	}

	if segmentationDurationFlag {
		if idx+5 > len(descriptor) {
			return "", false
		}
		idx += 5
	}

	if idx+2 > len(descriptor) {
		return "", false
	}

	segmentationUPIDLength := int(descriptor[idx+1])
	idx += 2

	if segmentationUPIDLength <= 0 || idx+segmentationUPIDLength > len(descriptor) {
		return "", false
	}

	segmentationUPIDData := descriptor[idx : idx+segmentationUPIDLength]
	return string(segmentationUPIDData), true
}
