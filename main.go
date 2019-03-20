package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/ilovelili/monacoind-proxy/config"
)

// SendFromArgs sendfrom arguments
type SendFromArgs struct {
	From   string // From account sendfrom
	To     string // To address sendto
	Amount float64
}

// Response response
type Response struct {
	Result []byte
}

// Service service wrapper
type Service struct {
	Config *config.Config
}

// SendFrom passthru to monacoind with a response delay
func (s *Service) SendFrom(r *http.Request, args *SendFromArgs, result *Response) (err error) {
	reqbody := strings.NewReader(fmt.Sprintf(`{"jsonrpc":"2.0","id":"curltext","method":"sendfrom","params":["%s","%s", %f]`, args.From, args.To, args.Amount))
	req, err := http.NewRequest("POST", s.Config.Endpoint, reqbody)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "text/plain;")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	time.Sleep(s.Config.GetDelay())
	*result = Response{Result: respbody}
	return
}

func main() {
	proxy := &Service{
		Config: config.GetConfig(),
	}

	rpcServer := rpc.NewServer()
	rpcServer.RegisterCodec(json.NewCodec(), "application/json")
	rpcServer.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	rpcServer.RegisterService(proxy, "proxy")

	router := mux.NewRouter()
	router.Handle("/proxy", rpcServer)
	fmt.Println("RPC proxy running on :1337")
	http.ListenAndServe(":1337", router)
}
