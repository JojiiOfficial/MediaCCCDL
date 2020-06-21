package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/JojiiOfficial/gaw"
	"github.com/PuerkitoBio/goquery"
)

var supportedFormats = []string{"mp4", "webm", "mp3", "opus"}
var supportedLangs = []string{"en", "de", "deu", "eng"}

func main() {
	if gaw.IsInStringArray("--help", os.Args[1:]) || gaw.IsInStringArray("-h", os.Args[1:]) {
		fmt.Println("Usage:\ncccdl [url] [--format] [--lang]")
		return
	}

	var format, lang string
	flag.StringVar(&format, "format", "mp4", "The format of the video/audio to download")
	flag.StringVar(&lang, "lang", "auto", "The language of the video/audio to download")

	// First param has to be the url
	flag.CommandLine.Parse(os.Args[2:])

	// Validate format
	if !gaw.IsInStringArray(format, supportedFormats) {
		fmt.Printf("Format '%s' not supported\n", format)
		os.Exit(1)
	}

	// Validate language
	if lang != "auto" && !gaw.IsInStringArray(lang, supportedLangs) {
		fmt.Printf("Language '%s' not supported\n", format)
		os.Exit(1)
	}

	if lang != "auto" {
		lang = formatLang(lang)
	}

	// Get link
	blogTitles, err := GetDownloadURL(os.Args[1], format, lang)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
	if len(blogTitles) == 0 {
		os.Exit(3)
		return
	}

	fmt.Printf(blogTitles)
}

// GetDownloadURL gets the dl url
func GetDownloadURL(url, format, lang string) (string, error) {
	// Get the HTML
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	// Convert HTML into goquery document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// Remember format
	isVideo := true

	typeSection := "video"
	if isAudioFormat(format) {
		typeSection = "audio"
		isVideo = false
	}

	// Get Video buttons
	dlSection := doc.Find(".downloads").Find("." + typeSection)

	// For video use correct ID
	if isVideo {
		dlSection = dlSection.Find("#" + strings.ToLower(format))
	}

	var link string
	// Find a content
	dlSection.Find("a").EachWithBreak(func(i int, sel *goquery.Selection) bool {
		href, has := sel.Attr("href")
		if !has {
			// Continue
			return true
		}

		// On VideoMode stop here
		if isVideo {
			link = href
			return false
		}

		// ----------------
		// Audio only
		// ----------------

		// Find correct language & audio version
		langSub := sel.Find(".language")
		if langSub == nil || len(langSub.Text()) == 0 {
			return true
		}

		// Lang does not match
		if lang != "auto" && strings.ToLower(langSub.Text()) != lang {
			return true
		}

		titleSub := sel.Find(".title")
		if titleSub == nil || len(titleSub.Text()) == 0 {
			return true
		}

		sTitle := strings.Split(titleSub.Text(), " ")
		if len(sTitle) == 1 {
			return true
		}

		version := strings.ToLower(sTitle[1])
		if format == version {
			link = href
			return false
		}

		return true
	})

	return link, nil
}

// Check if format is an audio format
func isAudioFormat(format string) bool {
	return gaw.IsInStringArray(format, []string{"mp3", "opus"})
}

// Return the correct language class name
func formatLang(lang string) string {
	if len(lang) == 0 {
		return ""
	}

	lang = strings.ToLower(lang)

	switch lang {
	case "de":
		return "deu"
	case "en":
		return "eng"
	}

	return lang
}

// GetLocale from system
// TODO maybe implement this to automatically select the language
func GetLocale() (string, error) {
	envlang, ok := os.LookupEnv("LANG")
	if ok {
		return strings.Split(envlang, ".")[0], nil
	}

	// Exec powershell Get-Culture on Windows.
	cmd := exec.Command("powershell", "Get-Culture | select -exp Name")
	output, err := cmd.Output()
	if err == nil {
		return strings.Trim(string(output), "\r\n"), nil
	}

	return "", fmt.Errorf("cannot determine locale")
}
