package fractal

import (
	"fmt"
	"net"
	"os"
)

var (
	envPod    = os.Getenv("HOSTNAME")
	envPodIP  = findIP(envPod)
	envPodSvc = os.Getenv("GENFRACT_PORT")
)

func findIP(h string) string {
	//look up addr
	s, err := net.LookupHost(h)
	if err != nil {
		return fmt.Sprintf("err:%v", err)
	}
	r := ""
	for _, h := range s {
		r = fmt.Sprintf("%s:%s", r, h)
	}
	return r
}
