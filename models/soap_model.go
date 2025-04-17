package models

import "encoding/xml"

type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    Body     `xml:"Body"`
}

type Body struct {
	XMLName                        xml.Name                       `xml:"Body"`
	RequestOtpBySmsServiceResponse RequestOtpBySmsServiceResponse `xml:"RequestOtpBySmsServiceResponse"`
}

type RequestOtpBySmsServiceResponse struct {
	XMLName                      xml.Name `xml:"RequestOtpBySmsServiceResponse"`
	RequestOtpBySmsServiceResult string   `xml:"RequestOtpBySmsServiceResult"`
}

type VerifyOtpSOAPResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		VerifyOtpResponse struct {
			VerifyOtpResult string `xml:"VerifyOtpResult"` // Extracts the inner text value
		} `xml:"VerifyOtpResponse"`
	} `xml:"Body"`
}
