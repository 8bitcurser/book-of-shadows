package handlers

import (
	"net/http"

	"book-of-shadows/views"
)

// KeeperDashboard renders the Keeper tools dashboard
func (h *Handler) KeeperDashboard(w http.ResponseWriter, r *http.Request) {
	component := views.KeeperDashboard()
	if err := component.Render(r.Context(), w); err != nil {
		h.logger.Printf("Failed to render keeper dashboard: %v", err)
		h.respondError(w, err)
	}
}

// ChaseTracker renders the chase tracker page
func (h *Handler) ChaseTracker(w http.ResponseWriter, r *http.Request) {
	component := views.ChaseTracker()
	if err := component.Render(r.Context(), w); err != nil {
		h.logger.Printf("Failed to render chase tracker: %v", err)
		h.respondError(w, err)
	}
}

// CombatTracker renders the combat tracker page
func (h *Handler) CombatTracker(w http.ResponseWriter, r *http.Request) {
	component := views.CombatTracker()
	if err := component.Render(r.Context(), w); err != nil {
		h.logger.Printf("Failed to render combat tracker: %v", err)
		h.respondError(w, err)
	}
}
