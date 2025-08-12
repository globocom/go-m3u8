//	Media Metadata Tags (Section 4.4.5 on RFC)
//
// Media Metadata tags provide information about the playlist that is
// not associated with specific Media Segments.  There MAY be more than
// one Media Metadata tag of each type in any Media Playlist.  The only
// exception to this rule is EXT-X-SKIP, which MUST NOT appear more than
// once.
package tags

import (
	"fmt"
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
	breakNotReadyLimit    = 20 * time.Millisecond
)

var (
	DateRangeTag       = "#EXT-X-DATERANGE"
	SkipTag            = "#EXT-X-SKIP"             //todo: has attributes
	PreLoadHintTag     = "#EXT-X-PRELOAD-HINT"     //todo: has attributes
	RenditionReportTag = "#EXT-X-RENDITION-REPORT" //todo: has attributes
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
	}

	playlist.Insert(dateRangeNode)
	return nil
}

func (e DateRangeEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	// Attribute X-<client-attribute> is a client-specific attribute and new ones must be added manually below (e.g., X-ASSET-URI)
	orderAttr := []string{"ID", "CLASS", "START-DATE", "END-DATE", "DURATION", "PLANNED-DURATION", "X-ASSET-URI", "SCTE35-OUT", "SCTE35-IN"}
	shouldQuoteAttr := map[string]bool{
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
	return pl.EncodeTagWithAttributes(builder, DateRangeTag, node.HLSElement.Attrs, orderAttr, shouldQuoteAttr)
}

// Returns the Ad Break's media sequence (string) and status (string).
//   - The Break's media sequence will be the media sequence of the first segment inside the break (or zero if Break is incomplete).
//   - The Break's status will be: "complete" or incomplete ("leavingDVRLimit" or "segmentsNotReady").
func getAdBreakDetails(playlist *pl.Playlist, dateRangeNode *internal.Node) (value, status string) {
	currentMediaSequence := fmt.Sprintf("%d", playlist.MediaSequence+playlist.SegmentsCounter)
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
