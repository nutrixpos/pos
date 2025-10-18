[![Go Reference](https://pkg.go.dev/badge/github.com/nutrixpos/pos.svg)](https://pkg.go.dev/github.com/nutrixpos/pos)

Nutrix is a point of sale management system RESTful api. It allows you to manage inventory, sales and products for your restaurant or shop.

You need to integrate with a GUI service to provide an interface for end users, you can find the standard GUI in the [frontend](./frontend) dir, it's a vuejs app, you can build it and put the content of the dist into /mnt/frontend

> Currently nutrix supports only restaurant style cycle, where the order is sent first to the kitchen awaiting to be started then begin processing the inventory.

**OpenAPI 3.x Docs :**
[`modules/core/specs.api.yaml`](modules/core/specs.api.yaml)

You can use the openapi docs to run [`Swagger`](https://swagger.io/) docs and test the api and you can also use the docs with a mock server like [`Prism`](https://github.com/stoplightio/prism) to develop the frontend without the need of installing nutrix.



# Getting started
To install the entire system including [posui](https://github.com/nutrixpos/posui) and [zitadel](https://zitadel.com) kindly follow the steps mentioned in the official website : [https://nutrixpos.com/getting_started.html](https://nutrixpos.com/getting_started.html)


### docker-compose
### Dependencies
- #### Mongo
    - Before running the web server or seeding, make sure to set the [MongoDB](https://www.mongodb.com/) credentials properly in **config.yaml**
- #### Zitadel
    -  Also make sure that a Zitadel instance is up and running, [Zitadel](https://zitadel.com/) is used for auth in the project.
    - Make sure to create a Zitadel api app inside your project and download the zitadel-key.json to a safe location
    - Make sure the domain, port and key path are set properly in **config.yaml**

        > **__Note__** that Zitadel needs to be reachable using the same domain from the backend and the frontend, docker produced notable connection issues regarding the hostname, which requires additional settings like adding a reverse proxy to be able to reach it using the same hostname as from the browser.
- #### google-chrome command (or equivalent on linux)
    - test by running `google-chrome --version` (tested on linux)

### Running the web server
- Run `go run .` in the backend root directory to run the server


### DB Seeding
- Run `go run . seed` in the backend directory which will prompt for entities to seed.
    > **Warning**  quiting the prompt with `ctrl+q` or `esc` will run the seeding process, if you want to quit, just deselect all the entities

