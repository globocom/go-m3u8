package playlist

import (
	"errors"
	"fmt"
	"regexp"
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

// Returns new Playlist instance with an empty doubly linked list
func NewPlaylist() *Playlist {
	return &Playlist{
		DoublyLinkedList:      new(internal.DoublyLinkedList),
		CurrentSegment:        nil,
		CurrentStreamInf:      nil,
		ProgramDateTime:       time.Time{},
		MediaSequence:         0,
		DiscontinuitySequence: 0,
		SegmentsCounter:       0,
		DVR:                   0,
	}
}

// Prints the playlist to stdout for debugging purposes
func (p *Playlist) Print() {
	if p.Head == nil || p.Tail == nil {
		log.Warn().Str("service", "go-m3u8/playlist.go").Msg("playlist is empty")
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

// Returns the Version (#EXT-X-VERSION) tag's value as a string
func (p *Playlist) VersionValue() string {
	node, found := p.Find("Version")
	if !found {
		return ""
	}
	return node.HLSElement.Attrs["#EXT-X-VERSION"]
}

// Returns the Version (#EXT-X-VERSION) tag as a Node if it exists, otherwise returns nil and false
func (p *Playlist) VersionTag() (*internal.Node, bool) {
	return p.Find("Version")
}

// Returns the MediaSequence (#EXT-X-MEDIA-SEQUENCE) tag's value as a string
func (p *Playlist) MediaSequenceValue() string {
	node, found := p.Find("MediaSequence")
	if !found {
		return ""
	}
	return node.HLSElement.Attrs["#EXT-X-MEDIA-SEQUENCE"]
}

// Returns the MediaSequence (#EXT-X-MEDIA-SEQUENCE) tag as a Node if it exists, otherwise returns nil and false
func (p *Playlist) MediaSequenceTag() (*internal.Node, bool) {
	return p.Find("MediaSequence")
}

// Returns the DiscontinuitySequence (#EXT-X-DISCONTINUITY-SEQUENCE) tag's value as a string
func (p *Playlist) DiscontinuitySequenceValue() string {
	node, found := p.Find("DiscontinuitySequence")
	if !found {
		return ""
	}
	return node.HLSElement.Attrs["#EXT-X-DISCONTINUITY-SEQUENCE"]
}

// Returns the DiscontinuitySequence (#EXT-X-DISCONTINUITY-SEQUENCE) tag as a Node if it exists, otherwise returns nil and false
func (p *Playlist) DiscontinuitySequenceTag() (*internal.Node, bool) {
	return p.Find("DiscontinuitySequence")
}

// Returns the VariableDefine (#EXT-X-DEFINE) tag as a Node if it exists, otherwise returns nil and false
func (p *Playlist) VariableDefineTag() (*internal.Node, bool) {
	return p.Find("VariableDefine")
}

// Returns all StreamInf (#EXT-X-STREAM-INF) nodes in the playlist
func (p *Playlist) Variants() []*internal.Node {
	return p.FindAll("StreamInf")
}

// Returns all Media (#EXT-X-MEDIA) nodes in the playlist (i.e. AUDIO groups, CLOSED-CAPTIONS groups, etc.)
func (p *Playlist) MediaGroups() []*internal.Node {
	return p.FindAll("Media")
}

// Returns all IFrameStreamInf (#EXT-X-I-FRAME-STREAM-INF) nodes in the playlist (i.e. keyframes)
func (p *Playlist) Keyframes() []*internal.Node {
	return p.FindAll("IFrameStreamInf")
}

// Returns all ExtInf (#EXTINF) nodes in the playlist
func (p *Playlist) Segments() []*internal.Node {
	return p.FindAll("ExtInf")
}

// Returns all Key (#EXT-X-KEY) nodes in the playlist
func (p *Playlist) EncryptionTags() []*internal.Node {
	return p.FindAll("Key")
}

// Returns all CueOut (#EXT-X-CUE-OUT) nodes in the playlist
func (p *Playlist) CueOutEvents() []*internal.Node {
	return p.FindAll("CueOut")
}

// Returns all CueIn (#EXT-X-CUE-IN) nodes in the playlist
func (p *Playlist) CueInEvents() []*internal.Node {
	return p.FindAll("CueIn")
}

// Returns all DateRange (#EXT-X-DATERANGE) nodes with SCTE35-OUT marking in the playlist
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

// Returns all DateRange (#EXT-X-DATERANGE) nodes with SCTE35-IN marking in the playlist
func (p *Playlist) SCTE35InTags() []*internal.Node {
	result := make([]*internal.Node, 0)
	nodes := p.FindAll("DateRange")
	for _, node := range nodes {
		if node.HLSElement.Attrs["SCTE35-IN"] != "" {
			result = append(result, node)
		}
	}
	return result
}

// Returns the first Comment node in the playlist whose value contains the given matchString.
//
//	Example: "# variants", "# AUDIO groups", etc
func (p *Playlist) Comment(matchString string) *internal.Node {
	nodes := p.FindAll("Comment")
	for _, node := range nodes {
		if strings.Contains(node.HLSElement.Attrs["Comment"], matchString) {
			return node
		}
	}
	return nil
}

// Returns true if node is inside ad break and false otherwise.
// When true, method also returns the DateRange (#EXT-X-DATERANGE) object for the Ad Break.
//
// For entering the Ad Break, we always have DateRange tag with SCTE-OUT and CueOut (#EXT-X-CUE-OUT) tag.
// However, for exiting the Ad Break, we have three possible manifests:
//
//   - DateRange SCTE-IN is ALWAYS present.
//   - No DateRange SCTE-IN. Exit is ONLY marked by CueIn (#EXT-X-CUE-IN) tag instead.
//   - SOMETIMES DateRange SCTE-IN is present, alongside the CueIn tag.
func (p *Playlist) FindNodeInsideAdBreak(node *internal.Node) (*internal.Node, bool) {
	current := node.Prev
	for current != nil {
		// node is inside Ad Break if it is preceded by a DateRange tag with attribute SCTE35-OUT
		if (current.HLSElement.Name == "DateRange") && (current.HLSElement.Attrs["SCTE35-OUT"] != "") {
			return current, true
		}

		// node is outside Ad Break if it is preceded by a CueIn tag or a DateRange tag with attribute SCTE35-IN
		if (current.HLSElement.Name == "CueIn") || (current.HLSElement.Name == "DateRange" && current.HLSElement.Attrs["SCTE35-IN"] != "") {
			return nil, false
		}

		current = current.Prev
	}

	return nil, false
}

// Returns the last DateRange (#EXT-X-DATERANGE) node with SCTE35-OUT marking (i.e. the last Ad Break) in the playlist
func (p *Playlist) FindLastAdBreak() (*internal.Node, bool) {
	adBreaks := p.Breaks()
	if len(adBreaks) == 0 {
		return nil, false
	}
	return adBreaks[len(adBreaks)-1], true
}

// DuplicateAdBreak checks if two ad breaks have the same START-DATE, indicating a duplicate.
// Two ad breaks are considered duplicates if they share the same START-DATE
// and the same PLANNED-DURATION.
func (p *Playlist) HasDuplicateAdBreak() bool {
	adBreaks := p.Breaks()
	if len(adBreaks) < 2 {
		return false
	}

	lastBreak := adBreaks[len(adBreaks)-1]
	previousBreak := adBreaks[len(adBreaks)-2]

	lastBreakStartDate := lastBreak.HLSElement.Attrs["START-DATE"]
	previousBreakStartDate := previousBreak.HLSElement.Attrs["START-DATE"]

	lastBreakDuration := lastBreak.HLSElement.Attrs["PLANNED-DURATION"]
	previousBreakDuration := previousBreak.HLSElement.Attrs["PLANNED-DURATION"]

	return lastBreakStartDate == previousBreakStartDate && lastBreakDuration == previousBreakDuration
}
