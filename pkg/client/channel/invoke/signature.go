/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package invoke

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/common/verifier"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"

	"github.com/pkg/errors"
)

//NewSignatureValidationHandler returns a handler that validates an endorsement
func NewSignatureValidationHandler(next ...Handler) *SignatureValidationHandler {
	return &SignatureValidationHandler{next: getNext(next)}
}

//SignatureValidationHandler for transaction proposal response filtering
type SignatureValidationHandler struct {
	next Handler
}

//Handle for Filtering proposal response
func (f *SignatureValidationHandler) Handle(requestContext *RequestContext, clientContext *ClientContext) {
	//Filter tx proposal responses
	err := f.validate(requestContext, clientContext)
	if err != nil {
		requestContext.Error = errors.WithMessage(err, "signature validation failed")
		return
	}

	// Delegate to next step if any
	if f.next != nil {
		f.next.Handle(requestContext, clientContext)
	}
}

func (f *SignatureValidationHandler) validate(requestContext *RequestContext, ctx *ClientContext) error {
	// defer utils.TimeCost("signature validation ", string(requestContext.Response.TransactionID))()
	for _, r := range requestContext.Response.Responses {
		if err := verifyProposalResponse(r, ctx); err != nil {
			return err
		}
	}

	return nil
}

func verifyProposalResponse(res *fab.TransactionProposalResponse, ctx *ClientContext) error {
	sv := &verifier.Signature{Membership: ctx.Membership}
	return sv.Verify(res)
}
