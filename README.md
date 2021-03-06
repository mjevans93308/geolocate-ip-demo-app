# Go service to perform geoIP lookup

This service is to be used as an example app of getting a Go app up and running quickly with minimal overhead. It is an HTTP-based API that receives an IP address and a white list of countries. The API makes an external call to a geoIP lookup service, and returns a status code if the IP's country is in the whitelist.

# Web Service documentation
We have one main endpoint, located at `/api/v1/checkiplocation`
Example invocation via cURL for the service running locally:
```
curl --location --request POST 'localhost:8080/api/v1/checkiplocation' \
--header 'Accept: application/json' \
--header 'Authorization: Basic YWRtaW46cGFzc3dvcmQ=' \
--header 'Content-Type: application/json' \
--data-raw '{
    "ip_address": "71.229.168.95",
    "country_names": [
        "United States",
        "Mexico",
        "Canada"
    ]
}'
```
## Possible Return Codes
| HTTP Status Code  | Reason |
| ------------- | ------------- |
| `302` | Country lookup via geoIP lookup succeeded, and it matched a country in the whitelist  |
| `404`  | Country lookup via geoIP lookup succeeded, but it did not match a country in the whitelist  |
| `400`  | Request body was malformed or not all values provided in the payload were valid  |

This service is secured via Basic Authentication. The basic auth creds in the above example are `admin:password`, but can and should be changed in the .env file provided in the top-level project directory.
We assume that the IP address provided in the JSON payload will adhere to IPv4 structure, but the system can already handle IPv6 addresses as well, if the need ever arises.

We assume that the `country_names` field provided will be a list of strings. We could have asked for a comma-separated list like `"Mexico,Canada,...etc"` but thought that a list would be easier to parse given the existence of companies with spaces in their name, like the `United States` or `Republic of Korea`. Escaping spaces in a single string could be done, but it adds complexity when it can probably be avoided.

The same consideration to spaces in a countries' name applies in our decision to provide our endpoint as a `POST` request, so that we can ask for and expect a JSON object rather than building it as a `GET` request, and leaving the values in the URL itself as query parameters like `http://localhost:8080/api/v1/checkiplocation?ip=x.x.x.x&country_names=Canada,Mexico,Ireland`.

# Setup
This service will need to provide geoip credentials of the form `User_ID` / `License_Key` for us to utilize their 3rd party web service. Instructions for generating those credentials can be found at https://dev.maxmind.com/geoip/geolite2-free-geolocation-data#Access. Once generated, those credentials should be placed in the .env file with the Basic Authentication credentials. The `.env` file will resemble this:
```# basic auth creds
BASIC_AUTH_USERNAME=foo
BASIC_AUTH_PASSWORD=bar

# external geoIP lookup api creds
MAXMIND_USER_ID=hello
MAXIND_LICENSE_KEY=world

# environment
ENVIRONMENT=test
```
Additional environment variables that handle sensitive information should also be placed here. Check how we're using the credentials already in the `.env` file for examples of how to use an environment variable once it is loaded in the workspace.

Before release of this service, make sure the `ENVIRONMENT` var is set to `prod`. Setting it to `test` at the moment logs the IP address along with whether the request succeeded or failed.

# Running
## Local Development
Commands for running this as a standalone go application for testing purposes and local development are:
```
go build -o out/app main/main.go && ./out/app serve
```

## Docker
The assumption is this service will be run via Docker, and a Dockerfile is provided in the top-level project directory. Commands for building the Docker image and running that image are as follows:
```
docker build -f Dockerfile -t geolocate-ip-demo-app:0.0.1 .
docker run --env-file ./geolocate-ip-demo-app.env --publish 8080:8080 geolocate-ip-demo-app:0.0.1
```
To get the .env file into the docker workspace we provide the .env file as a arg in the `docker run` command, and then inside the Dockerfile we copy it into our final image via 
```
COPY --from=builder /build/geolocate-ip-demo-app.env .
```
This allows us to not expose any sensitive information to console logging during build and release. The .env file cannot be committed to VCS for obvious reasons, so check the team's sharepoint for the latest version to use in your local workspace.

# Considerations and Assumptions
In order to avoid issues keeping the mapping data up to date we are directly utilizing the geolite2 web service (documented at https://dev.maxmind.com/geoip/docs/web-services?lang=en) rather than downloading a local copy of their database. This allows us to always be on the latest data they have.

Some pros/cons to this decision:
- Allows us to abstract away our data upkeep and validation needs to a 3rd party. This can be seen as either a pro or a con depending on how much we trust the 3rd party.
- Need to take into consideration changes to the 3rd party API, although they state that any future api development will be versioned according.
- Need to be mindful of how many geoip API calls we are making from our service due to the possibility of web service throttling on their end. Currently we cache each unique ip -> country mapping in a static map in Go, but obviously we would want to be mindful of the memory footprint of this map as we process calls. Using a standalone storage solution specifically tailored to our need (Key-Value schema and fast reads) like Redis would be a good long term solution, especially since we could configure the Redis table with an entry lifetime limit to manage the size of the table.
- If their service still isn't up to handling the breadth of calls we encounter then we might need to explore a different service architecture. We're currently utilizing the free service, but paying for a more robust and/or accurate version could be considered.
- If there comes a time when using the geoip web service isn't possible, we could explore automating the process of downloading their entire ip -> country mapping database and storing it in our environment for our use. This approach would be more complicated and involve the use of our own storage solution like DynamoDB or Redis. This would solve the throttling issue but we would still be dependent on a 3rd party.

# Future Ideas
- Implement monitoring of API statistics with something like statsd and tracing with Jaeger
- Implement RPC support using gRPC or Twirp
- Examine whether a more robust security device besides Basic Authentication should be explored. JWTs might be an option, but could also just be too much hassle for what they would provide.
- Better testing coverage, specifically around the api package
- Work with team where requests will originate to start sending `X-Request-ID` header so that we can log specific request details. Currently the only useful identifying data being sent with each request is the IP address, which is too sensitive for logging in production code.
- Add makefile support.
- Replace static Go map with standalone storage solution whose memory footprint can be more easily monitored and managed.

# Main 3rd party libraries used:
- go.uber.org/zap for logging
- github.com/spf13/viper for environment config support
- github.com/spf13/cobra for standalone CLI support
- inet.af/netaddr for IP validation/parsing support
- github.com/gin-gonic/gin for HTTP router and basic web service needs.
