// Code generated by go-swagger; DO NOT EDIT.

package machine

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/metal-pod/metal-go/api/models"
)

// NewSetMachineStateParams creates a new SetMachineStateParams object
// with the default values initialized.
func NewSetMachineStateParams() *SetMachineStateParams {
	var ()
	return &SetMachineStateParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewSetMachineStateParamsWithTimeout creates a new SetMachineStateParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewSetMachineStateParamsWithTimeout(timeout time.Duration) *SetMachineStateParams {
	var ()
	return &SetMachineStateParams{

		timeout: timeout,
	}
}

// NewSetMachineStateParamsWithContext creates a new SetMachineStateParams object
// with the default values initialized, and the ability to set a context for a request
func NewSetMachineStateParamsWithContext(ctx context.Context) *SetMachineStateParams {
	var ()
	return &SetMachineStateParams{

		Context: ctx,
	}
}

// NewSetMachineStateParamsWithHTTPClient creates a new SetMachineStateParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewSetMachineStateParamsWithHTTPClient(client *http.Client) *SetMachineStateParams {
	var ()
	return &SetMachineStateParams{
		HTTPClient: client,
	}
}

/*SetMachineStateParams contains all the parameters to send to the API endpoint
for the set machine state operation typically these are written to a http.Request
*/
type SetMachineStateParams struct {

	/*Body*/
	Body *models.V1MachineState
	/*ID
	  identifier of the machine

	*/
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the set machine state params
func (o *SetMachineStateParams) WithTimeout(timeout time.Duration) *SetMachineStateParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the set machine state params
func (o *SetMachineStateParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the set machine state params
func (o *SetMachineStateParams) WithContext(ctx context.Context) *SetMachineStateParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the set machine state params
func (o *SetMachineStateParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the set machine state params
func (o *SetMachineStateParams) WithHTTPClient(client *http.Client) *SetMachineStateParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the set machine state params
func (o *SetMachineStateParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the set machine state params
func (o *SetMachineStateParams) WithBody(body *models.V1MachineState) *SetMachineStateParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the set machine state params
func (o *SetMachineStateParams) SetBody(body *models.V1MachineState) {
	o.Body = body
}

// WithID adds the id to the set machine state params
func (o *SetMachineStateParams) WithID(id string) *SetMachineStateParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the set machine state params
func (o *SetMachineStateParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *SetMachineStateParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}