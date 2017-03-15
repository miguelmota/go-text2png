golang text to png

WIP

```bash
go build main.go
./main
```

[http://localhost:9361/image?text=hello](http://localhost:9361/image?text=hello)

query params

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

Example

[http://localhost:9361/image?text=Hello+World&size=24&fg=black&bg=white&width=125&height=40](http://localhost:9361/image?text=Hello+World&size=24&fg=black&bg=white&width=125&height=40)

<img src="./example.png" width="300" />

# License

MIT
