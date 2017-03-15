package main

import (
  //"bufio"
  "flag"
  //"fmt"
  "image"
  "image/draw"
  "image/png"
  //"image/jpeg"
  "bytes"
  "io/ioutil"
  "log"
  "math"
  "net/http"
  "strconv"
  //"os"

  "github.com/golang/freetype/truetype"
  "golang.org/x/image/font"
  "golang.org/x/image/math/fixed"
)

var (
  dpi      = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
  fontfile = flag.String("fontfile", "./luxisr.ttf", "filename of the ttf font")
  hinting  = flag.String("hinting", "none", "none | full")
  size     = flag.Float64("size", 12, "font size in points")
  spacing  = flag.Float64("spacing", 1.2, "line spacing (e.g. 2 means double spaced)")
)

func main() {
  // HTTP Request handler
  http.HandleFunc("/image", func (w http.ResponseWriter, r *http.Request) {
    flag.Parse()

    text := []string{r.URL.Query().Get("text")}

    // Read the font data.
    fontBytes, err := ioutil.ReadFile(*fontfile)
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
    fg, bg := image.White, image.Black
    const imgW, imgH = 160, 20
    rgba := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
    draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

    // Draw the text
    h := font.HintingNone
    switch *hinting {
    case "full":
      h = font.HintingFull
    }
    d := &font.Drawer{
      Dst: rgba,
      Src: fg,
      Face: truetype.NewFace(f, &truetype.Options{
        Size:    *size,
        DPI:     *dpi,
        Hinting: h,
      }),
    }
    x := 0
    y := 0
    dy := int(math.Ceil(*size * *spacing * *dpi / 72))
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

  host := ":8080"
  http.ListenAndServe(host, nil)
  /*

  // Save that RGBA image to disk.
  outFile, err := os.Create("out.png")
  if err != nil {
    log.Println(err)
    os.Exit(1)
  }
  defer outFile.Close()
  b := bufio.NewWriter(outFile)
  err = png.Encode(b, rgba)
  if err != nil {
    log.Println(err)
    os.Exit(1)
  }
  err = b.Flush()
  if err != nil {
    log.Println(err)
    os.Exit(1)
  }
  fmt.Println("Wrote out.png OK.")
  */
}
