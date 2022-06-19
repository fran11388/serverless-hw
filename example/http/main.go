package main
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)
func main(){
	streamname:="my-event-stream"
	restapiId:="px3bcgybfh"
	url:=fmt.Sprintf("https://%s.execute-api.ap-northeast-1.amazonaws.com/test/streams/%s/record",restapiId,streamname)
	payload:=strings.NewReader("{\n  \"Data\": {\n    \"client_id\": \"client831209\",\n    \"msg\": \"send from postman 2\",\n    \"issue_error\": false\n  },\n  \"PartitionKey\": \"1\"\n}")
	req, err := http.NewRequest("PUT", url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type","application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

}
