# WeatherAPI


## Build instructions

Build container

`docker build -t weather-api .`

Run it. Now the app runs on local port 8080.

`docker run -i -p 8080:8080 -e OPEN_WEATHER_API_KEY=<api_key> weather-api`

## API

There are 4 endpoints in total.

You can run them via postman collection provided in the repository

### Bookmarks

- GET /bookmarks

Endpoint will serve all stored bookmarks. There are essentially location names.
 ```
[
     {
         "name": "Piaseczno, PL"
     },
     {
         "name": "Washington, DC"
     }
 ]
```
- POST /bookmarks

Endpoint for adding bookmarks. You have to pass a json body to the request in form of:

```
{
	"name":"Piaseczno, PL"
}
```
If request is successful, 204 empty response will be returned.

### Weather query

- POST /query

Endpoint for fetching weather. You have to pass a json body to the request in form of:

```
{
	"name":"Piaseczno, PL"
}
```

This will a response along with a 200 OK status with body:

```
{
    "locationName": "Piaseczno, PL",
    "date": "2019-04-01",
    "type": "few clouds",
    "temp": 5.95,
    "min_temp": 3.33,
    "max_temp": 10
}
```

### Statistics

- GET /statistics?location={locationName}

You have to pass a location name in query, for instance: /statistics?location=Piaseczno, PL

This will return a 200 OK response along with a json body:

```
{
    "queriesCount": 1,
    "statistics": {
        "04-2019": {
            "temperatures": {
                "average": 5.95,
                "lowest": 3.33,
                "highest": 10
            },
            "typesOccurrences": {
                "few clouds": 1
            }
        }
    }
}
```

Statistics are grouped by month.





