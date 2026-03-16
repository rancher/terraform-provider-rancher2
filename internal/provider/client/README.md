# Client

The client for the Rancher provider creates the base http client for all communication with Rancher.

## Separation of concerns

The client needs to be separate from requests for that client.
The provider instantiates the client, which is a singleton, and each resource
builds and communicates request objects for the client.

This allows the resources to dynamically manage how and what they send to Rancher,
while the client manages how all communications work, eg. http vs https, TLS, etc.
The Client manages transport mechanisms while Requests manage what is transferred.

## Client Interface

The client interface gives a clear understanding of what methods a client should provide.
This allows the use of a TestClient which can be injected into the provider for testing.

## Http Client

The http client implements the client interface providing a net/http client.
It exposes all options of the http client and is what the provider will use in production.

## Test Client

The test client implements the client interface but doesn't provide an http client.
Instead, the test client logs requests and responses but doesn't do anything with them.
Tests are expected to inject the test client into the provider when instantiating it,
they inject response objects when the test exercises a resource's function.
