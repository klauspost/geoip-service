geoip-service
=============

A fast in-memory http microservice for looking up MaxMind GeoIP2 and GeoLite2 database.

This allows you to have your own IP to location lookup.

This implementation has been tested and handles more than 30,000(uncached)/70,000(cached) requests per second, and uses less than 100MB memory with no cache.

[![Build Status][1]][2]
[1]: https://travis-ci.org/klauspost/geoip-service.svg
[2]: https://travis-ci.org/klauspost/geoip-service

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
There is a [Docker Repository](https://registry.hub.docker.com/u/klauspost/geoip-service/) set up to easily deploy the service as a docker app.

To fetch the docker image, run:
```docker pull klauspost/geoip-service```

To get the server running, you must add the geo-database as a volume to /data/geodb. That will enable you to update the file without rebuilding the docker image.

Therefor, if you have placed the database at ```/local/GeoLite2-City.mmdb```, you can run the service at:

```docker run --rm -v /local/GeoLite2-City.mmdb:/data/geodb.mmdb klauspost/geoip-service"```

To map the service to port 3000 on your host, run:

```docker run --rm -p 127.0.0.1:3000:5000 -v /local/GeoLite2-City.mmdb:/data/geodb.mmdb klauspost/geoip-service"```

If you want to specify additional command line parameters, you can run the program like this:

```docker run --rm -p 127.0.0.1:3000:5000 -v /local/GeoLite2-City.mmdb:/data/geodb.mmdb klauspost/geoip-service app db="/data/geodb.mmdb" pretty=true```


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
