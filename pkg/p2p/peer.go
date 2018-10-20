package p2p

import (
	"errors"
	"net"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	hasPortRegexp = regexp.MustCompile(":([1-9]|[1-8][0-9]|9[0-9]|[1-8][0-9]{2}|9[0-8][0-9]|99[0-9]|[1-8][0-9]{3}|9[0-8][0-9]{2}|99[0-8][0-9]|999[0-9]|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$")
)

type blacklisting struct {
	when   time.Time
	due    time.Time
	reason int
}

return u, err
}
