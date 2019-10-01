package stringid

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

var pushRE = regexp.MustCompile(`^server/[A-Za-z0-9_\-]{20}$`)

func TestMiddleware(t *testing.T) {
	t.Parallel()

	s := httptest.NewServer(Middleware(WithPrefix("server/"))(readReq))
	defer s.Close()

	ids := make(map[string]bool)
	for i := 0; i < 10000; i++ {
		res, err := doReq(s)
		if err != nil {
			t.Fatalf("test %d expected no error, got: %v", i, err)
			continue
		}
		if !pushRE.MatchString(res) {
			t.Errorf("test %d received invalid result: %q", i, res)
			continue
		}
		if _, ok := ids[res]; ok {
			t.Errorf("test %d produced previously used ID: %q", i, res)
		}
	}
}

var uuidRE = regexp.MustCompile(`^server/[a-z0-9\-]{36}$`)

func TestMiddlewareUUID(t *testing.T) {
	t.Parallel()

	s := httptest.NewServer(Middleware(WithPrefix("server/"), WithGenerator(NewUUIDGenerator()))(readReq))
	defer s.Close()

	ids := make(map[string]bool)
	for i := 0; i < 10000; i++ {
		res, err := doReq(s)
		if err != nil {
			t.Fatalf("test %d expected no error, got: %v", i, err)
			continue
		}
		if !uuidRE.MatchString(res) {
			t.Errorf("test %d received invalid result: %q", i, res)
			continue
		}
		if _, ok := ids[res]; ok {
			t.Errorf("test %d produced previously used ID: %q", i, res)
		}
	}
}

func TestHeaderMiddleware(t *testing.T) {
	t.Parallel()

	s := httptest.NewServer(HeaderMiddleware("x-id")(readReq))
	defer s.Close()

	tests := []struct {
		headers []string
		exp     string
	}{
		{nil, ""},
		{[]string{"x", "id"}, ""},
		{[]string{"x-id", "foo"}, "foo"},
		{[]string{"x-foo", "foo", "x-id", "bar"}, "bar"},
	}

	for i, test := range tests {
		res, err := doReq(s, test.headers...)
		if err != nil {
			t.Errorf("test %d expected no error, got: %v", i, err)
			continue
		}
		if res != test.exp {
			t.Errorf("test %d expected %q, got: %q", i, test.exp, res)
		}
	}
}

var readReq = http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, FromRequest(req))
})

// doReq handles issuing a request to a url and returning the body of the
// response.
func doReq(s *httptest.Server, headers ...string) (string, error) {
	if len(headers)%2 != 0 {
		return "", errors.New("invalid headers length")
	}

	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		return "", err
	}

	for i := 0; i < len(headers); i += 2 {
		req.Header.Set(headers[i], headers[i+1])
	}

	cl := &http.Client{}
	res, err := cl.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", nil
	}

	return strings.TrimSpace(string(body)), nil
}
