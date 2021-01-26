## About

- Name: Lewis Jessett
- Go Experience: 1 year of commercial experience

## Setting up project

1. API cannot log into Postgres
    - Task said debug before contacting us so it is part of the task.
        - api database password doesn't match postgres container
2. Docker container running :)
3. PSQL environment variables aren't needed on api container.
4. Open documentation and play with postman

## Development
1. Create a package so it can be used as a lib in another project
2. Create basic structure of client
3. Documentation for delete does not match deleting a non existant resource is a 204 not 404. (also additional errors like invalid guid)
4. Cannot do health checks with curl on api
5. Needs to be easily extendable for new endpoints.

## Test Strategy
- Mainly integration tests as we can easily test against the docker container api.
- Unit tests where applicable.

## Attribution

General structure ideas taken from

https://github.com/google/go-github/tree/master/github
https://github.com/stripe/stripe-go/
