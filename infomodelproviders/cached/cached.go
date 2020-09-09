package cached

import (
	"fmt"
	"github.com/iddqdeika/pim"
	"infomodel-service/definitions"
	"log"
	"sync"
	"time"
)

const (
	cacheLifetimeDuration = time.Minute * 5
)

func NewInfomodelProvider(p definitions.InfomodelProvider) (definitions.InfomodelProvider, error) {
	if p == nil {
		return nil, fmt.Errorf("given infomodel provider must not be nil")
	}
	return &cachedInfomodelProvider{
		p:       p,
		records: make(map[string]*cachedRecord),
	}, nil
}

type cachedInfomodelProvider struct {
	p definitions.InfomodelProvider

	records map[string]*cachedRecord
	m       sync.RWMutex
}

func (c *cachedInfomodelProvider) GetByIdentifier(identifier string) (*pim.StructureGroup, error) {
	cr, ok := c.records[identifier]
	if !ok {
		err := c.createNewCachedRecord(identifier)
		if err != nil {
			return nil, err
		}
	}
	cr, ok = c.records[identifier]
	if !ok {
		return nil, fmt.Errorf("internal error: cant find cache record")
	}
	cr.m.RLock()
	outDated := cr.outDated()
	cr.m.RUnlock()
	if outDated {
		err := c.updateCacheRecordIfOutdated(cr, identifier)
		if err != nil {
			return nil, err
		}
	}
	cr.m.RLock()
	defer cr.m.RUnlock()
	return cr.im, nil

}

func (c *cachedInfomodelProvider) createNewCachedRecord(identifier string) error {
	cr := &cachedRecord{}

	err := c.updateCacheRecord(cr, identifier)
	if err != nil {
		return err
	}

	c.m.Lock()
	defer c.m.Unlock()
	c.records[identifier] = cr
	return nil
}

func (c *cachedInfomodelProvider) updateCacheRecord(cr *cachedRecord, identifier string) error {
	im, err := c.p.GetByIdentifier(identifier)
	if err != nil {
		return fmt.Errorf("cant get infomodel from PIM")
	}
	cr.update(im)
	return nil
}

func (c *cachedInfomodelProvider) updateCacheRecordIfOutdated(cr *cachedRecord, identifier string) error {
	cr.m.Lock()
	defer cr.m.Unlock()
	if !cr.outDated() {
		return nil
	}
	err := c.updateCacheRecord(cr, identifier)
	if err != nil {
		log.Println(err)
	}
	return nil
}

type cachedRecord struct {
	im       *pim.StructureGroup
	deadline time.Time
	m        sync.RWMutex
}

func (c *cachedRecord) update(im *pim.StructureGroup) {
	c.im = im
	c.deadline = time.Now().Add(cacheLifetimeDuration)
}

func (c *cachedRecord) outDated() bool {
	if c.deadline.Before(time.Now()) {
		return true
	}
	return false
}
