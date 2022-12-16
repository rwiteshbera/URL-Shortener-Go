# URL Shortener

This is a URL shortener built with Go and Redis.

## What is a URL shortener?

A URL shortener is a tool that takes a long, complex URL and creates a shorter, simpler version of it. This shorter URL can then be used to redirect to the original URL, making it easier to share and remember. URL shorteners are often used to save space on social media platforms, where there is a character limit for posts, or to track clicks on links.

## Prerequisites

- [Docker](https://www.docker.com/)

## Installation

1. Clone the repository:

```bash
git clone git@github.com:rwiteshbera/URL-Shortener-Go.git
```
2. Build and start the containers:
```bash
cd URL-Shortener-Go
docker-compose up
```
This will build the Go server and Redis containers and start them in the background.

## Usage
To shorten a URL, make a POST request to http://localhost:5000/shorten with a JSON body containing the `url` and `expiry` field:
The request body should look like this:
```json
{
  "url":"https://github.com/rwiteshbera",
  "expiry": 12 
}
```
Mention how longer do you want to make your custom link accessible. Once that time period is over, the link will be invalid. Value given on `expiry` field will be considered as hour.

This will return a JSON response containing the shortened URL:
```bash
{
  "response": {
    "url": "https://github.com/rwiteshbera",
    "short": "localhost:5000/b6e5cd",
    "expiry": 12,
    "rate_limit": 9,
    "rate_limit_reset": 30
  }
}
```

<img width="764" alt="Screenshot_20221216_094011" src="https://user-images.githubusercontent.com/73098407/208021586-c1c8aa00-b01e-40fb-b007-a6e797be6fa9.png">

In order to prevent abuse of the URL shortening service, I have implemented rate limiting for all users. This means that you will be limited to creating no more than 10 shortened URLs within a 30 minute window. If you exceed this limit, your IP address will be blocked for a period of time.

The rate limiting for this service is implemented using Redis. For each IP address that makes a request to the /shorten endpoint, a counter is incremented in Redis. If the counter exceeds the rate limit (10 in a 30 minute window), the IP address is added to a Redis set of blocked IPs.

If you receive a response with a status code of 429 (Too Many Requests), it means that you have exceeded the rate limit and your IP address has been blocked. You will need to wait until the expiration time has passed before making any further requests.
