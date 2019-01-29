# BV Test
Back End Developer Test

Information and requirements

- Developed in Golang v1.9.7
- install gcc if you are using windows -> http://tdm-gcc.tdragon.net/download
- dependencies file in glide.yaml

- #### api will run on localhost:8005


###### `Installation steps`

1 - clone the project from https://github.com/jimmivc/bvTest

2 - go to the project path `go/src/{path}/bvTest/`
 
3 - run `glide i` to install all dependencies

4 - run `go run spotycloud.go` to run the api or `go build` to generate executable file

If You want to run unit tests, just run `go test` inside the project path

###### `API Routes`
- /api/songs/artist/{artistName}
- /api/songs/{songName}
- /api/songs/genre/{genreName}
- /api/genres/summary
- /api/songs/byLength/{min}/{max}
