package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
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
	Id          string            `json:"id"`
	Scripts     map[string]script `json:"scripts"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
}

type script struct {
	Text     string `json:"text"`
	Id       string `json:"id"`
	Src      string `json:"src"`
	Speaker  string `json:"speaker"`
	Duration int    `json:"duration"`
	IsLast   bool   `json:"isLast"`
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

	defer os.Remove(fileName + ".wav")

	// Reads the audio file into memory.
	data, err := ioutil.ReadFile(fileName + ".wav")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	ctx := c.Request().Context()
	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:     speechpb.RecognitionConfig_LINEAR16,
			LanguageCode: "en-US",
			Model:        "phone_call",
			UseEnhanced:  true,
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

	s := story{
		Id:          id,
		Name:        "Apresentation Intro",
		Description: "You are talking with Laura, She's from New Zeland and want to met new people.",
		Scripts: map[string]script{
			"0": script{
				Id:      "0",
				Text:    "It wasn’t just that I was 41, which, let’s face it, isn’t old. It was that I was 41 and bored. And a little tired. And, at times, cantankerous. Crotchety, you might say.",
				Src:     "http://192.168.0.7:1323/audio/text_laura_1.mp3",
				Speaker: "speaker-laura",
			},
			"1": script{
				Id:       "1",
				Text:     "Exacerbating this problem was the fact that I had spent the entire span of my thirties at one place — a prestigious men’s magazine. I thought I had stability and security and swagger.",
				Src:      "http://192.168.0.7:1323/audio/text_you_1.mp3",
				Speaker:  "speaker-you",
				IsLast:   true,
				Duration: 5,
			},
			// "2": script{
			// 	Id:      "2",
			// 	Text:    "What I didn’t realize is that I had slowly started draining energy from the place where I worked instead of injecting it with my own. I was getting soft. I was getting lazy.",
			// 	Src:     "http://192.168.0.7:1323/audio/text_laura_2.mp3",
			// 	Speaker: "speaker-laura",
			// },
			// "3": script{
			// 	Id:       "3",
			// 	Text:     "A couple months into unemployment, I got a job at another prestigious men’s magazine.",
			// 	Src:      "http://192.168.0.7:1323/audio/text_you_2.mp3",
			// 	Speaker:  "speaker-you",
			// 	IsLast:   true,
			// 	Duration: 3,
			// },
		},
	}
	return c.JSON(http.StatusOK, s)
}
