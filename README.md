geoip-service
=============

A fast in-memory http microservice for looking up MaxMind GeoIP2 and GeoLite2 database.

This allows you to have your own IP to location lookup.

This implementation has been tested and handles more than 30,000(uncached)/70,000(cached) requests per second, and uses less than 100MB memory with no cache.

#Prerequisites
Requires a [go installation](https://golang.org/dl/).

A Database (choose one):
* [Free GeoLite 2 database](http://dev.maxmind.com/geoip/geoip2/geolite2/). Download the "MaxMind DB binary, gzipped" and unpack it.
* [GeoIP2 Downloadable Database](http://dev.maxmind.com/geoip/geoip2/downloadable/). This is a more detailed databse and should be compatible, but since I don't have access to that, I have been unable to verify this.

##Building the service
```go get github.com/klauspost/geoip-service```

This should build a "geoip-service" executable in your gopath.

##Running the service
Unpack the database to your current directory. Execute ```geoip-service -db=GeoLite2-City.mmdb```. This will start the service on port 5000 on your local computer.

##Using Docker
This project includes a Dockerfile. Place the source code and the .mmdb in the current directory.

To build the docker image, run:
```docker build -t geoip-service .```

To expose the service to other docker applications, assuming your database name is GeoLite2-City.mmdb run:
```docker run -it --rm geoip-service app db="GeoLite2-City.mmdb"```

To map the service to a port on your docker machine, run:
```docker run -it --rm -p 127.0.0.1:5000:5000 geoip-service app db="GeoLite2-City.mmdb"```

#Service Options

```
Usage of geoip-service:
  -db="GeoLite2-City.mmdb": File name of MaxMind GeoIP2 and GeoLite2 database
  -listen=":5000": Listen address and port, for instance 127.0.0.1:5000
  -lookup="city": Specify which value to look up. Can be 'city' or 'country' depending on which database you load.
  -pretty=false: Should output be formatted with newlines and intentation
  -threads=4: Number of threads to use. Defaults to number of detected cores
  -cache=0: How many seconds should requests be cached. Set to 0 to disable.

```
You can experiment with different cache options. It will store the results for a given IP for a number of seconds. That makes cached queries more than twice as fast as the initial lookup.

It will depend on your hardware and query pattern if it gives you a performance boost, and it may very well be faster not to have cache enabled. 

Also note, there is no RAM limit on the number of stores queries, so if you hit it with millions of completely different requests, your RAM use will rise significantly, so use with care.

#Using the service

Once the service is running, point your browser to ```http://localhost:5000/1.2.3.4```. You can replace "1.2.3.4" with the IP you would like to look up.

Currently the request above yields the following with "GeoLite2-City.mmdb" with the "pretty" option enabled:
```
{
  "Data": {
    "City": {
      "GeoNameID": 5804306,
      "Names": {
        "en": "Mukilteo",
        "ja": "ムキルテオ",
        "zh-CN": "马科尔蒂奥"
      }
    },
    "Continent": {
      "Code": "NA",
      "GeoNameID": 6255149,
      "Names": {
        "de": "Nordamerika",
        "en": "North America",
        "es": "Norteamérica",
        "fr": "Amérique du Nord",
        "ja": "北アメリカ",
        "pt-BR": "América do Norte",
        "ru": "Северная Америка",
        "zh-CN": "北美洲"
      }
    },
    "Country": {
      "GeoNameID": 6252001,
      "IsoCode": "US",
      "Names": {
        "de": "USA",
        "en": "United States",
        "es": "Estados Unidos",
        "fr": "États-Unis",
        "ja": "アメリカ合衆国",
        "pt-BR": "Estados Unidos",
        "ru": "Сша",
        "zh-CN": "美国"
      }
    },
    "Location": {
      "Latitude": 47.913,
      "Longitude": -122.3042,
      "MetroCode": 819,
      "TimeZone": "America/Los_Angeles"
    },
    "Postal": {
      "Code": "98275"
    },
    "RegisteredCountry": {
      "GeoNameID": 2077456,
      "IsoCode": "AU",
      "Names": {
        "de": "Australien",
        "en": "Australia",
        "es": "Australia",
        "fr": "Australie",
        "ja": "オーストラリア",
        "pt-BR": "Austrália",
        "ru": "Австралия",
        "zh-CN": "澳大利亚"
      }
    },
    "RepresentedCountry": {
      "GeoNameID": 0,
      "IsoCode": "",
      "Names": null,
      "Type": ""
    },
    "Subdivisions": [
      {
        "GeoNameID": 5815135,
        "IsoCode": "WA",
        "Names": {
          "en": "Washington",
          "es": "Washington",
          "fr": "État de Washington",
          "ja": "ワシントン州",
          "ru": "Вашингтон",
          "zh-CN": "华盛顿州"
        }
      }
    ],
    "Traits": {
      "IsAnonymousProxy": false,
      "IsSatelliteProvider": false
    }
  }
}
```

In case of success, the "Data" field will be filled. In case of an error, it will be put into an "Error" field, like this:

```
GET http://localhost:5000/1.2.3.4.5  ->

{
  "Error": "unable to decode ip"
}
```

For best performance you should utilize "keepalive" on the connection if you have a high throughput to avoid creating a new connection for every request. Consult the documentation on your language of choice for that.

#Credits
Uses  [GeoIP2 Reader for Go](https://github.com/oschwald/geoip2-golang) for its grunt work. This is the reason it is so fast!
