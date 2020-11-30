package root

import (
	"context"
	"github.com/iddqdeika/infomodel-service/config"
	"github.com/iddqdeika/infomodel-service/definitions"
	"github.com/iddqdeika/infomodel-service/infomodelproviders/cached"
	"github.com/iddqdeika/infomodel-service/infomodelproviders/pim"
	"github.com/iddqdeika/infomodel-service/web"
	"github.com/iddqdeika/rrr"
)

func New() rrr.Root {
	return &root{}
}

type root struct {
	ws definitions.WebService
}

func (r *root) Register() []error {
	errs := make([]error, 0)
	e := func(err error) {
		if err != nil {
			errs = append(errs, err)
		}
	}

	cfg, err := config.NewJsonCfg("cfg.json")
	e(err)
	if err != nil {
		return errs
	}

	//pim infomodel provider
	ip, err := pim.NewInfomodelProvider(cfg.Child("pim"))
	e(err)

	//caching wrapper
	cip, err := cached.NewInfomodelProvider(ip)
	e(err)

	ws, err := web.NewService(cfg.Child("web"), cip)
	e(err)
	r.ws = ws

	if len(errs) == 0 {
		return nil
	}
	return errs
}

func (r *root) Resolve(ctx context.Context) error {
	return r.ws.Run()
}

func (r *root) Release() error {
	return nil
}
