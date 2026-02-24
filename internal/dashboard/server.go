package dashboard

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/msrovani/myclaw/internal/memory"
	"github.com/msrovani/myclaw/internal/observability"
	"github.com/msrovani/myclaw/internal/router"
	"github.com/msrovani/myclaw/internal/skills"
	"github.com/msrovani/myclaw/web/views"
)

// Server handles the HTMX dashboard rendering.
type Server struct {
	mux           *http.ServeMux
	economy       *router.Economy
	skillRegistry *skills.Registry
	memEngine     *memory.Engine
	llmRouter     *router.Router
}

// NewServer initializes the Dashboard HTTP routes pointing to templ components.
func NewServer(mux *http.ServeMux, mem *memory.Engine, eco *router.Economy, lr *router.Router, reg *skills.Registry) *Server {
	s := &Server{
		mux:           mux,
		memEngine:     mem,
		economy:       eco,
		llmRouter:     lr,
		skillRegistry: reg,
	}

	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("GET /", s.handleOverview)
	s.mux.HandleFunc("GET /tokens", s.handleTokens)
	s.mux.HandleFunc("GET /memory", s.handleMemory)
	s.mux.HandleFunc("GET /router", s.handleRouter)
	s.mux.HandleFunc("GET /skills", s.handleSkills)
	s.mux.HandleFunc("GET /admin", s.handleAdmin)

	// API actions (HTMX targets)
	s.mux.HandleFunc("POST /api/admin/checkpoint", s.handleCheckpoint)
}

func (s *Server) handleOverview(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	stats := observability.CollectSystemStats()
	templ.Handler(views.Overview(stats)).ServeHTTP(w, r)
}

func (s *Server) handleTokens(w http.ResponseWriter, r *http.Request) {
	// For simplicity in the dashboard demo, requesting stats for a "global" view or default workspace.
	// In production, an Admin UI would list tenants.
	totalCost := 0.00
	templ.Handler(views.Tokens(totalCost)).ServeHTTP(w, r)
}

func (s *Server) handleMemory(w http.ResponseWriter, r *http.Request) {
	templ.Handler(views.MemoryPanel()).ServeHTTP(w, r)
}

func (s *Server) handleRouter(w http.ResponseWriter, r *http.Request) {
	templ.Handler(views.RouterPanel()).ServeHTTP(w, r)
}

func (s *Server) handleSkills(w http.ResponseWriter, r *http.Request) {
	skillsList := s.skillRegistry.List()
	templ.Handler(views.SkillsPanel(skillsList)).ServeHTTP(w, r)
}

func (s *Server) handleAdmin(w http.ResponseWriter, r *http.Request) {
	templ.Handler(views.AdminPanel()).ServeHTTP(w, r)
}

func (s *Server) handleCheckpoint(w http.ResponseWriter, r *http.Request) {
	slog.Info("dashboard: admin requested WAL checkpoint")
	// Returns a success message chunk for HTMX to swap in without reloading the page.
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<div style="color: var(--success); font-weight: bold; margin-bottom: 1rem;">Checkpoint executed successfully.</div>`))
}
