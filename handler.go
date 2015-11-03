package main

import (
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/Sirupsen/logrus"
	ps "github.com/shunsukeaihara/go-pocketsphinx"
	"github.com/unrolled/render"
	"golang.org/x/net/context"

	"github.com/shunsukeaihara/sphinx-httpserver/sphinx"
	"github.com/shunsukeaihara/sphinx-httpserver/utils"
)

var renderer = render.New(render.Options{})

type Rsponse struct {
	Response ps.Result `json:"response"`
}

func getSphinx(ctx context.Context) (*sphinx.PsInstance, string, string, error) {
	lang := sphinx.LangFromContext(ctx)
	sp, ok := sphinx.FromContext(ctx)
	if !ok {
		log.Errorln("speech recognition engine is not ready")
		return nil, lang, "speech recognition engine is not ready", errors.New("speech recognition engine is not ready")
	}
	ps, err := sp.GetSphinxFromLanguage(lang)
	if err != nil {
		log.Errorln(err)
		return nil, lang, "invalid accept-language", err

	}
	return ps, lang, "", nil
}

func dictationHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	ps, lang, errmsg, err := getSphinx(ctx)
	if err != nil {
		utils.WriteJsonErrorResponse(w, errmsg, 500)
		return
	}
	log.WithFields(log.Fields{
		"lang": lang,
	}).Info()

	log.Infoln(lang)
	ps.Lock()
	defer ps.Unlock()

	buf := make([]byte, 1024)
	ps.StartUtt()
	for {
		size, err := r.Body.Read(buf)
		if err != nil {
			break
		}
		ps.ProcessRaw(buf[:size], false, false)
	}
	ps.EndUtt()
	res, err := ps.GetHyp(false)
	if err != nil {
		utils.WriteJsonErrorResponse(w, "recognition error", 500)
		return
	}
	bytes, err := json.Marshal(Rsponse{res})
	if err != nil {
		utils.WriteJsonErrorResponse(w, "error", 500)
	}
	renderer.JSON(w, http.StatusOK, bytes)
	return
}
