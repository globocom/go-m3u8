//	Media Metadata Tags (Section 4.4.5 on RFC)
//
// Media Metadata tags provide information about the playlist that is
// not associated with specific Media Segments.  There MAY be more than
// one Media Metadata tag of each type in any Media Playlist.  The only
// exception to this rule is EXT-X-SKIP, which MUST NOT appear more than
// once.
package media

import (
	"fmt"
	"strings"
	"time"

	"github.com/globocom/go-m3u8/internal"
	pl "github.com/globocom/go-m3u8/playlist"
	"github.com/rs/zerolog/log"
)

var (
	DateRangeTag   = "#EXT-X-DATERANGE"    //todo: has attributes
	SkipTag        = "#EXT-X-SKIP"         //todo: has attributes
	PreLoadHintTag = "#EXT-X-PRELOAD-HINT" //todo: has attributes
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
			Name:    "DateRange",
			Attrs:   params,
			Details: map[string]string{},
		},
	}

	// date range tag with SCTE35-OUT attribute (i.e. indicates ad break start)
	if dateRangeNode.HLSElement.Attrs["SCTE35-OUT"] != "" {
		dateRangeNode.HLSElement.Details["FirstSegmentMediaSequence"] = getBreakMediaSequence(playlist, dateRangeNode)
	}

	playlist.Insert(dateRangeNode)
	return nil
}

func (e DateRangeEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	order := []string{"ID", "START-DATE", "PLANNED-DURATION", "END-DATE", "DURATION", "SCTE35-OUT", "SCTE35-IN"}
	return pl.EncodeTagWithAttributes(builder, DateRangeTag, node.HLSElement.Attrs, order)
}

// An EXT-X-DATERANGE SCTE35-OUT tag signals the start of an Ad Break
// The Ad Break's media sequence will be the media sequence of the first segment inside the break
// If we don't have information
func getBreakMediaSequence(playlist *pl.Playlist, dateRangeNode *internal.Node) string {
	currentMediaSequence := fmt.Sprintf("%d", playlist.MediaSequence+playlist.SegmentsCounter)

	// when ad break segments are leaving DVR limit, we lose the break's first segment's media sequence

	// if the playlist's PDT tag was parsed already, we check if the playlist PDT is equal or higher than the break's start date
	// in this case, the break is leaving the DVR and we will set the break media sequence to zero
	breakStartDate, _ := time.Parse(time.RFC3339Nano, dateRangeNode.HLSElement.Attrs["START-DATE"])
	if !playlist.ProgramDateTime.IsZero() {
		if playlist.ProgramDateTime.Equal(breakStartDate) || playlist.ProgramDateTime.After(breakStartDate) {
			log.Info().Msg("break leaving dvr limit, media sequence will be zero")
			return "0"
		}
	}

	// if the playlist's PDT tag is not parsed yet, we check if there are any media segments before the date range tag
	// if there aren't, the break is leaving the DVR and we will set the break media sequence to zero
	parsedSegmentsBeforeBreak := playlist.Segments()
	if len(parsedSegmentsBeforeBreak) == 0 {
		log.Info().Msg("no media segments present before ad break, media sequence will be zero")
		return "0"
	}

	// when date range tag exists, but we don't know if we have the break's first media segment yet

	// we check if the break's start date comes later than the estimated next segment's PDT
	// in this case, we will set the break media sequence to zero
	nextSegmentEstimatedPDT := playlist.ProgramDateTime.Add(time.Duration(playlist.DVR * float64(time.Second)))
	if (breakStartDate.Round(time.Millisecond)).After(nextSegmentEstimatedPDT.Round(time.Millisecond)) {
		log.Info().Msg("ad break does not contain segments yet, media sequence will be zero")
		return "0"
	}

	return currentMediaSequence
}
