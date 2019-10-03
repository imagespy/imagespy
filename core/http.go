package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/imagespy/imagespy/discovery"
	log "github.com/sirupsen/logrus"
)

type StorageWriter interface {
	Write(*discovery.Input) error
}

type discoverHandler struct {
	s StorageWriter
}

func (d *discoverHandler) discover(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warnf("read discover payload: %w", err)
		writeError(w, http.StatusBadRequest, errors.New("unable to read payload"))
		return
	}
	defer r.Body.Close()

	in := &discovery.Input{}
	err = json.Unmarshal(b, in)
	if err != nil {
		log.Warnf("unmarshal discover payload: %w", err)
		writeError(w, http.StatusBadRequest, errors.New("unable to decode payload into JSON"))
		return
	}

	errs := discovery.ValidateInput(in)
	if len(errs) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		for _, e := range errs {
			fmt.Fprintln(w, e.Error())
		}

		return
	}

	err = d.s.Write(in)
	if err != nil {
		log.Errorf("write discover payload: %w", err)
		writeError(w, http.StatusInternalServerError, errors.New("unable to write payload to storage"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func writeError(w http.ResponseWriter, c int, err error) {
	w.WriteHeader(c)
	w.Write([]byte(err.Error()))
}
