package telegram

import (
	"context"
	"fmt"
	"sync"

	"github.com/9seconds/mtg/conntypes"
	"github.com/9seconds/mtg/telegram/api"
	"github.com/9seconds/mtg/wrappers"
)

type middleTelegram struct {
	baseTelegram

	secret []byte
	mutex  sync.RWMutex
}

func (m *middleTelegram) update() error {
	secret, err := api.Secret()
	if err != nil {
		return fmt.Errorf("cannot fetch secret: %w", err)
	}

	v4Addresses, v4DefaultDC, err := api.AddressesV4()
	if err != nil {
		return fmt.Errorf("cannot fetch addresses for ipv4: %w", err)
	}

	v6Addresses, v6DefaultDC, err := api.AddressesV6()
	if err != nil {
		return fmt.Errorf("cannot fetch addresses for ipv6: %w", err)
	}

	m.mutex.Lock()
	m.secret = secret
	m.v4DefaultDC = v4DefaultDC
	m.V6DefaultDC = v6DefaultDC
	m.v4Addresses = v4Addresses
	m.v6Addresses = v6Addresses
	m.mutex.Unlock()

	return nil
}

func (m *middleTelegram) Dial(ctx context.Context,
	cancel context.CancelFunc,
	dc conntypes.DC,
	protocol conntypes.ConnectionProtocol) (wrappers.StreamReadWriteCloser, error) {
	if dc == 0 {
		dc = conntypes.DCDefaultIdx
	}

	return m.baseTelegram.dial(ctx, cancel, dc, protocol)
}
