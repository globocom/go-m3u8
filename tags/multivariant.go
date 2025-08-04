//	Multivariant Playlist Tags (Section 4.4.6 on RFC)
//
// Multivariant Playlist tags define the Variant Streams, Renditions,
// and other global parameters of the presentation.
//
// Multivariant Playlist tags MUST NOT appear in a Media Playlist;
// clients MUST fail to parse any Playlist that contains both a
// Multivariant Playlist tag and either a Media Playlist tag or a Media
// Segment tag.
package tags

import (
	"fmt"
	"strings"

	"github.com/globocom/go-m3u8/internal"
	pl "github.com/globocom/go-m3u8/playlist"
)

const (
	StreamInfName       = "StreamInf"
	MediaName           = "Media"
	IFrameStreamInfName = "IFrameStreamInf"
)

var (
	StreamInfTag       = "#EXT-X-STREAM-INF"
	MediaTag           = "#EXT-X-MEDIA"
	IFrameStreamInfTag = "#EXT-X-I-FRAME-STREAM-INF"
	SessionKey         = "#EXT-X-SESSION-KEY" //todo: has attributes
)

type (
	StreamInfParser       struct{}
	MediaParser           struct{}
	IFrameStreamInfParser struct{}
)

type (
	StreamInfEncoder       struct{}
	MediaEncoder           struct{}
	IFrameStreamInfEncoder struct{}
)

func (p StreamInfParser) Parse(tag string, playlist *pl.Playlist) error {
	playlist.CurrentStreamInf = pl.GetStreamInfData(pl.TagsToMap(tag))
	return nil
}

func (p MediaParser) Parse(tag string, playlist *pl.Playlist) error {
	params := pl.TagsToMap(tag)
	if len(params) < 1 {
		return fmt.Errorf("invalid media tag: %s", tag)
	}

	// The TYPE attribute is REQUIRED by RFC
	if params["TYPE"] == "" {
		return fmt.Errorf("TYPE attribute is required: %s", tag)
	}

	// Valid strings for TYPE are AUDIO, VIDEO, SUBTITLES, and CLOSED-CAPTIONS.
	if params["TYPE"] != "AUDIO" && params["TYPE"] != "VIDEO" && params["TYPE"] != "SUBTITLES" && params["TYPE"] != "CLOSED-CAPTIONS" {
		return fmt.Errorf("invalid TYPE attribute value: %s", params["TYPE"])
	}

	// The GROUP-ID attribute is REQUIRED by RFC
	if params["GROUP-ID"] == "" {
		return fmt.Errorf("GROUP-ID attribute is required: %s", tag)
	}

	// If the TYPE is CLOSED-CAPTIONS, the URI attribute MUST NOT be present
	if params["TYPE"] == "CLOSED-CAPTIONS" && params["URI"] != "" {
		return fmt.Errorf("URI attribute is not allowed for CLOSED-CAPTIONS type: %s", tag)
	}

	mediaNode := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name:  "Media",
			Attrs: params,
		},
	}
	playlist.Insert(mediaNode)

	return nil
}

func (p IFrameStreamInfParser) Parse(tag string, playlist *pl.Playlist) error {
	params := pl.TagsToMap(tag)
	if len(params) < 1 {
		return fmt.Errorf("invalid IFrameStreamInf tag: %s", tag)
	}

	// The BANDWIDTH attribute is REQUIRED by RFC
	if params["BANDWIDTH"] == "" {
		return fmt.Errorf("BANDWIDTH attribute is required: %s", tag)
	}

	// The CODECS attribute is REQUIRED by RFC
	if params["CODECS"] == "" {
		return fmt.Errorf("CODECS attribute is required: %s", tag)
	}

	// The URI attribute is REQUIRED by RFC
	if params["URI"] == "" {
		return fmt.Errorf("URI attribute is required: %s", tag)
	}

	IFrameStreamInfNode := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name:  "IFrameStreamInf",
			Attrs: params,
		},
	}
	playlist.Insert(IFrameStreamInfNode)
	return nil
}

func (e StreamInfEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	orderAttr := []string{"BANDWIDTH", "AVERAGE-BANDWIDTH", "CODECS", "RESOLUTION", "FRAME-RATE", "VIDEO-RANGE", "AUDIO", "VIDEO", "SUBTITLES", "CLOSED-CAPTIONS"}
	shouldQuoteAttr := e.shouldQuoteStreamInf(node)

	if err := pl.EncodeTagWithAttributes(builder, StreamInfTag, node.HLSElement.Attrs, orderAttr, shouldQuoteAttr); err != nil {
		return err
	}
	if node.HLSElement.URI != "" {
		_, err := builder.WriteString(node.HLSElement.URI + "\n")
		return err
	}
	return nil
}

func (e MediaEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	orderAttr := []string{"TYPE", "GROUP-ID", "LANGUAGE", "NAME", "DEFAULT", "AUTOSELECT", "CHANNELS", "URI"}
	shouldQuoteAttr := map[string]bool{
		"TYPE":       false,
		"GROUP-ID":   true,
		"LANGUAGE":   true,
		"NAME":       true,
		"DEFAULT":    false,
		"AUTOSELECT": false,
		"CHANNELS":   true,
		"URI":        true,
	}

	return pl.EncodeTagWithAttributes(builder, MediaTag, node.HLSElement.Attrs, orderAttr, shouldQuoteAttr)
}

func (e IFrameStreamInfEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	orderAttr := []string{"BANDWIDTH", "AVERAGE-BANDWIDTH", "CODECS", "RESOLUTION", "URI", "VIDEO-RANGE", "VIDEO", "SCORE"}
	shouldQuoteAttr := map[string]bool{
		"BANDWIDTH":         false,
		"AVERAGE-BANDWIDTH": false,
		"CODECS":            true,
		"RESOLUTION":        false,
		"URI":               true,
		"VIDEO-RANGE":       false,
		"VIDEO":             true,
		"SCORE":             false,
	}

	return pl.EncodeTagWithAttributes(builder, IFrameStreamInfTag, node.HLSElement.Attrs, orderAttr, shouldQuoteAttr)
}

func (e StreamInfEncoder) shouldQuoteStreamInf(node *internal.Node) map[string]bool {
	shouldQuoteAttr := map[string]bool{
		"BANDWIDTH":         false,
		"AVERAGE-BANDWIDTH": false,
		"CODECS":            true,
		"RESOLUTION":        false,
		"FRAME-RATE":        false,
		"VIDEO-RANGE":       false,
		"AUDIO":             true,
		"VIDEO":             true,
		"SUBTITLES":         true,
		"CLOSED-CAPTIONS":   true,
	}

	// the value can be either a quoted-string or an enumerated-string with the value NONE
	if node.HLSElement.Attrs["CLOSED-CAPTIONS"] == "NONE" {
		shouldQuoteAttr["CLOSED-CAPTIONS"] = false
	}

	return shouldQuoteAttr
}
