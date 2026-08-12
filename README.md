# Swarmpit agent

Swarmpit docker agent.

[![version](https://img.shields.io/github/release-pre/swarmpit/agent.svg)](https://github.com/swarmpit/agent/releases) 
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/swarmpit/agent/pulls)

## Run

```{r, engine='bash', count_lines}
docker run -d \
  --name agent \
  --volume /var/run/docker.sock:/var/run/docker.sock \
  swarmpit/agent:latest
```

### Parameters

- STATS_FREQUENCY - default to **30**
- EVENT_ENDPOINT - default to **http://app:8080/events**
- HEALTH_CHECK_ENDPOINT - default to **http://app:8080/version**
- SWARMPIT_EVENT_TOKEN - default to **empty** (no token sent)
- DEBUG_EVENT - default to **false**
- DEBUG_STATS - default to **false**
- DEBUG_HTTP - default to **false** (logs every incoming request, plus containers that have gone away)

### SWARMPIT_EVENT_TOKEN

Optional shared secret authorizing the agent's pushes to `EVENT_ENDPOINT`. When
empty, no token is sent and Swarmpit accepts the push as it always has. When set,
it is sent as the `X-Swarmpit-Event-Token` header and must match the Swarmpit
app's own `SWARMPIT_EVENT_TOKEN`, otherwise the push is rejected with 401.

It is deliberately not sent as `Authorization`: Swarmpit parses that header as a
JWT and rejects anything it cannot verify before the event endpoint's access rule
is reached.

Requires Swarmpit 2.1.5 or newer. Older versions ignore the header, so setting it
is harmless when pointing at one.

## Important!

In case you are deploying agent inside Swarmpit [stack](https://github.com/swarmpit/swarmpit/blob/master/docker-compose.yml)
with some sort of customization, make sure that Swarmpit service name (default to **app**) match domain name set by EVENT_ENDPOINT & HEALTH_CHECK_ENDPOINT. 