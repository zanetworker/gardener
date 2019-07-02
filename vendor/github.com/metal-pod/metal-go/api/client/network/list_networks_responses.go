// Code generated by go-swagger; DO NOT EDIT.

package network

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/metal-pod/metal-go/api/models"
)

// ListNetworksReader is a Reader for the ListNetworks structure.
type ListNetworksReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListNetworksReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewListNetworksOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewListNetworksDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListNetworksOK creates a ListNetworksOK with default headers values
func NewListNetworksOK() *ListNetworksOK {
	return &ListNetworksOK{}
}

/*ListNetworksOK handles this case with default header values.

OK
*/
type ListNetworksOK struct {
	Payload []*models.V1NetworkResponse
}

func (o *ListNetworksOK) Error() string {
	return fmt.Sprintf("[GET /v1/network][%d] listNetworksOK  %+v", 200, o.Payload)
}

func (o *ListNetworksOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListNetworksDefault creates a ListNetworksDefault with default headers values
func NewListNetworksDefault(code int) *ListNetworksDefault {
	return &ListNetworksDefault{
		_statusCode: code,
	}
}

/*ListNetworksDefault handles this case with default header values.

Error
*/
type ListNetworksDefault struct {
	_statusCode int

	Payload *models.HttperrorsHTTPErrorResponse
}

// Code gets the status code for the list networks default response
func (o *ListNetworksDefault) Code() int {
	return o._statusCode
}

func (o *ListNetworksDefault) Error() string {
	return fmt.Sprintf("[GET /v1/network][%d] listNetworks default  %+v", o._statusCode, o.Payload)
}

func (o *ListNetworksDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HttperrorsHTTPErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
