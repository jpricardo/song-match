package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"song-match-backend/domain"
	"song-match-backend/internal/audioutil"
	"song-match-backend/internal/ytbutil"
	"sort"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrackUsecase struct {
	trackRepository domain.TrackRepository
	contextTimeout  time.Duration

	// shutdown can wait for in-flight processing to complete before exiting.
	wg sync.WaitGroup
	// shutdown is closed by Shutdown() to signal goroutines to stop early.
	shutdown chan struct{}
}

func NewTrackUsecase(trackRepository domain.TrackRepository, timeout time.Duration) *TrackUsecase {
	return &TrackUsecase{
		trackRepository: trackRepository,
		contextTimeout:  timeout,
		shutdown:        make(chan struct{}),
	}
}

// Shutdown waits for all background processing goroutines to finish.
// Call this during graceful server shutdown so no track is left permanently
// stuck in "processing" status.
func (tu *TrackUsecase) Shutdown() {
	close(tu.shutdown)
	tu.wg.Wait()
}

func (tu *TrackUsecase) FindMatches(c context.Context, content []byte) ([]domain.Track, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	samples, sampleRate, err := audioutil.DecodeAudio(content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode audio: %w", err)
	}

	fingerprints, err := audioutil.ExtractFingerprints(samples, sampleRate)
	if err != nil {
		return nil, fmt.Errorf("failed to extract fingerprints: %w", err)
	}

	sampleHashes := audioutil.GenerateHashes(fingerprints)

	var sampleHashValues []string
	sampleHashMap := make(map[string][]float64)

	for _, h := range sampleHashes {
		sampleHashValues = append(sampleHashValues, h.HashValue)
		sampleHashMap[h.HashValue] = append(sampleHashMap[h.HashValue], h.Time)
	}

	matchedDBHashes, err := tu.trackRepository.GetMatchingHashes(ctx, sampleHashValues)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch matching hashes: %w", err)
	}

	// Group the matched hashes by Track ID so we can build a histogram per track
	trackHistograms := make(map[primitive.ObjectID]map[int]int)

	for _, dbHash := range matchedDBHashes {
		if _, exists := trackHistograms[dbHash.TrackID]; !exists {
			trackHistograms[dbHash.TrackID] = make(map[int]int)
		}

		if sampleTimes, exists := sampleHashMap[dbHash.HashValue]; exists {
			for _, sampleTime := range sampleTimes {
				offsetBin := int((dbHash.Time - sampleTime) * 10)
				trackHistograms[dbHash.TrackID][offsetBin]++
			}
		}
	}

	// Score each track based on its histogram
	type trackScore struct {
		trackID primitive.ObjectID
		score   int
	}
	var scores []trackScore

	for trackID, offsets := range trackHistograms {
		maxScore := 0
		for _, count := range offsets {
			if count > maxScore {
				maxScore = count
			}
		}

		if maxScore > 5 {
			scores = append(scores, trackScore{trackID: trackID, score: maxScore})
		}
	}

	// Sort highest score first
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// Collect all matched IDs and fetch in a single query.
	matchedIDs := make([]string, 0, len(scores))
	scoreByID := make(map[string]int, len(scores))
	for _, s := range scores {
		hex := s.trackID.Hex()
		matchedIDs = append(matchedIDs, hex)
		scoreByID[hex] = s.score
	}

	matchedTracks, err := tu.trackRepository.GetManyByIDs(ctx, matchedIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch matched tracks: %w", err)
	}

	// Re-sort fetched tracks to match score order, since $in does not
	// guarantee ordering.
	sort.Slice(matchedTracks, func(i, j int) bool {
		return scoreByID[matchedTracks[i].ID.Hex()] > scoreByID[matchedTracks[j].ID.Hex()]
	})

	if len(matchedTracks) > 0 {
		slog.Info("match found",
			"best_match", matchedTracks[0].Name,
			"score", scoreByID[matchedTracks[0].ID.Hex()],
			"total_candidates", len(matchedTracks),
		)
	} else {
		slog.Info("no match found")
	}

	return matchedTracks, nil
}

func (tu *TrackUsecase) GetMany(c context.Context) ([]domain.Track, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	return tu.trackRepository.Fetch(ctx)
}

func (tu *TrackUsecase) GetByID(c context.Context, id string) (domain.Track, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	return tu.trackRepository.GetByID(ctx, id)
}

func (tu *TrackUsecase) AddTrack(c context.Context, url string) (*domain.Track, error) {
	// Use a separate background context for the DB insert so
	// the write is not canceled when AddTrack returns and the request context
	// is torn down. The caller's context is only used to respect any upstream
	// cancellation that arrives before we even start.
	select {
	case <-c.Done():
		return nil, c.Err()
	default:
	}

	insertCtx, insertCancel := context.WithTimeout(context.Background(), tu.contextTimeout)
	defer insertCancel()

	t := &domain.Track{
		Name:   "Processing...", // Temporary name
		Url:    url,
		Status: "processing",
	}

	err := tu.trackRepository.Create(insertCtx, t)
	if err != nil {
		return nil, fmt.Errorf("failed to create track record: %w", err)
	}

	tu.wg.Go(func() {
		tu.processTrackInBackground(t.ID, url)
	})

	return t, nil
}

func (tu *TrackUsecase) processTrackInBackground(trackID primitive.ObjectID, url string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	log := slog.With("track_id", trackID.Hex(), "url", url)
	log.Info("background processing started")

	// Respect shutdown signal — abort early and leave the
	// track in "processing" rather than completing a half-written state. The
	// server restart will need to clean up stale "processing" records, but at
	// least no corrupt data is written.
	select {
	case <-tu.shutdown:
		log.Warn("background processing aborted: server shutting down")
		tu.trackRepository.UpdateTrackStatus(ctx, trackID.Hex(), "failed")
		return
	default:
	}

	wav, title, thumbnail, err := ytbutil.DownloadTrack(url)
	if err != nil {
		log.Error("background processing failed at download", "error", err)
		tu.trackRepository.UpdateTrackStatus(ctx, trackID.Hex(), "failed")
		return
	}

	samples, sampleRate, err := audioutil.DecodeAudio(wav)
	if err != nil {
		log.Error("background processing failed at decode", "error", err)
		tu.trackRepository.UpdateTrackStatus(ctx, trackID.Hex(), "failed")
		return
	}

	fingerprints, err := audioutil.ExtractFingerprints(samples, sampleRate)
	if err != nil {
		log.Error("background processing failed at fingerprinting", "error", err)
		tu.trackRepository.UpdateTrackStatus(ctx, trackID.Hex(), "failed")
		return
	}

	hashes := audioutil.GenerateHashes(fingerprints)

	// Update the database with the real data and mark it as "ready"
	track := &domain.Track{
		ID:           trackID,
		Name:         title,
		Thumbnail:    thumbnail,
		Status:       "ready",
		Fingerprints: fingerprints,
		Hashes:       hashes,
	}

	err = tu.trackRepository.UpdateTrackData(ctx, track)
	if err != nil {
		log.Error("background processing failed at DB update", "error", err)
		tu.trackRepository.UpdateTrackStatus(ctx, trackID.Hex(), "failed")
		return
	}

	log.Info("background processing complete", "title", title)
}

func (tu *TrackUsecase) DeleteByID(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	return tu.trackRepository.DeleteByID(ctx, id)
}
