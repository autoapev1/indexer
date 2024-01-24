package syncer

import (
	"context"
	"errors"
	"log/slog"

	"github.com/autoapev1/indexer/config"
	"github.com/autoapev1/indexer/eth"
	"github.com/autoapev1/indexer/storage"
)

var (
	ErrNoNetwork = errors.New("no network provided")
	ErrNoStore   = errors.New("no store provided")
)

type Syncer struct {
	config  config.Config
	network *eth.Network
	store   storage.Store
	ctx     context.Context
}

func NewSyncer(conf config.Config) *Syncer {
	return &Syncer{
		config: conf,
	}
}

func (s *Syncer) WithNetwork(n *eth.Network) *Syncer {
	s.network = n
	return s
}

func (s *Syncer) WithStore(st storage.Store) *Syncer {
	s.store = st
	return s
}

func (s *Syncer) WithContext(ctx context.Context) *Syncer {
	s.ctx = ctx
	return s
}

func (s *Syncer) Init() error {
	if s.network == nil {
		return ErrNoNetwork
	}

	if s.store == nil {
		return ErrNoStore
	}

	if s.ctx == nil {
		s.ctx = context.Background()
	}

	if !s.network.Ready() {
		slog.Warn("network not ready, initializing")
		err := s.network.Init()
		if err != nil {
			return err
		}
	}

	if !s.store.Ready() {
		slog.Warn("store not ready, initializing")
		err := s.store.Init()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Syncer) Sync(ctx context.Context) error {
	if err := s.Init(); err != nil {
		return err
	}

	if err := s.ArchiveSync(ctx); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return nil
	default:
	}

	return nil
}

func (s *Syncer) ArchiveSync(ctx context.Context) error {

	chainHeight, err := s.network.Web3.Eth.GetBlockNumber()
	if err != nil {
		return err
	}

	heights, err := s.store.GetHeights()
	if err != nil {
		slog.Error("failed to get db heights")
		return err
	}

	if (int64(chainHeight) - heights.Blocks) > 0 {
		slog.Info("chain height is higher than db blocktimestamp height, syncing block timestamps", "chainHeight", chainHeight, "dbHeight", heights.Blocks)
		bts, err := s.network.GetBlockTimestamps(ctx, heights.Blocks, int64(chainHeight))
		if err != nil {
			slog.Error("failed to get block timestamps", "error", err)
			return err
		}

		err = s.store.BulkInsertBlockTimestamp(bts)
		if err != nil {
			slog.Error("failed to ingest block timestamps", "error", err)
			return err
		}
	}

	if (int64(chainHeight) - heights.Pairs) > 0 {
		slog.Info("chain height is higher than db pair height, syncing pairs", "chainHeight", chainHeight, "dbHeight", heights.Pairs)
		pairs, err := s.network.GetPairs(ctx, heights.Pairs, int64(chainHeight))
		if err != nil {
			slog.Error("failed to get pairs", "error", err)
			return err
		}

		err = s.store.BulkInsertPairInfo(pairs)
		if err != nil {
			slog.Error("failed to ingest pairs", "error", err)
			return err
		}

		toFetchTokens, err := s.store.GetPairsWithoutTokenInfo()
		if err != nil {
			slog.Error("failed to get pairs without token info", "error", err)
			return err
		}

		tokens, err := s.network.GetTokenInfo(ctx, toFetchTokens)
		if err != nil {
			slog.Error("failed to get tokens", "error", err)
			return err
		}

		err = s.store.BulkInsertTokenInfo(tokens)
		if err != nil {
			slog.Error("failed to ingest tokens", "error", err)
			return err
		}
	}
	return nil
}
func (s *Syncer) LiveSync(ctx context.Context, bn int64) {}
func (s *Syncer) BlockOracle(ctx context.Context)        {}
