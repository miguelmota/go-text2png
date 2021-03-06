package main

/*
 * Sorry for bad code,
 * I'm still a beginner at golang.
 */

import (
  "fmt"
  "os"
  "image"
  "image/draw"
  "image/png"
  "bytes"
  "io/ioutil"
  "log"
  "math"
  "net/http"
  "html/template"
  "strconv"

  "github.com/golang/freetype/truetype"
  "golang.org/x/image/font"
  "golang.org/x/image/math/fixed"
)

func main() {
  // Cache html templates
  var templates = template.Must(template.ParseFiles("./index.html"))

  // HTTP Request handler for index
  http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")

    err := templates.ExecuteTemplate(w, "index.html", nil)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }
  })

  // HTTP Request handler for generating image
  http.HandleFunc("/image", func (w http.ResponseWriter, r *http.Request) {
    // Read query params
    text := []string{r.URL.Query().Get("text")}

    widthQ := r.URL.Query().Get("width")
    if widthQ == "" {
      widthQ = "200"
    }
    width64, _ := strconv.ParseInt(widthQ, 10, 64)
    width := int(width64)

    heightQ := r.URL.Query().Get("height")
    if heightQ == "" {
      heightQ = "20"
    }
    height64, _ := strconv.ParseInt(heightQ, 10, 64)
    height := int(height64)

    sizeQ := r.URL.Query().Get("size")
    if sizeQ == "" {
      sizeQ = "12"
    }
    size, err := strconv.ParseFloat(sizeQ, 64)
    if err != nil {
      size = 12.0
    }

    dpiQ := r.URL.Query().Get("dpi")
    if dpiQ == "" {
      dpiQ = "72"
    }
    dpi, err := strconv.ParseFloat(dpiQ, 64)
    if err != nil {
      dpi = 72.0
    }

    lineheightQ := r.URL.Query().Get("lineheight")
    if lineheightQ == "" {
      lineheightQ = "1.2"
    }
    lineheight, err := strconv.ParseFloat(lineheightQ, 64)
    if err != nil {
      lineheight = 1.2
    }

    fg := image.Black
    fgQ := r.URL.Query().Get("fg")
    if fgQ == "" {
      fgQ = "black"
    }
    if fgQ == "black" {
      fg = image.Black
    } else if fgQ == "white" {
      fg = image.White
    }

    bg := image.White
    bgQ := r.URL.Query().Get("bg")
    if bgQ == "" {
      bgQ = "white"
    }
    if bgQ == "black" {
      bg = image.Black
    }
    if bgQ == "none" || bgQ == "transparent" {
      bg = image.Transparent
    }

    hinting := r.URL.Query().Get("hinting")
    if hinting == "" {
      hinting = "none"
    }

    fontfile := "./luxisr.ttf"

    // Read the font data
    fontBytes, err := ioutil.ReadFile(fontfile)
    if err != nil {
      log.Println(err)
      return
    }
    f, err := truetype.Parse(fontBytes)
    if err != nil {
      log.Println(err)
      return
    }

    // Draw the background
    rgba := image.NewRGBA(image.Rect(0, 0, width, height))
    draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

    // Draw the text
    h := font.HintingNone
    switch hinting {
    case "full":
      h = font.HintingFull
    }
    d := &font.Drawer{
      Dst: rgba,
      Src: fg,
      Face: truetype.NewFace(f, &truetype.Options{
        Size:    size,
        DPI:     dpi,
        Hinting: h,
      }),
    }
    x := 0
    y := 0
    dy := int(math.Ceil(size * lineheight * dpi / 72))
    d.Dot = fixed.Point26_6{
      X: fixed.I(x),
      Y: fixed.I(y),
    }
    y += dy
    for _, s := range text {
      d.Dot = fixed.P(0, y)
      d.DrawString(s)
      y += dy
    }

    buffer := new(bytes.Buffer)
    if err := png.Encode(buffer, rgba); err != nil {
      log.Println("unable to encode image.")
    }

    // Write response
    w.Header().Set("Content-Type", "image/png")
    w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
    if _, err := w.Write(buffer.Bytes()); err != nil {
      log.Println("unable to write image.")
    }
  })

  port := os.Getenv("PORT")
  if port == "" {
    port = "9361"
  }

  fmt.Println("Listening on " + port)

  // Listen on port
  host := ":" + port
  http.ListenAndServe(host, nil)
}
