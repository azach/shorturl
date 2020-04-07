package pool

import (
	"github.com/azach/shorturl/lib/cache"
	"github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
	"math"
	"sync"
)

const minPoolSize = 3
const minPoolGenerationSize = 1
const idAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Pool struct {
	mux   sync.Mutex
	queue []string
	cache cache.Cache
}

func NewPool(cache cache.Cache) *Pool {
	return &Pool{
		queue: []string{},
		cache: cache,
	}
}

func (p *Pool) Get() string {
	p.mux.Lock()
	defer p.mux.Unlock()

	candidate := p.queue[0]
	p.queue = p.queue[1:]
	return candidate
}

func (p *Pool) Generate() {
	if len(p.queue) < minPoolSize {
		numToGenerate := int(math.Min(float64(minPoolGenerationSize), float64(minPoolSize-len(p.queue))))

		for i := 0; i < numToGenerate; i++ {
			candidate, err := shortid.Generate()
			if err != nil {
				logrus.Errorf("error generating candidate: %s", err)
				continue
			}

			logrus.Infof("created candidate: %s %v", candidate, len(p.queue))

			_, exists := p.cache.Get(candidate)
			if exists {
				logrus.Infof("candidate already exists %s", candidate)
				continue
			}

			p.queue = append(p.queue, candidate)
		}
	}
}
