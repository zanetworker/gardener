// Code generated by go-swagger; DO NOT EDIT.

package size

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/metal-pod/metal-go/api/models"
)

// FromHardwareReader is a Reader for the FromHardware structure.
type FromHardwareReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *FromHardwareReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewFromHardwareOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewFromHardwareDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewFromHardwareOK creates a FromHardwareOK with default headers values
func NewFromHardwareOK() *FromHardwareOK {
	return &FromHardwareOK{}
}

/*FromHardwareOK handles this case with default header values.

OK
*/
type FromHardwareOK struct {
	Payload *models.V1SizeMatchingLog
}

func (o *FromHardwareOK) Error() string {
	return fmt.Sprintf("[POST /v1/size/from-hardware][%d] fromHardwareOK  %+v", 200, o.Payload)
}

func (o *FromHardwareOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1SizeMatchingLog)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewFromHardwareDefault creates a FromHardwareDefault with default headers values
func NewFromHardwareDefault(code int) *FromHardwareDefault {
	return &FromHardwareDefault{
		_statusCode: code,
	}
}

/*FromHardwareDefault handles this case with default header values.

Error
*/
type FromHardwareDefault struct {
	_statusCode int

	Payload *models.HttperrorsHTTPErrorResponse
}

// Code gets the status code for the from hardware default response
func (o *FromHardwareDefault) Code() int {
	return o._statusCode
}

func (o *FromHardwareDefault) Error() string {
	return fmt.Sprintf("[POST /v1/size/from-hardware][%d] fromHardware default  %+v", o._statusCode, o.Payload)
}

func (o *FromHardwareDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HttperrorsHTTPErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
