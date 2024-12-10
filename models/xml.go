package models

import "encoding/xml"

type WeatherData struct {
	XMLName     xml.Name      `xml:"ObservationCollection"`
	Om          string        `xml:"om,attr"`
	Xmlns       string        `xml:"xmlns,attr"`
	Gml         string        `xml:"gml,attr"`
	Xlink       string        `xml:"xlink,attr"`
	Xsi         string        `xml:"xsi,attr"`
	Observation []Observation `xml:"member>Observation"`
}

type Observation struct {
	Metadata          Metadata          `xml:"metadata"`
	SamplingTime      TimeInfo          `xml:"samplingTime"`
	ResultTime        TimeInfo          `xml:"resultTime"`
	Procedure         Procedure         `xml:"procedure"`
	ObservedProperty  ObservedProperty  `xml:"observedProperty"`
	FeatureOfInterest FeatureOfInterest `xml:"featureOfInterest"`
	Result            Result            `xml:"result"`
}

type Metadata struct {
	General                General                `xml:"set>general"`
	IdentificationElements IdentificationElements `xml:"set>identification-elements"`
}

type General struct {
	Author  Author  `xml:"author"`
	Dataset Dataset `xml:"dataset"`
	Phase   Phase   `xml:"phase"`
	ID      ID      `xml:"id"`
	Parent  Parent  `xml:"parent"`
}

type IdentificationElements struct {
	Element []Element `xml:"element"`
}

type Author struct {
	Build   string `xml:"build,attr"`
	Name    string `xml:"name,attr"`
	Version string `xml:"version,attr"`
}

type Dataset struct {
	Name string `xml:"name,attr"`
}

type Phase struct {
	Name string `xml:"name,attr"`
}

type ID struct {
	Href string `xml:"href,attr"`
}

type Parent struct {
	Href string `xml:"href,attr"`
}

type TimeInfo struct {
	TimeInstant TimeInstant `xml:"TimeInstant"`
}

type TimeInstant struct {
	TimePosition string `xml:"timePosition"`
}

type Procedure struct {
	Href string `xml:"href,attr"`
}

type ObservedProperty struct {
	RemoteSchema string `xml:"remoteSchema,attr"`
}

type FeatureOfInterest struct {
	FeatureCollection FeatureCollection `xml:"FeatureCollection"`
}

type FeatureCollection struct {
	Location Location `xml:"location"`
}

type Location struct {
	Point Point `xml:"Point"`
}

type Point struct {
	Pos string `xml:"pos"`
}

type Result struct {
	Elements Elements `xml:"elements"`
}

type Elements struct {
	Element []Element `xml:"element"`
}

type Element struct {
	Name      string      `xml:"name,attr"`
	Uom       string      `xml:"uom,attr"`
	Value     string      `xml:"value,attr"`
	Qualifier []Qualifier `xml:"qualifier"`
}

type Qualifier struct {
	Name  string `xml:"name,attr"`
	Uom   string `xml:"uom,attr"`
	Value string `xml:"value,attr"`
}
