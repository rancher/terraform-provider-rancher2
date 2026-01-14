# Rancher Client

The Rancher Client resource controls how the Rancher provider communicates with Rancher.

## In Memory Resource

This resource is unusual in how it works because its job is to generate configurations in memory.
This gives the user a tranparent way to configure the http client used to talk to Rancher. It also 
provides a pathway for adding new clients, such as the testing client used in our unit tests, and
it provides clear logging for Rancher interactions.

## Why isn't this incorporated into the provider package?

Normally, Terraform providers allow the user to configure access and authentication at the provider level.
However, in most providers, authorization is managed outside of Terraform and expected to be handled before
the provider is configured. We want to support the workflow of installing Rancher with the Helm chart and
smoothly moving to configuring Rancher in the same Terraform "apply" process. We need the configuration of
the provider to follow resource dependency chain to enable that workflow. In previous iterations of this 
provider we had a "bootstrap" resource which accomplished a series of actions including this one. We have
chosen to break that down for transparency sake and to allow more control of how the provider is configured.
The bootstrap resource had weaknesses in how it worked forcing users to configure the provider multiple times.
This forced users to juggle provider configurations with provider aliases and often resulted in authentication bugs.

## Benefits

There are several benefits this design choice provides. Primarily, it allows users to change the configuration
of the client at apply time, respecting the resource dependency chain. This means that, rather than the provider
being loaded during read before any resources in your configuruation have a chance to trigger, you can use apply
time data in your client configuration. This enables a workflow of installing Rancher in the same "apply" as the
helm chart deployment. The most appealing secondary benefit to users is the ability to separate out authentication
logic, allowing users to configure multiple authentication mechanisms and configurations and configure resources to
use the client configured with the appropriate level of access and authorization for that resource. The most important
benefit for provider developers is the ability to code different APIs for each resource, allowing each resource to
be in control of which API it talks to. This allows developers to migrate resources from one API to another independently
of each other without risk of breaking unrelated resources.

## Read function rehydrates

Since this resource is actualized in memory, it is guaranteed that some read function won't find any clients,
which creates a perpetual apply cycle. To avoid this, the read function will actually rehydrate the client registry
from state. This is usually taboo for a read function since it shouldn't alter actual resources, but since this resource
is sort of internal for the provider it won't cause any adverse effects. The read function will look for and create
missing clients in the registry if they exist in state. This allows "create" to manage adding new clients from config and
"update" to manage updates from the config rather than clients missing from memory.

## Environment Pass-through

This resource can be configured or partially configured from environment variables. This potentially removes some of
the benefits of having a client resource that is separate from the provider, but it allows users to manage their
secrets separately from their Terraform config. Data from environment variables aren't added to state! This means 
that the environment variables used to configure the client will need to be available each time the provider is run.
The update function detects environment variables and will calculate differences between only the state and actual config,
however, if the environment variable isn't present then it won't have any way of knowing that it should ignore that 
variable when calculating differences and may attempt to update the client with empty or bad data.

## Testing doesn't use test client

The tests written don't use the test client for implementation because there is no need at this point for the test client.
The test client's job is to prevent the need for actual client/server interactions. It should be implemented and used in 
any resource where the provider sends requests or expects responses, which isn't necessary in this resource. This resource
only manages the client registry in memory.
