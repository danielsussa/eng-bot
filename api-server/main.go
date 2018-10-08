package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"acesso.io/acessorh/lib/uuid"

	speech "cloud.google.com/go/speech/apiv1"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

type story struct {
	Id          string             `json:"id"`
	Scripts     map[string]*script `json:"scripts"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
}

type script struct {
	Text     string   `json:"text"`
	Id       string   `json:"id"`
	Grade    []string `json:"grade"`
	Src      string   `json:"src"`
	Speaker  string   `json:"speaker"`
	Duration float64  `json:"duration"`
	Interval string   `json:"interval"`
	IsLast   bool     `json:"isLast"`
}

var client *speech.Client

func main() {
	var err error
	// Creates a client.

	cfg, err := google.JWTConfigFromJSON([]byte(os.Getenv("SPEECH_CREDENTIAL")), speech.DefaultAuthScopes()...)
	if err != nil {
		panic(err)
	}

	ts := cfg.TokenSource(context.Background())
	opt := option.WithTokenSource(ts)

	client, err = speech.NewClient(context.Background(), opt)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	////HTTP/////////
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		Skipper:      middleware.DefaultSkipper,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{}))

	e.GET("/story/:id", getStory)

	//e.Static("/audio", "api-server/public")
	e.GET("/audio/:id", getAudio)

	e.POST("/story/audio/send/:id", sendAudio)

	e.Logger.Fatal(e.Start(":1323"))
}

func sendAudio(c echo.Context) error {
	// SAVE TO FILE
	fileName := uuid.New().String()
	file, err := c.FormFile("audio")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(fileName + ".mp3")
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	defer os.Remove(fileName + ".mp3")

	//CONVERT
	_, err = exec.Command("sh", "-c", fmt.Sprintf("ffmpeg -i %s.mp3 -sample_rate 24000 -y %s.wav", fileName, fileName)).Output()

	if err != nil {
		return err
	}

	//defer os.Remove(fileName + ".wav")

	// Reads the audio file into memory.
	data, err := ioutil.ReadFile(fileName + ".wav")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	ctx := c.Request().Context()
	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			LanguageCode:    "en-US",
			MaxAlternatives: 3,
			Model:           "phone_call",
			//UseEnhanced:  true,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	})

	if err != nil {
		log.Fatalf("failed to recognize: %v", err)
	}

	// Prints the results.
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}

	return c.String(200, "ok")
}

func getAudio(c echo.Context) error {
	time.Sleep(time.Millisecond * time.Duration(rand.Int31n(6000)))
	fmt.Println(c.Param("id"))
	return c.File(fmt.Sprintf("api-server/public/%s", c.Param("id")))
}

// get story
func getStory(c echo.Context) error {
	time.Sleep(2 * time.Second)
	// User ID from path `users/:id`
	id := c.Param("id")

	story := story{
		Id:          id,
		Name:        "Apresentation Intro",
		Description: "You are talking with Laura, She's from New Zeland and want to met new people.",
		Scripts: map[string]*script{
			"0": &script{
				Id:   "0",
				Text: "I want you to take this seriously! Phoebe is very very important to me, ok? And I wanna make sure that you are gonna take care of her.",
				Src:  "http://192.168.0.7:1323/audio/text_laura_1.mp3",

				Speaker: "speaker-laura",
			},
			"1": &script{
				Id:      "1",
				Text:    "Joe, I love Phoebe. She’s the single most important thing in my life. I’d die before I let anything happen to her.",
				Src:     "http://192.168.0.7:1323/audio/text_244.mp3",
				Speaker: "speaker-you",
				Grade:   []string{"full", "half", "none"},
				//IsLast:   true,
				Duration: 8,
			},
			"2": &script{
				Id:      "2",
				Text:    "That’s what I wanted to hear! Because she’s family, ok, and now you’re gonna be family, and there is nothing more important in the whole world, than family.",
				Src:     "http://192.168.0.7:1323/audio/text_laura_2.mp3",
				Speaker: "speaker-laura",
				//IsLast:  true,
			},
			"3": &script{
				Id:       "3",
				Text:     "It’s your favorite sister!",
				Src:      "http://192.168.0.7:1323/audio/text_you_2.mp3",
				Speaker:  "speaker-you",
				Grade:    []string{"full", "none", "none", "none", "none"},
				IsLast:   true,
				Duration: 3,
			},
		},
	}

	init := make(map[string]float64)

	for _, script := range story.Scripts {
		startInterval := init[script.Speaker]
		init[script.Speaker] += script.Duration
		endInterval := init[script.Speaker]

		script.Interval = fmt.Sprintf("#t=%3.2f,%3.2f", startInterval, endInterval)
	}

	return c.JSON(http.StatusOK, story)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////

const PAIR_ACTIVE = 1
const PAIR_INACTIVE = 2

// maximum elements in a filepath
const FILEPATH_MAX = 64

type Pair struct {
	fst     byte
	snd     byte
	status  byte
	__align byte
}

type Pairs []Pair

func (pair *Pair) String() string {
	return fmt.Sprintf("%c%c|%d ", pair.fst, pair.snd, pair.status)
}

func detectLength(tokens []string) int {
	result := 0
	for i := 0; i < len(tokens); i++ {
		l := len(tokens[i]) - 1
		if l > 0 {
			result += l
		}
	}
	return result
}

// This helper func is missing in Go 1.0 (String type)
func splitWithRegexp(s string, re *regexp.Regexp) []string {
	if len(re.String()) > 0 && len(s) == 0 {
		return []string{""}
	}
	matches := re.FindAllStringIndex(s, -1)
	strings := make([]string, 0, len(matches))
	beg := 0
	end := 0
	for _, match := range matches {
		end = match[0]
		if match[1] != 0 {
			strings = append(strings, s[beg:end])
		}
		beg = match[1]
	}
	if end != len(s) {
		strings = append(strings, s[beg:])
	}
	return strings
}

func splitFilepath(s string) []string {
	if len(s) == 0 {
		return []string{}
	}
	strs := make([]string, 0, FILEPATH_MAX)
	slc := s
	ix := 0
	for {
		ix = strings.IndexRune(slc, '/')
		if ix == -1 {
			break
		}
		strs = append(strs, slc[0:ix])
		slc = slc[ix+1:]
	}
	strs = append(strs, slc[ix+1:])
	return strs
}

func NewPairsFromArray(tokens []string) Pairs {
	dl := detectLength(tokens)
	pairs := make(Pairs, dl)

	k := 0
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		for j := 0; j < len(t)-1; j++ {
			pairs[k].fst = t[j]
			pairs[k].snd = t[j+1]
			pairs[k].status = PAIR_ACTIVE
			k++
		}
	}
	return pairs
}

func NewPairsFromString(str string) Pairs {
	return NewPairsFromArray([]string{str})
}

func NewPairsFromStringTokens(str string, re regexp.Regexp) Pairs {
	ss := splitWithRegexp(str, &re)
	return NewPairsFromArray(ss)
}

func NewPairsFromFilepath(str string) Pairs {
	ss := splitFilepath(str)
	return NewPairsFromArray(ss)
}

func (self Pairs) String() string {
	buff := bytes.NewBufferString("")
	for _, p := range self {
		buff.WriteString(p.String())
	}
	return buff.String()
}

func (self Pairs) Reactivate() {
	for _, p := range self {
		p.status = PAIR_ACTIVE
	}
}

func (a Pair) Equal(b Pair) bool {
	return (a.fst == b.fst && a.snd == b.snd &&
		(a.status&b.status&PAIR_ACTIVE) == PAIR_ACTIVE)
}

func (self Pairs) Match(other Pairs) float64 {
	matches := 0
	len_self := len(self)
	len_other := len(other)
	sum := len_self + len_other
	if sum == 0 {
		return 1.0
	}
	for i := 0; i < len_self; i++ {
		for j := 0; j < len_other; j++ {
			if self[i].Equal(other[j]) {
				matches++
				other[j].status = PAIR_INACTIVE
				break
			}
		}
	}
	return float64(2*matches) / float64(sum)
}

func MatchStrings(stra, strb string) float64 {
	a := NewPairsFromString(stra)
	b := NewPairsFromString(strb)
	return a.Match(b)
}

func MatchStringsTokens(stra, strb string, re *regexp.Regexp) float64 {
	a := NewPairsFromStringTokens(stra, *re)
	b := NewPairsFromStringTokens(strb, *re)
	return a.Match(b)
}
