## Data quality check for recycling containers/centres


#Usage

`go build`

then run the executable.

#Required

The app relies on osmtogejson for data conversion
e.g. `npm install osmtogeojson`


The app will fetch data from overpass API convert them to GeoJSON (using osmtogejson) a serve them using Leaflet. 
It is preset for the Czech Republic, but that can be changed by setting different geocodes in `geocodes.go`.




Demo: [Czech Recycling Containers](https://thartek.alwaysdata.net/)


