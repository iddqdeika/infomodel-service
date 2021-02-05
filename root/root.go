package root

import (
	"context"

	"infomodel-service/definitions"
	"infomodel-service/infomodelproviders/cached"
	"infomodel-service/infomodelproviders/pim"
	"infomodel-service/web"

	"github.com/iddqdeika/reactivetools/misc"
	"github.com/iddqdeika/rrr"
	"github.com/iddqdeika/rrr/helpful"
)

func New() rrr.Root {
	return &root{}
}

type root struct {
	infomodelService definitions.WebService
	echoService      definitions.WebService
}

func (r *root) Register() []error {
	errs := make([]error, 0)
	e := func(err error) {
		if err != nil {
			errs = append(errs, err)
		}
	}

	cfg, err := helpful.NewJsonCfg("cfg.json")
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
	r.infomodelService = ws

	//echo service
	es, err := misc.NewEchoService(cfg.Child("echo_service"), helpful.DefaultLogger)
	e(err)
	r.echoService = es

	if len(errs) == 0 {
		return nil
	}
	return errs
}

func (r *root) Resolve(ctx context.Context) error {
	return rrr.ComposeErrors("Resolce", rrr.RunServices(ctx, r.infomodelService, r.echoService)...)
}

func (r *root) Release() error {
	return nil
}
