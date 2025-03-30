package proxy

import "net/http"

func (pd *ProxyDelivery) deleteHeaders(req *http.Request) {
	req.Header.Del("Proxy-Connection")
}
