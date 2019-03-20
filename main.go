package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
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
	Transaction string
}

// Service service wrapper
type Service struct {
	Config *config.Config
}

// SendFrom passthru to monacoind with a response delay
func (s *Service) SendFrom(r *http.Request, args *SendFromArgs, result *Response) (err error) {
	rawreq := fmt.Sprintf(`{"jsonrpc":"1.0","id":"curltext","method":"sendfrom","params":["%s","%s", %f]}`, args.From, args.To, args.Amount)
	url := s.Config.Endpoint

	reqbody := strings.NewReader(rawreq)
	req, err := http.NewRequest("POST", url, reqbody)
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

	// todo: json unmarshall it check response structure
	*result = Response{Transaction: string(respbody)}
	time.Sleep(s.Config.GetDelay())
	return
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)

	go func() {
		s := <-sigs
		log.Printf("RECEIVED SIGNAL: %s", s)
		os.Exit(1)
	}()

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
