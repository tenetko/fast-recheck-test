package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/KosyanMedia/delta/pkg/iata"
	"github.com/KosyanMedia/delta/pkg/types/integration"
)

type Response struct {
	Offers []*Offer `xml:"variant"`
}

type Offer struct {
	Segments []*Segment `xml:"segment"`
}

type Segment struct {
	Flights []*Flight `xml:"flight"`
}

type Flight struct {
	Origin           string `xml:"departure"`
	Destination      string `xml:"arrival"`
	RecheckBaggage   bool   `xml:"baggageRecheck"`
	VirtualInterline *bool  `xml:"virtualInterline"`
}

func Parse(fileName string, recheckBaggageAfter bool, virtualInterlineAfter bool) ([]*integration.FlightLeg, [][]*integration.TransferTerms) {
	xmlFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}

	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)
	var res Response
	xml.Unmarshal(byteValue, &res)

	// Соберём flight_legs и flight_terms только из первого сегмента первого оффера, для теста
	segment := res.Offers[0].Segments[0]
	legs := make([]*integration.FlightLeg, 0)
	transferTerms := make([][]*integration.TransferTerms, 1)

	for flightIdx, flight := range segment.Flights {

		// Если в конфиге указан флаг recheckBaggageAfter == true и если мы нашли флайт с признаком речека,
		// то перемещаем признак речека в предыдущий флайт, а в текущем флайте меняем признак речека на false.

		if recheckBaggageAfter && flight.RecheckBaggage && flightIdx > 0 {
			segment.Flights[flightIdx-1].RecheckBaggage = true
			segment.Flights[flightIdx].RecheckBaggage = false
		}

		// Попутно сделаем то же самое для признаков интерлайна, чтобы метод IsVirtualInterline()
		// возвращал правильные значения признака интерлайна после сдвига индекса.
		//
		// Если выполняются все следующие условия:
		// - в конфиге указан флаг virtualInterlineAfter == true ;
		// - партнёр передал нам признаки интерлайна (во флайте есть непустой тег VirtualInterline);
		// - этот тег равен true;
		//
		// то переещаем признак интерлайна в предыдущий флайт, а в текущем флайте меняем признак речека на false.

		if virtualInterlineAfter && flight.VirtualInterline != nil && *flight.VirtualInterline && flightIdx > 0 {
			*segment.Flights[flightIdx-1].VirtualInterline = true
			*segment.Flights[flightIdx].VirtualInterline = false
		}

		// Если флаги recheckBaggageAfter или virtualInterlineAfter == false, то ничего не делаем,
		//признаки речека и интерлайна уже стоят как надо.
	}

	// Пройдёмся по массиву флайтов и сформируем массив FlightLegs:

	for _, flight := range segment.Flights {
		leg := &integration.FlightLeg{
			Origin:         iata.NewLocationIATACode(flight.Origin),
			Destination:    iata.NewLocationIATACode(flight.Destination),
			RecheckBaggage: flight.RecheckBaggage,
		}

		legs = append(legs, leg)
	}

	// Ещё раз пройдёмся по массиву флайтов и сформируем массив TransferTerms:

	for flightIdx, _ := range segment.Flights {
		segmentIdx := 0 // первый и единственный

		if flightIdx > 0 {
			idx := flightIdx - 1

			// Эти строки мы можем убрать из исходного кода фаста,
			// т.к. мы уже должны были сделать перестановку признаков интерлайна
			// в правильном порядке в строках 70-71

			//if virtualInterlineAfter {
			//idx = flightIdx
			//}

			transferTerms[segmentIdx] = append(transferTerms[segmentIdx], &integration.TransferTerms{
				IsVirtualInterline: segment.Flights[idx].IsVirtualInterline(),
			})
		}
	}

	return legs, transferTerms
}

func (f Flight) IsVirtualInterline() bool {
	if f.VirtualInterline == nil {
		return f.RecheckBaggage
	}
	return *f.VirtualInterline
}

func main() {
	legs, transferTerms := Parse("xml_vi_rb/false-true-false.xml", true, true)
	for _, leg := range legs {
		fmt.Println(leg)
	}

	for _, terms := range transferTerms {
		for _, element := range terms {
			fmt.Println(element)
		}
	}
}
