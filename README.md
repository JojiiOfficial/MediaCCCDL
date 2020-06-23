# MediaCCCDl
(mirrored from [gitlab](https://git.jojii.de/jojii/mediacccdl))
A tool to retrieve download URLs for [media.ccc.de](media.ccc.de) videos.

# Usage
`mediacccdl <url> [--format] [--lang]`

### Supported formats:
<b>Video:</b>
- mp4
- webm
<b>Audio:</b>
- mp3
- opus

### Supported languages
- deu/de
- eng/en

# Exit codes
0 - no error <br>
1 - flag error <br>
2 - HTTP error <br>
3 - No URL found <br>

# Examples
Download a video/audio
```bash
wget $(mediacccdl "<URL>" --format mp4/mp3)
```

Download in german/english
```bash
wget $(mediacccdl "<URL>" --lang de/en)
```

Download a list of videos parallel
```bash
cat ccc | parallel wget '$(mediacccdl {})'
```
