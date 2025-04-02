package go_m3u8

import (
	"time"

	"github.com/globocom/go-m3u8/internal"
)

type Playlist struct {
	*internal.DoublyLinkedList
	CurrentDateRange *internal.Node
	CurrentSegment   *Segment
	CurrentStreamInf *StreamInf
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
	return node.Attrs["#EXT-X-VERSION"]
}

func (p *Playlist) Version() (*internal.Node, bool) {
	return p.Find("Version")
}

func (p *Playlist) MediaSequenceValue() string {
	node, found := p.Find("MediaSequence")
	if !found {
		return ""
	}
	return node.Attrs["#EXT-X-MEDIA-SEQUENCE"]
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
		if node.Attrs[scteOutAttribute] != "" {
			result = append(result, node)
		}
	}
	return result
}

func (p *Playlist) ReplaceBreaksURI(transform func(string) string) error {
	startCondition := func(node *internal.Node) bool {
		return node.Name == "DateRange" && node.Attrs["SCTE35-OUT"] != ""
	}
	endCondition := func(node *internal.Node) bool {
		return node.Name == "DateRange" && node.Attrs["SCTE35-IN"] != ""
	}
	transformFunc := func(node *internal.Node) {
		if node.Name == "ExtInf" && node.URI != "" {
			node.URI = transform(node.URI)
		}
	}
	return p.ModifyNodesBetween(startCondition, endCondition, transformFunc)
}
