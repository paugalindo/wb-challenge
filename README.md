# Wallbox employees EV Pooling Service Challenge

Thank you for your interest in joining Wallbox engineering team. The next step in the process is for you 
to complete the code challenge described below.

We understand you may have other commitments and time constraints as well. Due to that, you will have 
**10 days** to finish the code challenge. It's not a hard limit. If for some reason you reach the deadline, 
but you are close to finish it, tell us, and it won't be a problem :smile:

Please, read carefully this document to know what we are expecting from you, the problem and how to submit the solution.

### What we're looking for

:warning: We don't want you to invest too many hours, we all have a life outside work :smile: **Tell us how many hours you dedicated in the readme**

The challenge is the way we have to know your technical skills. The _details matter_. A summary of what
we expect would be:
  * Write down your tech decisions in [explanation](#explanation) section. Being aware of it will help us review your challenge, tradeoffs you have considered, improvements in future iterations, etc...
  * A good testing strategy it's key for us. More than a high % of coverage we need somebody that understands what needs to be tested and what not.
  * We like well-structured code (who doesn't :stuck_out_tongue:)
  * We care very much about the domain layer. We expect rich logic there using meaningful language.
  * We like code that is easy to read and follow.
  * Follow best practices using anything you need to provide a solution that matches the problem.
  * Try to find the simplest solution possible, we recommend not applying extra features.
  * If you lack time to implement stuff feel free to comment on the pending things in the readme.
  * Our services follow DDD and Hexagonal, we need somebody with a good understanding of these from day 0 for this specific position.
  * Feel free to use any skeleton project that can help you to reduce the number of hours invested in developing this challenge. 

:exclamation: **We look for a solution that taking into account the challenges we face every day (concurrence, event-driven...) allows us to see what is the candidate able to do and his knowledge about this.** :exclamation:

Bonus points:
  * Show us what you know about event-driven but don't go too crazy with complex implementations.
  * Probably you've done a good unit test coverage, why not show us how you do acceptance testing? keep it simple, no need to cover everything.

:information_source: Feel free to use whatever programming language you think is best to solve the problem, and you are comfortable with.

## Problem

Design and implement a system to manage electric vehicle (EV) pooling.

Wallbox recently opened its new factory close to its headquarters. Communication
between teams is key and we often need to move from one place to another.
To achieve that, we have a fleet of EVs ready to use for our employees.
As saving energy is one of our main goals, we propose sharing cars with multiple
groups of people. This is an opportunity to optimize the use of resources by introducing car
pooling.

You have been assigned to build the car availability service that will be used
to track the available seats in cars.

Cars have a different amount of seats available. They can accommodate groups of
up to 4, 5 or 6 people.

People request cars in groups of 1 to 6. People in the same group want to ride
in the same car. You can assign any group to any car that has enough empty seats
for them. If it's not possible to accommodate them, they're willing to wait until
there's a car available for them. Once a car is available for a group, they should immediately 
enter and drive the car. You cannot ask them to change the car (i.e. swap them to make space for another group). 
The trip order should be "First come, first serve".

For example, a group of 6 people is waiting for a car. They cannot enter a car with less than 6 available seats 
(you can not split the group), so they need to wait. This means that smaller groups after them could enter a car with 
fewer available seats before them.

## API

To simplify the challenge and remove language restrictions, this service must
provide a REST API that will be used to interact with it.

This API must comply with the following contract:

### GET /status

Indicate the service has started up correctly and is ready to accept requests.

Responses:

* **200 OK** When the service is ready to receive requests.

### PUT /evs

Load the list of available EVs in the service and remove all previous data
(existing journeys and EVs). This method may be called more than once during
the life cycle of the service.

**Body** _required_ The list of EVs to load.

**Content Type** `application/json`

Sample:

```json
[
  {
    "id": 1,
    "seats": 4
  },
  {
    "id": 2,
    "seats": 6
  }
]
```

Responses:

* **200 OK** When the list is registered correctly.
* **400 Bad Request** When there is a failure in the request format, expected
  headers, or the payload can't be unmarshalled.

### POST /journey

A group of people requests to perform a journey.

**Body** _required_ The group of people that wants to perform the journey

**Content Type** `application/json`

Sample:

```json
{
  "id": 1,
  "people": 4
}
```

Responses:

* **200 OK** or **202 Accepted** When the group is registered correctly.
* **400 Bad Request** When there is a failure in the request format or the
  payload can't be unmarshalled.

### POST /dropoff

A group of people requests to be dropped off whether they traveled or not.

**Body** _required_ The ID of the group

**Content Type** `application/json`

Sample:

```json
{
  "id": 1
}
```

Responses:

* **200 OK** or **204 No Content** When the group is unregistered correctly.
* **404 Not Found** When the group cannot be found.
* **400 Bad Request** When there is a failure in the request format or the
  payload can't be unmarshalled.

### POST /locate

Given a group ID such as `ID=X`, return the car the group is traveling
with, or no car if they are still waiting to be served.

**Body** _required_ The ID of the group

**Content Type** `application/json`

Sample:

```json
{
  "id": 1
}
```

**Accept** `application/json`

Responses:

* **200 OK** With the car as the payload when the group is assigned to a car.
* **204 No Content** When the group is waiting to be assigned to a car.
* **404 Not Found** When the group cannot be found.
* **400 Bad Request** When there is a failure in the request format or the
  payload can't be unmarshalled.

## Evaluation criteria

The scoring system is partially automated because it needs to pass a series of automated checks 
and scoring before it gets evaluated by the team.

- The `testing` test step in the `.gitlab-ci.yml` must pass in master before you
  submit your solution. This is a public check that can be used to make sure that other tests
  will run successfully on your solution. **This step needs to run without
  modification**.

- ___"further tests"___ will be used to prove that your solution works correctly.
  These are not visible to you as a candidate and will be run once you submit
  the solution.

:information_source: If you consider that your solution is good enough even though some test is falling, 
don't hesitate to submit your solution. However, we encourage you to finish properly the test and try to 
accomplish a green test pipeline. But, it's up to you :grin:

## Tooling

Wallbox uses Gitlab and Gitlab CI for our backend development work.
In this repository, you may find a [.gitlab-ci.yml](./.gitlab-ci.yml) file which
contains some tooling that would simplify the setup and testing of the
deliverable.

Additionally, you will find a basic Dockerfile which you could use as a
baseline. Be sure to modify it as much as needed, but keep the exposed port
as it is.

**Remember that the entry point should bootstrap your application to be able to start receiving requests.**

Feel free to modify the repository as much as need to include or remove
dependencies.

:warning: The challenge needs to be self-contained so we can evaluate it. What does it mean? All dependencies (for example Redis, MySQL, wherever...)
should be inside docker image.

:warning: Avoid dependencies and tools that would require changes to the
`testing` step of [.gitlab-ci.yml](./.gitlab-ci.yml), such as
`docker-compose`

## What to do when I finish the challenge

Once you finish it, open a **Merge Request** and send a message to HR contact to let them know that 
your challenge can be evaluated.

## Explanation

I have designed the application using Domain-Driven Design, Hexagonal Architecture, CQS and event-driven flows. To prevent write concurrency, I have implemented a mutex. For the domain, I found it appropriate to model vehicles and their seats, as well as groups and people. Both seats and people are extendable, allowing for additional characteristics and functionalities, such as a baby seat.

### Project Structure
At the root of the project, you will find: Go dependencies, a script to initialize the database, a test file that demonstrates the application flow, a simple implementation of a command bus, the `/cmd` directory that contains the application's entry point and the dependency bootstrap and the `/internal` directory that holds the core application code.

Within `/internal`, the domain entities, domain events, and implementations of infrastructure and commands/queries are organized into their respective directories.

### External Dependencies
I have minimized external dependencies as much as possible. I've used Postgres to store the data because it was the database I was familiar with from my last job, NATS as the service for events due to its simplicity, and Gin to handle HTTP route multiplexing.

### Execution
I have added both Postgres and NATS to the provided Dockerfile, and to build and run the image, I have used the following commands from the root of the project:
 - `docker build -t wb-challenge -f ./_ci/build/coding-challenge/Dockerfile ./`
 - `docker run -d -p 80:80 -p 5432:5432 -p 4222:4222 --name wb-challenge wb-challenge`

To run the flow test, it is necessary to have the Docker container with the external dependencies up and running (without the main golang app).

### Testing
I have unit-tested the domain, commands and queries. However, due to time constraints, I have not tested all error cases for the commands/queries. Additionally, I opted not to implement specific integration tests for the repositories or the event producer/consumer, as I believe these are adequately covered by the main acceptance test, which simulates a full application flow. That said, as an improvement of this test, I should add testing to ensure that domain events are triggered within this flow, although this part is already tested unitarily in the commands.

### Areas for Improvement
Several aspects of the application still need enhancement:
 - Logging and Monitoring: Currently nonexistent. For a production-ready system, I consider this an essential component.
 - Environment Variables and Secrets: The current implementation lacks proper handling for these, and it is not ready for production use.
 - Application Initialization and Dependency Injection: The initialization and dependency injection process should be reworked. Specifically, the signaling logic for proper shutdown and the adoption of a package like Wire would make dependency injection more reusable and maintainable.
 - Testing: The testing coverage can be improved by adding more test cases and using a framework to organize them more effectively.
 - Command Bus: I would rewrite the command bus using generics to avoid type checks in the command handlers
 - Repository & View: To follow CQRS, repositories should be separated from the views, as they currently share implementation and different entities should be created for reading and writing.
 - Sagas: Given the eventual consistency of distributed transactions for vehicle assignment and group dropoff, it would be necessary to add compensation sagas to avoid potential inconsistencies in case the second part of the transaction is not executed properly.
 - Tooling: Additional tooling should also be added to set up different environments and dependencies as needed.

### Edge Cases
I've also encountered some edge cases that are not defined in the README and that I implemented as I thought best, such as what to return when requesting `/locate` a group that has already dropped off, or when a group of people who have already done a dropoff requests a journey again.

### Time Spent
In total, I will have dedicated about 24 hours of work spread over 4 days. I would have liked to dedicate a bit more time, but with 3 young daughters, the available time is limited.

## Help
If you need any specific help or doubt, please reach us sending an email to the following address: `backend.hiring.help@wallbox.com`, 
please add in the email's subject the challenge identifier (repository's name).

Good luck and happy programming!
