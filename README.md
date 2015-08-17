# ulog

## ulog build server steps

1. Install environments below
  - Sublime 3
  - Sublime Package Control
  - GoSublime (inside Package Control)

2. Setup GOPATH
  - Open 'GoSubline.subline-settings' (from menu __Package Settings | GoSublime | Settings - User__) 
  - Use the sample below and __change the actual GOPATH to your local path__ 
    ``` json
    {
    "env": { "GOPATH": "D:/<local_ulog_dir>/ulog_server/" },    
    }
    ```

3. open `/ulog_server/src/ulog_sv/main.go` 
  - run `go get gopkg.in/mgo.v2`
  - run `go install`
