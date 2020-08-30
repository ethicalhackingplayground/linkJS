package main

import (
"fmt"
"time"
"sync"
"regexp"
"os"
"bufio"
"net/http"
"io/ioutil"
"flag"
)

func main () {

	LINK_REGEX:=`(https?|ftp|file)://[-A-Za-z0-9\+&@#/%?=~_|!:,.;]*[-A-Za-z0-9\+&@#/%=~_|]`

	var concurrency int	
	flag.IntVar(&concurrency, "c", 30, "Set concurrency for greater speed")
	flag.Parse()

	var wg sync.WaitGroup
	for i:=0; i<=concurrency; i++ {
		wg.Add(1)
		go func () {
			get_links(LINK_REGEX)
			wg.Done()	
		}()
		wg.Wait()
	}
}

func get_links(regex string) {

	time.Sleep(time.Millisecond * 10)
	client:= &http.Client{}
	

	scanner:=bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		jsLink:=scanner.Text()
		
		req,err:=http.NewRequest("GET", jsLink,nil)
		if err != nil {
			return
		}

		resp,err:=client.Do(req)
		if err != nil {
			return
		}

		bodyBuffer,err:=ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		bodyString:=string(bodyBuffer)

		// Regex Magic
		re:=regexp.MustCompile(regex)
		match:=re.FindStringSubmatch(bodyString)
		if match != nil {
			fmt.Printf("%q\n", match[0])
		}
	}
}











