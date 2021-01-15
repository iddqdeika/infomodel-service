package web

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"infomodel-service/definitions"

	"github.com/iddqdeika/pim"
)

const (
	getInfomodelByIdentifierPath = "/getInfomodelByIdentifier"
	identifierQueryParam         = "identifier"
)

func NewService(config definitions.Config, ip definitions.InfomodelProvider) (definitions.WebService, error) {
	if config == nil {
		return nil, fmt.Errorf("config must not be nil")
	}
	if ip == nil {
		return nil, fmt.Errorf("infomodel provider must not be nil")
	}
	port, err := config.GetInt("port")
	if err != nil {
		return nil, err
	}

	return &webService{
		cfg: Config{Port: port},
		ip:  ip,
	}, nil
}

type Config struct {
	Port int
}

type webService struct {
	cfg Config
	ip  definitions.InfomodelProvider
}

func (ws *webService) Run() error {
	log.Println("http handler init")
	//addr, err := nethelper.GetCurrentAddr(ws.cfg.Port)
	//if err != nil {
	//	return err
	//}
	http.HandleFunc(getInfomodelByIdentifierPath, ws.getInfomodelByIdentifier)

	log.Println("http handler started on port: " + strconv.Itoa(ws.cfg.Port))
	return http.ListenAndServe(":"+strconv.Itoa(ws.cfg.Port), nil)
}

func (ws *webService) getInfomodelByIdentifier(w http.ResponseWriter, req *http.Request) {
	ids, ok := req.URL.Query()[identifierQueryParam]
	if !ok {
		w.Write([]byte("identifier query param must be set"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(ids) != 1 {
		w.Write([]byte("only one identifier query param must be set"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	identifier := ids[0]
	g, err := ws.ip.GetByIdentifier(identifier)
	if err != nil {
		w.Write([]byte("internal server error: " + err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch req.Header.Get("Accept") {
	case "application/xml":
		im := xmlConvertInfomodel(g)
		data, err := xml.Marshal(im)
		if err != nil {
			w.Write([]byte("internal server error: " + err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	default:
		im := jsonConvertInfomodel(g)
		data, err := json.Marshal(im)
		if err != nil {
			w.Write([]byte("internal server error: " + err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
}

func jsonConvertInfomodel(g *pim.StructureGroup) *definitions.JsonInfomodelDTO {
	fs := make(map[string]*definitions.FeatureDTO)
	for _, feature := range g.Features {
		f := &definitions.FeatureDTO{
			Name:         feature.Name,
			DataType:     feature.DataType,
			PresetValues: feature.PresetValues,
			Mandatory:    feature.Mandatory,
			Multivalued:  feature.Multivalued,
		}
		fs[feature.Name] = f
	}
	im := &definitions.JsonInfomodelDTO{
		Identifier: g.Identifier,
		Features:   fs,
	}
	return im
}

func xmlConvertInfomodel(g *pim.StructureGroup) *definitions.XmlInfomodelDTO {
	var arr []definitions.FeatureDTO
	for _, feature := range g.Features {
		f := definitions.FeatureDTO{
			Name:         feature.Name,
			DataType:     feature.DataType,
			PresetValues: feature.PresetValues,
			Mandatory:    feature.Mandatory,
			Multivalued:  feature.Multivalued,
		}
		arr = append(arr, f)
	}
	im := &definitions.XmlInfomodelDTO{
		Identifier: g.Identifier,
		Features:   arr,
	}
	return im
}
