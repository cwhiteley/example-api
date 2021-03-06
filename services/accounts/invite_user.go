package accounts

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/RichardKnop/example-api/util/response"
	"github.com/RichardKnop/example-api/logger"
)

// Handles requests to invite a new user
// POST /v1/invitations
func (s *Service) inviteUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the request context
	authenticatedUser, err := GetAuthenticatedUser(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// Request body cannot be nil
	if r.Body == nil {
		response.Error(w, "Request body cannot be nil", http.StatusBadRequest)
		return
	}

	// Read the request body
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Unmarshal the request body into the request prototype
	invitationRequest := new(InvitationRequest)
	if err = json.Unmarshal(payload, invitationRequest); err != nil {
		logger.ERROR.Printf("Failed to unmarshal invitation request: %s", payload)
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new invited user account
	invitation, err := s.InviteUser(authenticatedUser, invitationRequest)
	if err != nil {
		logger.ERROR.Printf("Invite user error: %s", err)
		response.Error(w, err.Error(), getErrStatusCode(err))
		return
	}

	// Create invitation response
	invitationResponse, err := NewInvitationResponse(invitation)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response
	response.WriteJSON(w, invitationResponse, 201)
}
