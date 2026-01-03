package application

import (
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/embedding"
)

// SemanticSearchResult represents a search result with similarity score
type SemanticSearchResult struct {
	Facility   *domain.Facility `json:"facility"`
	Similarity float32          `json:"similarity"`
}

// SemanticSearchUseCase handles semantic search for facilities
type SemanticSearchUseCase struct {
	repo     domain.FacilityRepository
	embedder *embedding.EmbeddingService
}

// NewSemanticSearchUseCase creates a new semantic search use case
func NewSemanticSearchUseCase(repo domain.FacilityRepository) *SemanticSearchUseCase {
	return &SemanticSearchUseCase{
		repo:     repo,
		embedder: embedding.NewEmbeddingService(),
	}
}

// Search performs a semantic search for facilities matching the query
// Example queries: "canchas techadas para lluvia", "piscina climatizada", "tenis noche"
func (uc *SemanticSearchUseCase) Search(clubID, query string, limit int) ([]*SemanticSearchResult, error) {
	if limit <= 0 {
		limit = 10
	}

	// Generate embedding for query
	queryEmbedding := uc.embedder.GenerateEmbedding(query)

	// Search in database using pgvector
	results, err := uc.repo.SemanticSearch(clubID, queryEmbedding, limit)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	searchResults := make([]*SemanticSearchResult, len(results))
	for i, r := range results {
		searchResults[i] = &SemanticSearchResult{
			Facility:   r.Facility,
			Similarity: r.Similarity,
		}
	}

	return searchResults, nil
}

// GenerateAndStoreEmbedding generates an embedding for a facility and stores it
func (uc *SemanticSearchUseCase) GenerateAndStoreEmbedding(facility *domain.Facility) error {
	// Build text from facility data for embedding
	text := buildFacilityText(facility)

	// Generate embedding
	emb := uc.embedder.GenerateEmbedding(text)

	// Store in database
	return uc.repo.UpdateEmbedding(facility.ID, emb)
}

// GenerateAllEmbeddings generates embeddings for all facilities (batch operation)
func (uc *SemanticSearchUseCase) GenerateAllEmbeddings(clubID string) (int, error) {
	facilities, err := uc.repo.List(clubID, 1000, 0)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, fac := range facilities {
		if err := uc.GenerateAndStoreEmbedding(fac); err != nil {
			// Log error but continue with other facilities
			continue
		}
		count++
	}

	return count, nil
}

// buildFacilityText creates a searchable text representation of a facility
func buildFacilityText(f *domain.Facility) string {
	text := f.Name + " "
	text += string(f.Type) + " "

	// Add specifications
	if f.Specifications.SurfaceType != nil {
		text += *f.Specifications.SurfaceType + " "
	}
	if f.Specifications.Lighting {
		text += "iluminacion nocturno noche "
	}
	if f.Specifications.Covered {
		text += "techado cubierto lluvia "
	}
	for _, eq := range f.Specifications.Equipment {
		text += eq + " "
	}

	// Add location
	text += f.Location.Name + " "
	text += f.Location.Description

	return text
}
