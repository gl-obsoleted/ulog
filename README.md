# ulog

## ulog build server steps

1. Install environments below
  - Sublime 3
  - Sublime Package Control
  - GoSublime (inside Package Control)

2. Setup GOPATH
  - Open `GoSubline.subline-settings` (from menu __Package Settings | GoSublime | Settings - User__) 
  - Use the sample below and __change the actual GOPATH to your local path__  
    ``` json
    {
    "env": { "GOPATH": "D:/<local_ulog_dir>/ulog_server/" },    
    }
    ```

3. Build Go server 
  - open `/ulog_server/src/ulog_sv/main.go` 
  - Select `Tools | Build Systems | GoSublime`
  - run `Ctrl - B` to open the build console
  - run `go get gopkg.in/mgo.v2`
  - run `go install`

4. Build Go test client
  - open `ulog_server/src/ulog_test/main.go`
  - run `go install`
  
5. Run the program
  - run server db `ulog_server_db/run_mongo_test_db.bat`
  - run server exe `ulog_server/bin/ulog_sv.exe`
  - run test client `ulog_server/bin/ulog_test.exe`


