package main

import (
    "encoding/xml"
    "fmt"
    "os"
)

const (
    SoapSchema    string = "http://schemas.xmlsoap.org/soap/envelope/"
    XSISchema     string = "http://www.w3.org/2001/XMLSchema-instance"
    XSDSchema     string = "http://www.w3.org/2001/XMLSchema"
    XMLNSSchema   string = "https://advertising.criteo.com/API/v201305"
    AuthToken     string = "763145448959408512"
    AppToken      string = "5458719955272208384"
    ClientVersion string = "sephora-criteoRep-1.0"
)

func main() {
    type Header struct {
        XMLNSSchema   string `xml:"xmlns,attr"`
        AuthToken     string `xml:"authToken,omitempty"`
        AppToken      string `xml:"appToken,omitempty"`
        ClientVersion string `xml:"clientVersion,omitempty"`
    }
    type GetAccountFunc struct {
        XMLNSSchema string `xml:"xmlns,attr"`
    }
    type Call struct {
        XMLName        xml.Name        `xml:"soap:Envelope"`
        SoapSchema     string          `xml:"xmlns:soap,attr"`
        XSISchema      string          `xml:"xmlns:xsi,attr"`
        XSDSchema      string          `xml:"xmlns:xsd,attr"`
        Header         Header          `xml:"soap:Header>apiHeader"`
        GetAccountFunc *GetAccountFunc `xml:"soap:Body>getAccount,omiempty"`
    }

    v := &Call{
        SoapSchema: SoapSchema,
        XSISchema:  XSISchema,
        XSDSchema:  XSDSchema,
        Header: Header{
            XMLNSSchema:   XMLNSSchema,
            AuthToken:     AuthToken,
            AppToken:      AppToken,
            ClientVersion: ClientVersion,
        },
    }
    v.GetAccountFunc = &GetAccountFunc{
        XMLNSSchema: XMLNSSchema,
    }

    enc := xml.NewEncoder(os.Stdout)
    enc.Indent("  ", "    ")
    if err := enc.Encode(v); err != nil {
        fmt.Printf("error: %v\n", err)
    }

}

//`
//  <person id="13">
//      <name>
//          <first>John</first>
//          <last>Doe</last>
//      </name>
//      <age>42</age>
//      <Married>false</Married>
//      <City>Hanga Roa</City>
//      <State>Easter Island</State>
//      <!-- Need more details. -->
//  </person>
//`
