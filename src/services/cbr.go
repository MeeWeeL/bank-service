package services

import (
    "bytes"
    "fmt"
    "net/http"
    "github.com/beevik/etree"
)

type CBRService struct {
    soapEndpoint string
}

func NewCBRService() *CBRService {
    return &CBRService{
        soapEndpoint: "https://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx",
    }
}

func (s *CBRService) GetKeyRate() (float64, error) {
    soapRequest := `<?xml version="1.0" encoding="utf-8"?>
    <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
      <soap:Body>
        <KeyRate xmlns="http://web.cbr.ru/"/>
      </soap:Body>
    </soap:Envelope>`

    resp, err := http.Post(s.soapEndpoint, "text/xml", bytes.NewBufferString(soapRequest))
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()

    doc := etree.NewDocument()
    if _, err := doc.ReadFrom(resp.Body); err != nil {
        return 0, err
    }

    rate := doc.FindElement("//KeyRateResult")
    if rate == nil {
        return 0, fmt.Errorf("key rate not found in response")
    }

    var result float64
    _, err = fmt.Sscanf(rate.Text(), "%f", &result)
    return result, err
}