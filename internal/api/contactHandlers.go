package api

import (
	"context"
	"time"

	"github.com/clintjedwards/basecoat/internal/models"
	"github.com/clintjedwards/basecoat/internal/storage"
	"github.com/clintjedwards/basecoat/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetContact returns a single contact by key
func (api *API) GetContact(ctx context.Context, request *proto.GetContactRequest) (*proto.GetContactResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.GetContactResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &proto.GetContactResponse{}, status.Error(codes.FailedPrecondition, "id required")
	}

	contactRaw, err := api.db.GetContact(api.db, account, request.Id)
	if err != nil {
		if err == storage.ErrEntityNotFound {
			return &proto.GetContactResponse{}, status.Error(codes.NotFound, "contact requested not found")
		}
		return &proto.GetContactResponse{}, status.Error(codes.Internal, "failed to retrieve contact from database")
	}

	contact := models.Contact{}
	contact.FromStorage(&contactRaw)

	return &proto.GetContactResponse{Contact: contact.ToProto()}, nil
}

// ListContacts returns a list of all contacts's metadata.
func (api *API) ListContacts(ctx context.Context, _ *proto.ListContactsRequest) (*proto.ListContactsResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.ListContactsResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	contactsRaw, err := api.db.ListContacts(api.db, account, 0, 0)
	if err != nil {
		return &proto.ListContactsResponse{}, status.Error(codes.Internal, "failed to retrieve contacts from database")
	}

	protoContacts := []*proto.Contact{}
	for _, contactRaw := range contactsRaw {
		var contact models.Contact
		contact.FromStorage(&contactRaw)
		protoContacts = append(protoContacts, contact.ToProto())
	}

	return &proto.ListContactsResponse{Contacts: protoContacts}, nil
}

// CreateContact registers a new contact
func (api *API) CreateContact(ctx context.Context, request *proto.CreateContactRequest) (*proto.CreateContactResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.CreateContactResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Name == "" {
		return &proto.CreateContactResponse{}, status.Error(codes.FailedPrecondition, "contact name required")
	}

	contact := models.NewContact(account, request.Name)
	contact.Phone = request.Phone
	contact.Email = request.Email

	err := api.db.InsertContact(api.db, contact.ToStorage())
	if err != nil {
		if err == storage.ErrEntityExists {
			return &proto.CreateContactResponse{}, status.Error(codes.AlreadyExists, "could not save contact; contact already exists")
		}
		log.Error().Err(err).Msg("could not save contact")
		return &proto.CreateContactResponse{}, status.Error(codes.Internal, "could not save contact")
	}

	log.Info().Str("id", contact.ID).Str("contact", contact.Name).Msg("contact created")

	return &proto.CreateContactResponse{
		Contact: contact.ToProto(),
	}, nil
}

// UpdateContact updates an already existing contact
func (api *API) UpdateContact(ctx context.Context, request *proto.UpdateContactRequest) (*proto.UpdateContactResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.UpdateContactResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &proto.UpdateContactResponse{}, status.Error(codes.FailedPrecondition, "contact id required")
	}

	err := api.db.UpdateContact(api.db, account, request.Id, storage.UpdatableContactFields{
		Name:     request.Name,
		Email:    request.Email,
		Phone:    request.Phone,
		Modified: ptr(time.Now().UnixMilli()),
	})
	if err != nil {
		if err == storage.ErrEntityNotFound {
			return &proto.UpdateContactResponse{}, status.Error(codes.NotFound, "contact requested not found")
		}
		log.Error().Err(err).Msg("could not save contact")
		return &proto.UpdateContactResponse{}, status.Error(codes.Internal, "could not save contact")
	}

	log.Debug().Str("id", request.Id).Msg("contact updated")
	return &proto.UpdateContactResponse{}, nil
}

func (api *API) DeleteContact(ctx context.Context, request *proto.DeleteContactRequest) (*proto.DeleteContactResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.DeleteContactResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &proto.DeleteContactResponse{}, status.Error(codes.FailedPrecondition, "contact id required")
	}

	err := api.db.DeleteContact(api.db, account, request.Id)
	if err != nil {
		if err == storage.ErrEntityNotFound {
			return &proto.DeleteContactResponse{}, status.Error(codes.NotFound, "could not delete contact; contact key not found")
		}

		log.Error().Err(err).Msg("could not delete contact")
		return &proto.DeleteContactResponse{}, status.Error(codes.Internal, "could not delete contact")
	}

	return &proto.DeleteContactResponse{}, nil
}
