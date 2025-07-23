# Handling Ad Breaks: DATERANGE with SCTE-35 Marking

An Ad Break is present on the media playlist by way of SCTE-35 Marking. 

- Break **start** is marked by a `EXT-X-DATERANGE` tag with attribute `SCTE-OUT` AND a `EXT-X-CUE-OUT` tag. Both tags are ALWAYS present.
- Break **end** is marked by a `EXT-X-DATERANGE` tag with attribute `SCTE-IN` AND/OR a `EXT-X-CUE-IN` tag. One or both tags are ALWAYS present.
- The **start** and **end** tag(s) are usually followed by a PDT `EXT-X-PROGRAM-DATE-TIME` tag.
- Inside the Break, between the **start** and **end** markings, we have multiple the Break segments.

In some cases, it is important to know the media sequence of when the Ad Break starts (e.g. manifest manipulation for server-side dynamic ad insertion).
For this reason, the `EXT-X-DATERANGE` tag that marks the Break **start** (`tags/media/metadata.go`) will hold the following additional metadata:

- *Start Media Sequence:* This value is ALWAYS equal to the media sequence of the first segment inside the Break.
- *Status:* Determines whether the Break is complete (**start** and **end** tags + all Break segments in between are present).

Additionally, some manifests may have incomplete Ad Break elements at the time of request. In some cases, the Break **start** and **end** tags are present, but some Break segments have already left the manifest due to the DVR limit. In other cases, the Break **start** tags are at the end of the manifest, but the next Break segments don't exit yet.

This leaves us with two possible Ad Break scenarios to handle, when parsing the media playlist.

## 1. Complete Ad Breaks

Most manifests will contain *complete* Ad Breaks (one or more) that have clear start and finish markings and all Break segments present inside it.

**Test File:** `/mocks/media/scte35/withCompleteAdBreak.m3u8`
```
#EXTM3U
#EXT-X-VERSION:11
#EXT-X-MEDIA-SEQUENCE:363987560
#EXT-X-INDEPENDENT-SEGMENTS
#EXT-X-TARGETDURATION:7
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T12:44:44.233333Z
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363987560.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363987561.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363987562.ts
#EXTINF:5.9333, no desc
channel-audio_1=96000-video=3442944-363987563.ts
## splice_insert(auto_return)
#EXT-X-DATERANGE:ID="4026559039-1747140304",START-DATE="2025-05-13T12:45:04.566666Z",PLANNED-DURATION=20,SCTE35-OUT=0xFC3025000000000BB800FFF01405F0006A3F7FEFFE3775B370FE001B77400001010100001AC3CE61
#EXT-X-CUE-OUT:20
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T12:45:04.566666Z
#EXTINF:3.6666, no desc
channel-audio_1=96000-video=3442944-363987564.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363987565.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363987566.ts
#EXTINF:6.7333, no desc
channel-audio_1=96000-video=3442944-363987567.ts
## Auto Return Mode
#EXT-X-CUE-IN
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T12:45:24.566666Z
#EXTINF:2.8666, no desc
channel-audio_1=96000-video=3442944-363987568.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363987569.ts
```

## 2. Partial Ad Breaks

There are cases where the client requests the media playlist, and the manifest will have *partial* Ad Breaks, either at the beginning (segments outside playlist DVR) or at the end (segments not generated yet). These require special handling when parsing the media playlist information. We cannot, for example, calculate the Break's *Start Media Sequence* if we don't have the Break's first segment present on the manifest.

### 2.1. End of Playlist: Ad Break Is Not Ready Yet

It is possible that when the client requests the playlist, there will be an Ad Break **start** at the end of the manifest, but no segments for the Ad Break will have been generated yet. 

In the example below, the `EXT-X-DATERANGE` (`SCTE-OUT`) tag is already on the manifest (as might be the `EXT-X-CUE-OUT` tag) but we don't have the next media segment yet.

