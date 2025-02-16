[![Go Reference](https://pkg.go.dev/badge/github.com/nutrixpos/pos.svg)](https://pkg.go.dev/github.com/nutrixpos/pos)

Nutrix is a point of sale management system RESTful api. It allows you to manage inventory, sales and products for your restaurant or shop.

You need to integrate with a separage GUI service to provide an interface for end users, you can find suggested GUIs in the [GUI](#gui) section.

> Currently nutrix supports only restaurant style cycle, where the order is sent first to the kitchen awaiting to be started then begin processing the inventory.

**OpenAPI 3.x Docs :**
[`modules/core/specs.api.yaml`](modules/core/specs.api.yaml)

You can use the openapi docs to run [`Swagger`](https://swagger.io/) docs and test the api and you can also use the docs with a mock server like [`Prism`](https://github.com/stoplightio/prism) to develop the frontend without the need of installing nutrix.



# Getting started

### docker-compose
You can use the [DevOps](https://github.com/nutrixpos/devops) repo to install the entire application with docker-compose including the ui and databases.

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


### GUI
Since nutrix is an api based project, you will need a GUI to let end users interact with the api, you are free to create your own GUI and integrate it with the api. if you have an api you can open a discussion with a reference url.

Following are suggested GUI(s) :

- [https://github.com/nutrixpos/pos-frontend](https://github.com/nutrixpos/pos-frontend)


