<div align=center><img src="gopher.png" height=100px></div>
<h1 align=center>
go-m3u8
<div align=center>
<img src="https://github.com/globocom/go-m3u8/actions/workflows/go.yml/badge.svg">
<a href="https://goreportcard.com/report/github.com/globocom/go-m3u8"><img src="https://goreportcard.com/badge/github.com/globocom/go-m3u8"/></a>
<img src="https://img.shields.io/github/go-mod/go-version/globocom/go-m3u8">
</div>
</h1>


### ⚠️ Work in Progress!

_This project is currently in pre-release and is subject to changes._

## About

A m3u8 parser library in Go to facilitate the manipulation of HTTP Live Streaming (HLS) manifests.

The library provides:

- **Doubly Linked List Architecture -** Efficient playlist representation with ordered element access and manipulation.
- **RFC 8216 Compliance -** Support for HLS tags based on the official [RFC](https://datatracker.ietf.org/doc/html/draft-pantos-hls-rfc8216bis) documentation.
- **Live Streaming Support -** Optimized for live streaming manifests only (no VoD support yet).

## Architecture

### Doubly Linked List

The m3u8 manifest is parsed into a **doubly linked list**. Each node on the list represents an element of the original HLS playlist (e.g. a tag, a comment, etc).

This data structure allows us to access the manifest in a sorted manner, to retrieve information and apply operations to its content.

Some available operations are:

- **Create** a new node, to represent an element of the playlist.
- **Insert** a new node into the playlist, either at the end or before/after an existing node.
- **Find** a specific node or a list of nodes, based on the element name.

### HLS Elements

To guarantee scalibility, our library considers a data structure for the HLS elements that follows the [RFC documentation](https://tools.ietf.org/html/rfc8216). 

The [**tags**](/tags) package implements the currently supported [Playlist Tags](https://datatracker.ietf.org/doc/html/draft-pantos-hls-rfc8216bis#section-4.4):

1. **basic -** Basic Tags (Section 4.4.1).
- `#EXTM3U`
- `#EXT-X-VERSION`

2. **exclusive -** Media or Multivariant Playlist Tags (Section 4.4.2).
- `#EXT-X-INDEPENDENT-SEGMENTS`
- `#EXT-X-DEFINE`

3. **media -** Media Playlist, Metadata and Segment Tags (Sections 4.4.3 to 4.4.5).
- `#EXT-X-DATERANGE`
- `#EXT-X-TARGETDURATION`
- `#EXT-X-MEDIA-SEQUENCE`
- `#EXT-X-DISCONTINUITY-SEQUENCE`
- `#EXTINF`
- `#EXT-X-DISCONTINUITY`
- `#EXT-X-PROGRAM-DATE-TIME`
- `#EXT-X-KEY`

4. **multivariant -** Multivariant Playlist Tags (Section 4.4.6).
- `#EXT-X-STREAM-INF`

5. **others -** The tags in this section are "non-official" and are not listed in the RFC, e.g. tags added to the manifest by the live stream packaging service.
- `#EXT-X-CUE-OUT`
- `#EXT-X-CUE-IN`
- Packager specific tags.
- In-manifest comments (begin with `#` and are NOT tags).

## Getting Started

### Installation

To import the library to your Go project, run:

```sh
go install github.com/globocom/go-m3u8
```

The following [**Makefile**](/Makefile) commands are available:

```Makefile
make test 		# Run test suites
make lint 		# Run code linter
```

The [**testlocal**](/testlocal/) folder contains instructions on how to setup and test the library for local experimentation.

### Decoding a Playlist

The `ParsePlaylist` method receives a `io.ReadCloser` object as argument and returns a `Playlist` object.

You may decode a manifest that is in string format:
```go
manifest := `#EXTM3U
#EXT-X-VERSION:3
# variants
#EXT-X-STREAM-INF:BANDWIDTH=479000,AVERAGE-BANDWIDTH=435000,CODECS="mp4a.40.2,avc1.64001F",RESOLUTION=512x288,FRAME-RATE=30
channel_01.m3u8`

manifestReader := io.NopCloser(strings.NewReader(manifest))

playlist, err := m3u8.ParsePlaylist(manifestReader)
if err != nil {
	panic(err)
}
```

Or read the manifest file directly:
```go
file, _ := os.Open("multivariant.m3u8")

playlist, err := go_m3u8.ParsePlaylist(file)
if err != nil {
	panic(err)
}
```

### Encoding a Playlist

The `EncodePlaylist` method parses a `Playlist` object back into string format.

```go
manifest, err := go_m3u8.EncodePlaylist(playlist)
if err != nil {
	panic(err)
}
```

## Usage 

For complete details on the available methods, please read [the original release notes](https://github.com/globocom/go-m3u8/releases/tag/v0.1.0).

### Collecting Ad Break Data

Collect information on ad breaks present on the manifest when SCTE-35 ad insertion is used.

```go
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	go_m3u8 "github.com/globocom/go-m3u8"
	m3u8_pl "github.com/globocom/go-m3u8/playlist"
	m3u8_tags "github.com/globocom/go-m3u8/tags"
)

type AdBreak struct {
	Timestamp     string
	MediaSequence int
	Status        string
}

func main() {
	file, _ := os.Open("playlist.m3u8")
	p, err := go_m3u8.ParsePlaylist(file)

	if err != nil {
		panic(err)
	}

	latestBreak := GetLatestBreakData(p)

	fmt.Printf("%+v\n", latestBreak)
}

func GetLatestBreakData(manifest *m3u8_pl.Playlist) AdBreak {
	adBreaks := manifest.Breaks()

	if len(adBreaks) == 0 {
		return AdBreak{}
	}

	latestAdBreak := adBreaks[len(adBreaks)-1]

	breakStatus := latestAdBreak.HLSElement.Details["Status"]
	if (breakStatus == m3u8_tags.BreakStatusLeavingDVR) || (breakStatus == m3u8_tags.BreakStatusNotReady) {
		return AdBreak{}
	}

	if latestAdBreak.HLSElement.Attrs["PLANNED-DURATION"] == "" {
		return AdBreak{Status: "invalid"}
	}

	startMediaSequence, err := strconv.Atoi(latestAdBreak.HLSElement.Details["StartMediaSequence"])
	if err != nil {
		return AdBreak{Status: "invalid"}
	}

	startDate, err := time.Parse(time.RFC3339Nano, latestAdBreak.HLSElement.Attrs["START-DATE"])
	if err != nil {
		return AdBreak{Status: "invalid"}
	}

	adBreak := AdBreak{
		MediaSequence: startMediaSequence,
		Timestamp:     fmt.Sprintf("%d", startDate.Unix()),
		Status:        "valid",
	}

	return adBreak
}
```

### Adding Ad Break Markers

Insert SCTE-35 ad break markers at specific segments.

```go
package main

import (
	"os"

	go_m3u8 "github.com/globocom/go-m3u8"
	m3u8_tags "github.com/globocom/go-m3u8/tags"
)

func main() {
	file, _ := os.Open("playlist.m3u8")
	p, err := go_m3u8.ParsePlaylist(file)
	if err != nil {
		panic(err)
	}

	// Find all segments from media playlist
	// Suppose there are 12 segments, each with a duration of 3.2s
	segments := p.Segments()

	// Find a specific segment to insert ad break before
	segmentNodeBreakStart := segments[2]

	// Create date range tag with SCTE35
	attrs := map[string]string{
		"ID":               "1-1747402436",
		"START-DATE":       "2025-05-16T13:33:56.266666Z",
		"PLANNED-DURATION": "16.0",
		"SCTE-OUT":         "0x123abc456def",
	}
	dateRangeNode := p.NewNode(m3u8_tags.DateRangeName, "", attrs, nil)

	// Create event cue out tag
	cueOutNode := p.NewNode(m3u8_tags.EventCueOutName, "", map[string]string{m3u8_tags.EventCueOutTag: "16.0"}, nil)

	// Insert ad break start
	p.InsertBefore(segmentNodeBreakStart, dateRangeNode)
	p.InsertAfter(dateRangeNode, cueOutNode)

	// Find the last segment inside ad break
	segmentNodeBreakEnd := segments[8]

	// Create event cue in tag
	cueInNode := p.NewNode(m3u8_tags.EventCueInName, "", map[string]string{m3u8_tags.EventCueInTag: ""}, nil)

	// Insert ad break end
	p.InsertAfter(segmentNodeBreakEnd, cueInNode)

  // Encode the playlist back into manifest format
	manifest, err := go_m3u8.EncodePlaylist(p)
	if err != nil {
		panic(err)
	}

	print(manifest)
}
```

### Adding Discontinuity Information

Insert discontinuity tags when SCTE-35 ad break markers are present.

```go
package main

import (
	"os"

	go_m3u8 "github.com/globocom/go-m3u8"
	m3u8_tags "github.com/globocom/go-m3u8/tags"
)

func main() {
	file, _ := os.Open("playlist.m3u8")
	p, err := go_m3u8.ParsePlaylist(file)
	if err != nil {
		panic(err)
	}

	// Find all ad break markers (date range tags with scte-out & cue-in events)
	dateRangeNodes := p.Breaks()
	cueInNodes := p.CueInEvents()

	// Insert discontinuity before each ad break start
	for _, dateRangeNode := range dateRangeNodes {
		discontinuityNode := p.NewNode(m3u8_tags.DiscontinuityName, "", nil, nil)
		p.InsertBefore(dateRangeNode, discontinuityNode)
	}

	// Insert discontinuity after each ad break end
	for _, cueInNode := range cueInNodes {
		discontinuityNode := p.NewNode(m3u8_tags.DiscontinuityName, "", nil, nil)
		p.InsertBefore(cueInNode, discontinuityNode)
	}

	// Encode the playlist back into manifest format
	manifest, err := go_m3u8.EncodePlaylist(p)
	if err != nil {
		panic(err)
	}

	print(manifest)
}
```

### Updating Encryption Keys

Rotate encryption keys for enhanced security.

```go
package main

import (
	"os"
	go_m3u8 "github.com/globocom/go-m3u8"
)

func main() {
	file, _ := os.Open("playlist.m3u8")
	p, err := go_m3u8.ParsePlaylist(file)
	if err != nil {
		panic(err)
	}

	// Find all encryption key tags
	keyNodes := p.EncryptionTags()
	
	for _, keyNode := range keyNodes {
		// Update the key URI with new key server
		if keyNode.HLSElement.Attrs != nil {
			keyNode.HLSElement.Attrs["URI"] = "https://new-key-server.com/key.bin"
			keyNode.HLSElement.Attrs["IV"] = "0x12345678901234567890123456789012"
		}
	}

	p.Print()
}
```

## Contributing

As this is an open-source project, we encourage and support any community contributions!

Feel free to report bugs, request features and recommend improvements. Multiple repository mantainers are keeping an eye on opened issues and pull requests.

## License

Released under the [MIT License](/LICENSE).