**Test File:** `/mocks/media/scte35/withAdBreakNewNotReady.m3u8`
```
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-MEDIA-SEQUENCE:363969827
#EXT-X-INDEPENDENT-SEGMENTS
#EXT-X-TARGETDURATION:7
#EXT-X-PROGRAM-DATE-TIME:2025-05-12T13:06:05.433333Z
(...)
channel-audio_1=96000-video=3442944-363969991.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363969992.ts
#EXTINF:6.2333, no desc
channel-audio_1=96000-video=3442944-363969993.ts
## splice_insert()
#EXT-X-DATERANGE:ID="1-1747055968",START-DATE="2025-05-12T13:19:30.466666Z",PLANNED-DURATION=60.033333,SCTE35-OUT=0xFC3025000000000BB802FFF01405000000017FEFFF8D788E687E00527178000100000000A4C46C9A
```

We need to confirm if the Ad Break is ready or not - that is, if the next media segment will be inside the Ad Break.

To do this, we estimate the next media segment's PDT (playlist DVR + playlist PDT) and check if it equals the `EXT-X-DATERANGE`'s `START-DATE`.

- If yes, the next media segment is inside the Ad Break. The Break's *Start Media Sequence* is the segment's media sequence and *Status* is `"complete"`.
  - **Test File:** `/mocks/media/scte35/withAdBreakNewReady.m3u8`
- If no, we *cannot* assume the next media segment is inside the Break. The Break's *Start Media Sequence* is `"0"` and *Status* is `"segmentsNotReady"`.
  - **Test File:** `/mocks/media/scte35/withAdBreakNewNotReady.m3u8`

Later, when the first segment for the Ad Break has been generated, we will have the `EXT-X-DATERANGE` (`SCTE-OUT`), `EXT-X-CUEOUT` and `EXT-X-PROGRAM-DATE-TIME` tags, followed by the first Break segment `EXTINF`. As normally, the Break's *Start Media Sequence* is the newest segment's media sequence and *Status* is `"complete"`.

**Test File:** `/mocks/media/scte35/withAdBreakNewReadyWithSegment.m3u8`
```
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-MEDIA-SEQUENCE:363969828
#EXT-X-INDEPENDENT-SEGMENTS
#EXT-X-TARGETDURATION:7
#EXT-X-PROGRAM-DATE-TIME:2025-05-12T13:06:10.233333Z
(...)
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363969991.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363969992.ts
#EXTINF:6.2333, no desc
channel-audio_1=96000-video=3442944-363969993.ts
## splice_insert()
#EXT-X-DATERANGE:ID="1-1747055968",START-DATE="2025-05-12T13:19:28.466666Z",PLANNED-DURATION=60.033333,SCTE35-OUT=0xFC3025000000000BB802FFF01405000000017FEFFF8D788E687E00527178000100000000A4C46C9A
#EXT-X-CUE-OUT:60.033333
#EXT-X-PROGRAM-DATE-TIME:2025-05-12T13:19:28.466666Z
#EXTINF:3.3666, no desc
channel-audio_1=96000-video=3442944-363969994.ts
```

### 2.2. Start of Playlist: Ad Break Leaves the DVR Limit

It is possible that when the client requests the playlist, there is an Ad Break that has already partially left the DVR range or is about to.

In the example below, the Ad Break is just about to leave the DVR, with only two non-Break segments remaining.

**Test File:** `/mocks/media/scte35/withAdBreakBeforeDVRLimit.m3u8`
```
#EXTM3U
#EXT-X-VERSION:11
## Created with Unified Streaming Platform  (version=1.14.4-30793)
#EXT-X-MEDIA-SEQUENCE:363991004
#EXT-X-INDEPENDENT-SEGMENTS
#EXT-X-TARGETDURATION:7
#USP-X-TIMESTAMP-MAP:MPEGTS=5361582439,LOCAL=2025-05-13T17:20:15.433300Z
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T17:20:15.433333Z
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363991004.ts
#EXTINF:6.4, no desc
channel-audio_1=96000-video=3442944-363991005.ts
## splice_insert(auto_return)
#EXT-X-DATERANGE:ID="4026559336-1747156826",START-DATE="2025-05-13T17:20:26.633333Z",PLANNED-DURATION=20,SCTE35-OUT=0xFC3025000000000BB800FFF01405F0006B687FEFFE90174E80FE001B774000010101000021F71DA8
#EXT-X-CUE-OUT:20
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T17:20:26.633333Z
#EXTINF:3.2, no desc
channel-audio_1=96000-video=3442944-363991006.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363991007.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363991008.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363991009.ts
#EXTINF:2.4, no desc
channel-audio_1=96000-video=3442944-363991010.ts
## Auto Return Mode
#EXT-X-CUE-IN
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T17:20:46.633333Z
#EXTINF:7.2, no desc
channel-audio_1=96000-video=3442944-363991011.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363991012.ts
```

