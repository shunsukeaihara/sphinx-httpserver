package sphinx

import (
	"sync"

	ps "github.com/shunsukeaihara/go-pocketsphinx"
	"golang.org/x/net/context"
)

type key int

const psKey key = 0

type Sphinx map[string][]*PsInstance

type PsInstance struct {
	*ps.PocketSphinx
	mu   sync.Mutex
	lang string
}

func NewSphinx(cfgMap map[string]ps.Config, cpunum int) Sphinx {
	ret := Sphinx{}
	for lang, config := range cfgMap {
		sli := make([]*PsInstance, 0, cpunum)
		for i := 0; i < cpunum; i++ {
			sphinx := ps.NewPocketSphinx(config)
			sli = append(sli, &PsInstance{sphinx, sync.Mutex{}, lang})
		}
		ret[lang] = sli
	}

	return ret
}

func NewContext(ctx context.Context, sp Sphinx) context.Context {
	return context.WithValue(ctx, psKey, sp)
}

func FromContext(ctx context.Context) (Sphinx, bool) {
	ps, ok := ctx.Value(psKey).(Sphinx)
	return ps, ok
}

func (t Sphinx) GetSphinxFromLanguage(lang string) (*PsInstance, error) {

	return nil, nil
}
