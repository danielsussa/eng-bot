package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
	Duration int16  `json:"duration"`
	IsLast   bool   `json:"ssLast"`
}

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		Skipper:      middleware.DefaultSkipper,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.GET("/story/:id", getStory)
	e.Static("/audio", "api-server/public")
	e.Logger.Fatal(e.Start(":1323"))
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
				Src:     "http://localhost:1323/audio/text_laura_1.mp3",
				Speaker: "speaker-laura",
			},
			"1": script{
				Id:      "1",
				Text:    "Exacerbating this problem was the fact that I had spent the entire span of my thirties at one place — a prestigious men’s magazine. I thought I had stability and security and swagger.",
				Src:     "http://localhost:1323/audio/text_you_1.mp3",
				Speaker: "speaker-you",
			},
			"2": script{
				Id:      "2",
				Text:    "What I didn’t realize is that I had slowly started draining energy from the place where I worked instead of injecting it with my own. I was getting soft. I was getting lazy.",
				Src:     "http://localhost:1323/audio/text_laura_2.mp3",
				Speaker: "speaker-laura",
			},
			"3": script{
				Id:      "3",
				Text:    "A couple months into unemployment, I got a job at another prestigious men’s magazine.",
				Src:     "http://localhost:1323/audio/text_you_2.mp3",
				Speaker: "speaker-you",
			},
		},
	}
	return c.JSON(http.StatusOK, s)
}