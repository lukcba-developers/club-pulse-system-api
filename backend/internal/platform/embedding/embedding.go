package embedding

import (
	"math"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Dimension defines the embedding vector size
const Dimension = 256

// EmbeddingService generates text embeddings using a lightweight TF-IDF approach
// This is a zero-cost alternative to paid embedding APIs (OpenAI, etc.)
type EmbeddingService struct {
	vocabulary map[string]int
	idf        map[string]float32
}

// NewEmbeddingService creates a new embedding service with a base vocabulary
func NewEmbeddingService() *EmbeddingService {
	// Base vocabulary for sports/facilities domain
	baseVocab := []string{
		// Sports
		"tenis", "futbol", "basketball", "natacion", "golf", "paddle", "squash", "volleyball",
		"badminton", "handball", "hockey", "rugby", "atletismo", "gimnasia", "yoga", "pilates",
		"crossfit", "spinning", "aerobics", "boxeo", "karate", "judo", "taekwondo",
		// Facility types
		"cancha", "piscina", "gimnasio", "campo", "pista", "sala", "court", "pool", "gym", "field",
		// Features
		"techado", "cubierto", "iluminacion", "nocturno", "noche", "climatizado", "exterior",
		"interior", "cesped", "sintetico", "natural", "arcilla", "cemento", "madera", "parquet",
		"lluvia", "sol", "sombra", "ventilacion", "aire", "acondicionado",
		// Equipment
		"raqueta", "pelota", "red", "arco", "aro", "colchoneta", "pesas", "maquinas",
		// Size
		"grande", "mediano", "pequeno", "profesional", "amateur", "olimpico", "estandar",
		// Time
		"manana", "tarde", "noche", "madrugada", "horario", "disponible", "reserva",
		// Quality
		"premium", "vip", "exclusivo", "familiar", "infantil", "adultos", "mayores",
	}

	vocab := make(map[string]int)
	for i, word := range baseVocab {
		vocab[word] = i % Dimension
	}

	return &EmbeddingService{
		vocabulary: vocab,
		idf:        make(map[string]float32),
	}
}

// GenerateEmbedding creates a vector representation of the input text
func (s *EmbeddingService) GenerateEmbedding(text string) []float32 {
	embedding := make([]float32, Dimension)

	// Tokenize and normalize
	tokens := s.tokenize(text)
	if len(tokens) == 0 {
		return embedding
	}

	// Count term frequency
	tf := make(map[string]int)
	for _, token := range tokens {
		tf[token]++
	}

	// Build embedding using vocabulary positions
	for token, count := range tf {
		// Get vocabulary position (hash if not in vocab)
		pos := s.getPosition(token)

		// TF component (log normalization)
		tfWeight := float32(1 + math.Log(float64(count)))

		// IDF component (use 1.0 if unknown, allowing unseen terms)
		idfWeight := float32(1.0)
		if w, ok := s.idf[token]; ok {
			idfWeight = w
		}

		// Add to embedding position
		embedding[pos] += tfWeight * idfWeight
	}

	// L2 normalize the embedding
	s.normalize(embedding)

	return embedding
}

// tokenize splits text into normalized tokens
func (s *EmbeddingService) tokenize(text string) []string {
	// Lowercase
	text = strings.ToLower(text)

	// Remove accents
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	text, _, _ = transform.String(t, text)

	// Split by non-alphanumeric
	re := regexp.MustCompile(`[^a-z0-9]+`)
	words := re.Split(text, -1)

	// Filter empty and short words
	var tokens []string
	for _, word := range words {
		if len(word) > 2 {
			tokens = append(tokens, word)
		}
	}

	return tokens
}

// getPosition returns the embedding dimension position for a token
func (s *EmbeddingService) getPosition(token string) int {
	if pos, ok := s.vocabulary[token]; ok {
		return pos
	}
	// Hash unknown tokens to a position
	return s.hash(token) % Dimension
}

// hash creates a simple hash for unknown tokens
func (s *EmbeddingService) hash(token string) int {
	h := 0
	for _, c := range token {
		h = 31*h + int(c)
	}
	if h < 0 {
		h = -h
	}
	return h
}

// normalize applies L2 normalization to the embedding
func (s *EmbeddingService) normalize(embedding []float32) {
	var sum float64
	for _, v := range embedding {
		sum += float64(v * v)
	}

	if sum == 0 {
		return
	}

	norm := float32(math.Sqrt(sum))
	for i := range embedding {
		embedding[i] /= norm
	}
}

// CosineSimilarity calculates the cosine similarity between two embeddings
func CosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return float32(dot / (math.Sqrt(normA) * math.Sqrt(normB)))
}