We have two ways of confirming if the Break is about to leave the DVR limit.

1) If the playlist's PDT tag was already parsed, we check if the playlist PDT is equal or higher than the `EXT-X-DATERANGE`'s `START-DATE`.
2) If the playlist's PDT tag was not parsed yet, we check if there aren't any media segments remaining before the `EXT-X-DATERANGE` tag.

If either of these are true, then we set the Break's *Start Media Sequence* as `"0"` and *Status* as `"leavingDVRLimit"`.

**Test File:** `/mocks/media/scte35/withAdBreakOutsideDVRLimit.m3u8`
```
#EXTM3U
#EXT-X-VERSION:11
## Created with Unified Streaming Platform  (version=1.14.4-30793)
#EXT-X-MEDIA-SEQUENCE:363991008
#EXT-X-INDEPENDENT-SEGMENTS
#EXT-X-TARGETDURATION:7
#USP-X-TIMESTAMP-MAP:MPEGTS=5363310439,LOCAL=2025-05-13T17:20:34.633300Z
## splice_insert(auto_return)
#EXT-X-DATERANGE:ID="4026559336-1747156826",START-DATE="2025-05-13T17:20:26.633333Z",PLANNED-DURATION=20,SCTE35-OUT=0xFC3025000000000BB800FFF01405F0006B687FEFFE90174E80FE001B774000010101000021F71DA8
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T17:20:34.633333Z
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363991008.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363991009.ts
#EXTINF:2.4, no desc
channel-audio_1=96000-video=3442944-363991010.ts
## Auto Return Mode
#EXT-X-CUE-IN
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T17:20:46.633333Z
#EXTINF:7.2, no desc
channel-audio_1=96000-video=3442944-363991011.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363991012.ts
```

### 3. Case Study: Step-By-Step of Ad Break Leaving the DVR Limit

This section explains, step-by-step the changes in the manifest when the Ad Break is about to leave the DVR.

#### 3.1. Media Sequence: Current Segment is the FIRST media segment INSIDE the Break

On this moment, the manifest is on the limit of the DVR - that is, all non-Break segments have already left and the current media sequence is the Break's first segment.

**BEFORE:**
```
#EXTM3U
#EXT-X-VERSION:11
#EXT-X-MEDIA-SEQUENCE:363992684
#EXT-X-INDEPENDENT-SEGMENTS
#EXT-X-TARGETDURATION:7
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T19:34:39.433299Z
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363992684.ts
#EXTINF:5.3667, no desc
channel-audio_1=96000-video=3442944-363992685.ts
## splice_insert(auto_return)
#EXT-X-DATERANGE:ID="4026559475-1747164889",START-DATE="2025-05-13T19:34:49.599999Z",PLANNED-DURATION=20,SCTE35-OUT=0xFC3025000000000BB800FFF01405F0006BF37FEFFEBB581B38FE001B7740000101010000B80E326E
#EXT-X-CUE-OUT:20
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T19:34:49.599999Z
#EXTINF:4.2333, no desc
channel-audio_1=96000-video=3442944-363992686.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363992687.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363992688.ts
#EXTINF:6.1666, no desc
channel-audio_1=96000-video=3442944-363992689.ts
## Auto Return Mode
#EXT-X-CUE-IN
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T19:35:09.599999Z
#EXTINF:3.4333, no desc
channel-audio_1=96000-video=3442944-363992690.ts
#EXTINF:4.8, no desc
```

**NOW:**
```
#EXTM3U
#EXT-X-VERSION:11
#EXT-X-MEDIA-SEQUENCE:363992686
#EXT-X-INDEPENDENT-SEGMENTS
#EXT-X-TARGETDURATION:7
## splice_insert(auto_return)
#EXT-X-DATERANGE:ID="4026559475-1747164889",START-DATE="2025-05-13T19:34:49.599999Z",PLANNED-DURATION=20,SCTE35-OUT=0xFC3025000000000BB800FFF01405F0006BF37FEFFEBB581B38FE001B7740000101010000B80E326E
#EXT-X-CUE-OUT:20
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T19:34:49.599999Z
#EXTINF:4.2333, no desc
channel-audio_1=96000-video=3442944-363992686.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363992687.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363992688.ts
#EXTINF:6.1666, no desc
channel-audio_1=96000-video=3442944-363992689.ts
## Auto Return Mode
#EXT-X-CUE-IN
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T19:35:09.599999Z
#EXTINF:3.4333, no desc
channel-audio_1=96000-video=3442944-363992690.ts
#EXTINF:4.8, no desc
```

