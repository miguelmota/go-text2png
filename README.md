# text2png

> golang server to convert text to PNG image.

Uses [go-freetype](https://github.com/golang/freetype/).

## Instructions

```bash
go build main.go
./main
```

Query parameter options with defaults

```
text=""
width=200
height=20
lineheight=1.2
bg=white (white|black|transparent)
fg=black (white|black)
dpi=72
hinting=none (none|full)
```

## Example

[http://localhost:9361/image?text=Hello+World&size=24&fg=black&bg=white&width=125&height=40](http://localhost:9361/image?text=Hello+World&size=24&fg=black&bg=white&width=125&height=40)

<img src="./example.png" height="40" />

# License

MIT
