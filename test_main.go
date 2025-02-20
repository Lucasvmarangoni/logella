package main

import (
	"errors"
	"fmt"
	"net/http"

	logger "github.com/Lucasvmarangoni/logella/config/log"
	"github.com/Lucasvmarangoni/logella/router"
	test_errs "github.com/Lucasvmarangoni/logella/teste"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

func init() {
	logger.ConfigDefault()
}

func test_handler(w http.ResponseWriter, r *http.Request) {
}

// func main() {
// 	_, err := test_textError()
// 	log.Error().Err(test_errs.Test_Trace(err).Test_Stack()).Msg(fmt.Sprint(test_errs.Test_Trace(err).Code)) // Log
// }

func test_textRouter() {
	Router := router.NewRouter()

	r := chi.NewRouter()
	Router.Method(router.POST).Prefix("prefix").InitializeRoute(r, "/tanto", test_handler)
}

func test_textError() (string, error) {
	t := &test{}
	_, err := t.test_service()
	// test_errs.Test_Trace(err).Test_Trace()
	return "", test_errs.Test_Trace(err)
}

type test struct {
}

func (t *test) test_service() (string, error) {
	err := test_repo()
	return "", test_errs.Test_Trace(err)
}

func test_repo() error {
	return test_errs.Test_Wrap(errors.New("test error"), 500)
}
