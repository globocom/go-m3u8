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

// Returns the Version tag's value as a string
func (p *Playlist) VersionValue() string {
	node, found := p.Find("Version")
	if !found {
		return ""
	}
	return node.HLSElement.Attrs["#EXT-X-VERSION"]
}

// Returns the Version tag as a Node if it exists, otherwise returns nil and false
func (p *Playlist) VersionTag() (*internal.Node, bool) {
	return p.Find("Version")
}

// Returns the MediaSequence tag's value as a string
func (p *Playlist) MediaSequenceValue() string {
	node, found := p.Find("MediaSequence")
	if !found {
		return ""
	}
	return node.HLSElement.Attrs["#EXT-X-MEDIA-SEQUENCE"]
}

// Returns the MediaSequence tag as a Node if it exists, otherwise returns nil and false
func (p *Playlist) MediaSequenceTag() (*internal.Node, bool) {
	return p.Find("MediaSequence")
}

// Returns the DiscontinuitySequence tag's value as a string
func (p *Playlist) DiscontinuitySequenceValue() string {
	node, found := p.Find("DiscontinuitySequence")
	if !found {
		return ""
	}
	return node.HLSElement.Attrs["#EXT-X-DISCONTINUITY-SEQUENCE"]
}

// Returns the DiscontinuitySequence tag as a Node if it exists, otherwise returns nil and false
func (p *Playlist) DiscontinuitySequenceTag() (*internal.Node, bool) {
	return p.Find("DiscontinuitySequence")
}

// Returns the Variable Define tag as a Node if it exists, otherwise returns nil and false
func (p *Playlist) VariableDefineTag() (*internal.Node, bool) {
	return p.Find("VariableDefine")
}

// Returns all StreamInf nodes in the playlist
func (p *Playlist) Variants() []*internal.Node {
	return p.FindAll("StreamInf")
}

// Returns all Media nodes in the playlist (i.e. AUDIO groups, CLOSED-CAPTIONS groups, etc.)
func (p *Playlist) MediaGroups() []*internal.Node {
	return p.FindAll("Media")
}

// Returns all IFrameStreamInf nodes in the playlist (i.e. keyframes)
func (p *Playlist) Keyframes() []*internal.Node {
	return p.FindAll("IFrameStreamInf")
}

// Returns all ExtInf nodes in the playlist
func (p *Playlist) Segments() []*internal.Node {
	return p.FindAll("ExtInf")
}

// Returns all Key nodes in the playlist
func (p *Playlist) EncryptionTags() []*internal.Node {
	return p.FindAll("Key")
}

// Returns all CueOut nodes in the playlist
func (p *Playlist) CueOutEvents() []*internal.Node {
	return p.FindAll("CueOut")
}

// Returns all CueIn nodes in the playlist
func (p *Playlist) CueInEvents() []*internal.Node {
	return p.FindAll("CueIn")
}

// Returns all DateRange nodes with SCTE35-OUT marking in the playlist
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

// Returns all DateRange nodes with SCTE35-IN marking in the playlist
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
// When true, method also returns the DateRange object for the Ad Break.
//
// For entering the Ad Break, we always have DateRange tag with SCTE-OUT and CueOutEvent tag.
// However, for exiting the Ad Break, we have three possible manifests:
//
//   - DateRange SCTE-IN is ALWAYS present.
//   - No DateRange SCTE-IN. Exit is ONLY marked by CueInEvent tag instead.
//   - SOMETIMES DateRange SCTE-IN is present, alongside the CueInEvent tag.
func (p *Playlist) FindNodeInsideAdBreak(node *internal.Node) (*internal.Node, bool) {
	current := node
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
