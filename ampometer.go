package main

import (
    "io/ioutil"
    "log"
    "strconv"
    "errors"
    "bufio"
    "strings"
    "os"
    "time"
    )

func findpx() (int, error) {
    pids, err := ioutil.ReadDir("/proc")
    if err != nil {
        log.Fatal(err)
    }

    for _, pid := range pids {
        if pid.IsDir() && pid.Name()[0] >= '1' && pid.Name()[0] <= '9' {

            cmd, err := ioutil.ReadFile("/proc/" + pid.Name() + "/comm")
            if err == nil {
                if string(cmd) == "px-storage\n" {
                    p, _ := strconv.Atoi(pid.Name())
                    return p, nil
                }
            }
        }
    }

    return -1, errors.New("unable to find px-storage process")
}


func iostats(pid int) (int64, int64) {
    var wc int64
    var wb int64
    ios, err := os.Open("/proc/" + strconv.Itoa(pid) + "/io")
    if err == nil {
        defer ios.Close()
        scanner := bufio.NewScanner(ios)
        scanner.Split(bufio.ScanLines)
        for scanner.Scan() {
            f := strings.Split(scanner.Text(), ":")
            if f[0] == "wchar" {
                wc, _ = strconv.ParseInt(strings.TrimSpace(f[1]), 10, 64)
            }
            if f[0] == "write_bytes" {
                wb, _ = strconv.ParseInt(strings.TrimSpace(f[1]), 10, 64)
            }
        }
    }
    
    return wc, wb
}

func amprate(pid int) {
    var prev_c int64
    var curr_c int64
    var prev_b int64
    var curr_b int64
    var rate_c int64
    var rate_b int64
    
    for {
        prev_c = curr_c
        prev_b = curr_b
        
        curr_c, curr_b = iostats(pid)
        
        rate_c = curr_c - prev_c
        rate_b = curr_b - prev_b
        
        //fmt.Println("prev_c:", prev_c, "curr_c:", curr_c, "rate_c:", rate_c, "prev_b:", prev_b, "curr_b:", curr_b, "rate_b:", rate_b)
        log.Printf("Amplification Factor: %d   [Requested: %.1f MB  Written: %.1f MB]", rate_b / rate_c, float64(rate_c)/1024.0/1024.0, float64(rate_b)/1024.0/1024.0)
        time.Sleep(1 * time.Second)
    }
    
}    

func main() {
    
    p, err := findpx()

    if err != nil {
        log.Fatal(err)
    }

    amprate(p)

}
