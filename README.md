# userService

postman - `https://www.postman.com/winter-meteor-514146/workspace/user-service`

## Installation

Make sure that docker is installed (for windows or mac install docker desktop)

1. Clone this repo
2. Copy `.env.template.local` file contents to file `.env.local`, and write variables values
3. Run `source .env.local`
4. Run `docker-compose up --build -d` to start backend
5. Run `make migrate` to apply all migrations to database
6. Run `make feed` to insert test data to database
7. Run `docker-compose ps` to check if all services are UP
