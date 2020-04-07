package pool

import (
	"github.com/azach/shorturl/lib/cache"
	"github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
	"math"
	"sync"
)

// These can be tuned according to scalability/performance needs
const defaultMinPoolSize = 100
const defaultMinPoolGenerationSize = 10

type Options struct {
	minPoolSize           int
	minPoolGenerationSize int
}

type Pool struct {
	mux                   sync.Mutex
	queue                 []string
	cache                 cache.Cache
	minPoolSize           int
	minPoolGenerationSize int
}

func NewPool(cache cache.Cache, options *Options) *Pool {
	minPoolSize := defaultMinPoolSize
	if options.minPoolSize > 0 {
		minPoolSize = options.minPoolSize
	}

	minPoolGenerationSize := defaultMinPoolGenerationSize
	if options.minPoolGenerationSize > 0 {
		minPoolGenerationSize = options.minPoolGenerationSize
	}

	return &Pool{
		queue:                 []string{},
		cache:                 cache,
		minPoolSize:           minPoolSize,
		minPoolGenerationSize: minPoolGenerationSize,
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
	if len(p.queue) < p.minPoolSize {
		numToGenerate := int(math.Min(float64(p.minPoolGenerationSize), float64(p.minPoolSize-len(p.queue))))

		for i := 0; i < numToGenerate; i++ {
			candidate, err := shortid.Generate()
			if err != nil {
				logrus.Errorf("error generating candidate: %s", err)
				continue
			}

			_, exists := p.cache.Get(candidate)
			if exists {
				logrus.Infof("candidate already exists %s", candidate)
				continue
			}

			logrus.Infof("created candidate: %s %v", candidate, len(p.queue))

			p.queue = append(p.queue, candidate)
		}
	}
}