To avoid duplicate PDT tags, the playlist's first `EXT-X-PROGRAM-DATE` tag, which was tracking the media sequence, LEAVES the manifest, and the Break **start** PDT tag will take over accompaning the next media segments.

#### 3.2. Media Sequence: FIRST media segment INSIDE the Break has LEFT the Playlist

The `EXT-X-CUE-OUT` tag LEAVES the manifest alongside the first Break segment (i.e. the playlist media sequence is, at least, the SECOND media segment INSIDE the Break).

The `EXT-X-DATERANGE` (`SCTE-OUT`) tag will STAY during the Break and LEAVE only when the Break ends.

```
#EXTM3U
#EXT-X-VERSION:11
#EXT-X-MEDIA-SEQUENCE:363992688
#EXT-X-INDEPENDENT-SEGMENTS
#EXT-X-TARGETDURATION:7
## splice_insert(auto_return)
#EXT-X-DATERANGE:ID="4026559475-1747164889",START-DATE="2025-05-13T19:34:49.599999Z",PLANNED-DURATION=20,SCTE35-OUT=0xFC3025000000000BB800FFF01405F0006BF37FEFFEBB581B38FE001B7740000101010000B80E326E
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T19:34:58.633333Z
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363992688.ts
#EXTINF:6.1666, no desc
channel-audio_1=96000-video=3442944-363992689.ts
## Auto Return Mode
#EXT-X-CUE-IN
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T19:35:09.599999Z
#EXTINF:3.4333, no desc
channel-audio_1=96000-video=3442944-363992690.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363992691.ts
#EXTINF:4.8, no desc
```

#### 3.3. Media Sequence: Current segment is the FIRST media segment OUTSIDE the Break

To avoid duplicate PDT tags, the Break start PDT `EXT-X-PROGRAM-DATE` tag, which was tracking the media sequence, LEAVES the manifest, and the Break **end** PDT tag will take over accompaning the next media segments.

If there is a `EXT-X-DATERANGE` tag with `SCTE-IN`, the `EXT-X-DATERANGE` (`SCTE-OUT`) tag LEAVES the manifest. Otherwise, it leaves in the next media segment.

```
#EXTM3U
#EXT-X-VERSION:11
#EXT-X-MEDIA-SEQUENCE:363992690
#EXT-X-INDEPENDENT-SEGMENTS
#EXT-X-TARGETDURATION:7
## splice_insert(auto_return)
#EXT-X-DATERANGE:ID="4026559475-1747164889",START-DATE="2025-05-13T19:34:49.599999Z",PLANNED-DURATION=20,SCTE35-OUT=0xFC3025000000000BB800FFF01405F0006BF37FEFFEBB581B38FE001B7740000101010000B80E326E
## Auto Return Mode
#EXT-X-CUE-IN
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T19:35:09.599999Z
#EXTINF:3.4333, no desc
channel-audio_1=96000-video=3442944-363992690.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363992691.ts
#EXTINF:4.8, no desc
```

#### 3.4. Media Sequence: Current segment is the SECOND media segment OUTSIDE the Break

The `EXT-X-DATERANGE` (`SCTE-OUT`) and `EXT-X-CUE-IN` tags LEAVE the manifest.

```
#EXTM3U
#EXT-X-VERSION:11
## Created with Unified Streaming Platform  (version=1.14.4-30793) [0db1ee867823a80c252dd4604dd94aa9d482a83e6ae968df1da76cc43321ff2d]
#EXT-X-MEDIA-SEQUENCE:363992691
#EXT-X-INDEPENDENT-SEGMENTS
#EXT-X-TARGETDURATION:7
#USP-X-TIMESTAMP-MAP:MPEGTS=6090366439,LOCAL=2025-05-13T19:35:13.033300Z
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T19:35:13.033333Z
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363992691.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363992692.ts
```

Now, we await until the next Ad Break leave the manifest, when this process will repeat.