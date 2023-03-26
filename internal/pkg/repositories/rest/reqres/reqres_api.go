package reqres

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Fajar-Islami/go-simple-user-crud/internal/pkg/repositories/rest"
)

type ReqResAPI interface {
	GetListUser(ctx context.Context, params ReqListUser) (res ResListUsers, err error)
}
type ReqResAPIImpl struct {
	URI  string
	Opts rest.Opts
}

func NewReqResAPI(URI string, Opts rest.Opts) ReqResAPI {
	return &ReqResAPIImpl{
		URI:  URI,
		Opts: Opts,
	}
}

func (lur *ReqResAPIImpl) GetListUser(ctx context.Context, params ReqListUser) (res ResListUsers, err error) {
	uri := fmt.Sprintf("%s?page=%d&per_page=%d", lur.URI, params.Page, params.PerPage)

	respHttp, err := rest.DoRequest(ctx, http.MethodGet, uri, lur.Opts)
	if err != nil {
		log.Println(err)
		return res, err
	}

	defer respHttp.Body.Close()
	err = json.NewDecoder(respHttp.Body).Decode(&res)
	if err != nil {
		log.Println(err)
		return res, err
	}

	return res, nil
}
