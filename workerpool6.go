package main

import("fmt"; "os"; "bufio"; "net/http"; "log"; "time"; "sync")

var allData []dataURL

type dataURL struct{
  url string
}

func fetchStatus(w http.ResponseWriter, r *http.Request){

  start := time.Now()

  var wg sync.WaitGroup

  jobsChannel := make(chan string, 23)

  for w:=1; w<20; w++{
    wg.Add(1)
    go worker(jobsChannel, &wg)
  }

  for j,_:= range allData{
    jobsChannel <- allData[j].url
  }
  close(jobsChannel)

  wg.Wait()
  elapsed := time.Since(start)
  fmt.Println("TOOK ==================================>> %s\n", elapsed)
}


func worker(jobsChannel <- chan string, wg *sync.WaitGroup){

  defer wg.Done()

  for data := range jobsChannel{
    resp, err := http.Get(data)
    if err != nil{
      fmt.Println(err)
      return
    }
    time.Sleep(1000*time.Millisecond)
    fmt.Println(string(data), "-----> ", resp.Status)
  }
}


func readFile(path string)(result []string){
    //read the lines of urlList.txt file
    file, err := os.Open(path)
    if err != nil {
        fmt.Println(err)
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }

    return lines
}


func main(){

  lines := readFile("urlList.txt")

  for i:=1; i<len(lines); i++{
    data := dataURL{url: lines[i]}
    allData = append(allData, data)
  }


  http.HandleFunc("/", fetchStatus)
  log.Fatal(http.ListenAndServe(":8080", nil))
}
