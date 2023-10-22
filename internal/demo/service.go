package demo

import (
	"context"
	"crypto/md5"
	"encoding/base64"

	"github.com/ItsNotGoodName/radiomux/internal/webrpc"
)

type StateService struct {
	webrpc.StateService
}

func NewStateService(stateService webrpc.StateService) StateService {
	return StateService{
		StateService: stateService,
	}
}

func (s StateService) StateMediaSet(ctx context.Context, req *webrpc.SetStateMedia) error {
	// Do not use user supplied URI
	if req.Uri != nil {
		hash := md5.Sum([]byte(*req.Uri))
		hashURI := "https://example.com/" + base64.RawURLEncoding.EncodeToString(hash[:])
		req.Uri = &hashURI
	}

	return s.StateService.StateMediaSet(ctx, req)
}
