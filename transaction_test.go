package echo

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type debugTransaction struct {
	*bytes.Buffer
}

func (b *debugTransaction) Begin(ctx context.Context) error {
	_, err := b.WriteString(`[Begin]`)
	return err
}

func (b *debugTransaction) Rollback(ctx context.Context) error {
	_, err := b.WriteString(`[Rollback]`)
	return err
}

func (b *debugTransaction) Commit(ctx context.Context) error {
	_, err := b.WriteString(`[Commit]`)
	return err
}

func (b *debugTransaction) End(ctx context.Context, succeed bool) error {
	if succeed {
		return b.Commit(ctx)
	}
	return b.Rollback(ctx)
}

func TestTransaction(t *testing.T) {
	b := bytes.NewBuffer(nil)
	tx := &debugTransaction{b}
	w := NewTransaction(tx)
	f := func(rollback bool) error {
		b.WriteString(`<f:Begin`)
		w.Begin(nil)
		b.WriteString(`>`)
		var e error
		if rollback {
			b.WriteString(`<f:Rollback`)
			w.Rollback(nil)
			e = ErrNotFound
		} else {
			b.WriteString(`<f:Commit`)
			w.Commit(nil)
		}
		b.WriteString(`>`)
		return e
	}
	f2 := func(rollback bool, rollback2 bool) error {
		b.WriteString(`<f2:Begin`)
		w.Begin(nil)
		b.WriteString(`>`)
		e := f(rollback2)
		if e != nil {
			return e
		}
		if rollback {
			b.WriteString(`<f2:Rollback`)
			w.Rollback(nil)
			e = ErrNotFound
		} else {
			b.WriteString(`<f2:Commit`)
			w.Commit(nil)
		}
		b.WriteString(`>`)
		return e
	}

	// -------------
	// - Commit
	// -------------

	b.WriteString(`<Begin`)
	w.Begin(nil)
	b.WriteString(`>`)
	err := f2(false, false)
	if err == nil {
		b.WriteString(`<Commit`)
	} else {
		b.WriteString(`<Rollback`)
	}
	w.End(nil, err == nil)
	b.WriteString(`>`)
	assert.Equal(t, `<Begin[Begin]><f2:Begin><f:Begin><f:Commit><f2:Commit><Commit[Commit]>`, b.String())

	// -------------
	// - Rollback
	// -------------

	b.Reset()

	b.WriteString(`<Begin`)
	w.Begin(nil)
	b.WriteString(`>`)
	err = f2(false, true)
	if err == nil {
		b.WriteString(`<Commit`)
	} else {
		b.WriteString(`<Rollback`)
	}
	w.End(nil, err == nil)
	b.WriteString(`>`)
	assert.Equal(t, `<Begin[Begin]><f2:Begin><f:Begin><f:Rollback[Rollback]><Rollback>`, b.String())
}
