package wizard

import (
	"log"
	"net/http"

	"book-of-shadows/storage"
	"book-of-shadows/views"
)

// Handler holds dependencies for wizard handlers
type Handler struct {
	store  storage.InvestigatorStore
	logger *log.Logger
}

// New creates a new wizard Handler with dependencies
func New(store storage.InvestigatorStore, logger *log.Logger) *Handler {
	return &Handler{
		store:  store,
		logger: logger,
	}
}

// BaseStep handles the base step of the wizard
func (h *Handler) BaseStep(w http.ResponseWriter, r *http.Request) {
	params := r.Context().Value("params").([]string)
	key := params[0]

	component := views.BaseStep(nil)
	if key != "" && key != "new" {
		investigator, err := h.store.GetInvestigator(r, key)
		if err != nil {
			h.logger.Printf("Failed to get investigator: %v", err)
			// Continue with nil investigator for new character
		} else {
			component = views.BaseStep(investigator)
		}
	}

	if err := component.Render(r.Context(), w); err != nil {
		h.logger.Printf("Failed to render base step: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// AttrStep handles the attributes step of the wizard
func (h *Handler) AttrStep(w http.ResponseWriter, r *http.Request) {
	params := r.Context().Value("params").([]string)
	key := params[0]

	investigator, err := h.store.GetInvestigator(r, key)
	if err != nil {
		h.logger.Printf("Failed to get investigator: %v", err)
		http.Error(w, "Investigator not found", http.StatusNotFound)
		return
	}

	component := views.AttrStep(investigator)
	if err := component.Render(r.Context(), w); err != nil {
		h.logger.Printf("Failed to render attr step: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// SkillStep handles the skills step of the wizard
func (h *Handler) SkillStep(w http.ResponseWriter, r *http.Request) {
	params := r.Context().Value("params").([]string)
	key := params[0]

	investigator, err := h.store.GetInvestigator(r, key)
	if err != nil {
		h.logger.Printf("Failed to get investigator: %v", err)
		http.Error(w, "Investigator not found", http.StatusNotFound)
		return
	}

	component := views.SkillStep(investigator)
	if err := component.Render(r.Context(), w); err != nil {
		h.logger.Printf("Failed to render skill step: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// TalentStep handles the talent selection step of the wizard
func (h *Handler) TalentStep(w http.ResponseWriter, r *http.Request) {
	params := r.Context().Value("params").([]string)
	key := params[0]

	investigator, err := h.store.GetInvestigator(r, key)
	if err != nil {
		h.logger.Printf("Failed to get investigator: %v", err)
		http.Error(w, "Investigator not found", http.StatusNotFound)
		return
	}

	component := views.TalentStep(investigator)
	if err := component.Render(r.Context(), w); err != nil {
		h.logger.Printf("Failed to render talent step: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}