package main

import (
  "fmt"
  "os"
  "net/http"
  "math/rand"
  "strings"
  "time"

  "github.com/labstack/echo/v4"
  "github.com/rs/zerolog"
  "github.com/rs/zerolog/log"
)

var (
  db map[string]string
  charset string
  rand_len int
)

func init() {
  db = make(map[string]string)
  // z-base-32
  charset = "yeabt3nk4dmu5rcwhfp7gqs68x9"
  rand_len = 5
}

func random_string() string {
  rand.Seed(time.Now().UnixNano())

  sb := strings.Builder{}
  sb.Grow(rand_len)
  for i := 0; i < rand_len; i++ {
    sb.WriteByte(charset[rand.Intn(len(charset))])
  }

  return sb.String()
}


func main() {
    e := echo.New()

    e.POST("/short", func(c echo.Context) error {
      url := c.FormValue("url")

      if _, ok := db[url]; !ok {
        s := random_string()
        db[s] = url
        db[url] = s
      }

      log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
      log.Info().Str("url", url).Str("code", db[url]).Msg("")

      return c.String(http.StatusOK, db[url])
    })

    e.GET("/:id", func(c echo.Context) error {
      if _, ok := db[c.Param("id")]; !ok {
        return c.String(http.StatusInternalServerError, fmt.Sprintf("Could not find %s", c.Param("id")))
      }

      return c.Redirect(http.StatusMovedPermanently, db[c.Param("id")])
    })

    e.Logger.Fatal(e.Start(":3002"))
}
