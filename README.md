# Monacoind-proxy

Monacoind proxy is a rcp reverse proxy to reproduce timeout issue from Monacoind side

## Configure

- `delay`: Defines the reponse delay to reproduce timeout issue
- `endpoint`: Defines the monacoind api to forward to

## Example how to use JSON-RPC server written in Go with JSON-RPC client written in python

```python
import json
import requests

def rpc_call(url, method, args):
    headers = {'content-type': 'application/json'}
    payload = {
        "method": method,
        "params": [args],
        "jsonrpc": "2.0",
        "id": 1,
    }
    response = requests.post(url, data=json.dumps(payload), headers=headers).json()
    return response['result']

url = 'http://localhost:1337/proxy'

sendfromArgs = {'From': 'monappy-28357','To': 'pBGTgV7zaCpiXEGZERwcVuJib6Asbtwdfc', 'Amount': '1.000'}

print rpc_call(url, 'proxy.SendFrom', sendfromArgs)

```

## Contact

<m_ju@indiesquare.me>