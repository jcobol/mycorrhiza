package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/internal/shroom"
	"github.com/bouncepaw/mycorrhiza/internal/user"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/web/viewutil"
)

type countingWriter struct {
	*httptest.ResponseRecorder
	headers int
}

func (cw *countingWriter) WriteHeader(status int) {
	cw.headers++
	cw.ResponseRecorder.WriteHeader(status)
}

// oldHandlerUploadBinary simulates the behaviour before the fix where processing
// continued after an early error.
func oldHandlerUploadBinary(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	rq.ParseMultipartForm(10 << 20)
	var (
		hyphaName = util.HyphaNameFromRq(rq, "upload-binary")
		h         = hyphae.ByName(hyphaName)
		u         = user.FromRequest(rq)
		lc        = l18n.FromRequest(rq)
		_, _, err = rq.FormFile("binary")
		meta      = viewutil.MetaFrom(w, rq)
	)
	if err != nil {
		viewutil.HttpErr(meta, http.StatusInternalServerError, hyphaName, err.Error())
	}
	if err := shroom.CanAttach(u, h, lc); err != nil {
		viewutil.HttpErr(meta, http.StatusInternalServerError, hyphaName, err.Error())
	}
	if err != nil {
		return
	}
}

func TestHandlerUploadBinaryStopsAfterError(t *testing.T) {
	origAuth := cfg.UseAuth
	cfg.UseAuth = true
	t.Cleanup(func() { cfg.UseAuth = origAuth })

	viewutil.Init()

	rr1 := &countingWriter{ResponseRecorder: httptest.NewRecorder()}
	req := httptest.NewRequest(http.MethodPost, "/upload-binary/foo", nil)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=foo")

	// Old behaviour produced multiple header writes.
	oldHandlerUploadBinary(rr1, req)
	if rr1.headers < 2 {
		t.Fatalf("expected old handler to write header twice, got %d", rr1.headers)
	}

	rr2 := &countingWriter{ResponseRecorder: httptest.NewRecorder()}
	handlerUploadBinary(rr2, req)
	if rr2.headers != 1 {
		t.Fatalf("expected single header write, got %d", rr2.headers)
	}
}
