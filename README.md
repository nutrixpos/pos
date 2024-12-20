[![Go Reference](https://pkg.go.dev/badge/github.com/elmawardy/nutrix.svg)](https://pkg.go.dev/github.com/elmawardy/nutrix)

Nutrix is a point of sale management system RESTful api. It allows you to manage inventory, sales and products for your restaurant or shop.

You need to integrate with a separage GUI service to provide an interface for end users, you can find suggested GUIs in the [GUI](#gui) section.

**OpenAPI 3.x Docs :**
[`modules/core/specs.api.yaml`](modules/core/specs.api.yaml)

# Getting started

### Prerequisites
- #### Mongo
    - Before running the web server or seeding, make sure to set the [MongoDB](https://www.mongodb.com/) credentials properly in **config.yaml**
- #### Zitadel
    -  Also make sure that a Zitadel instance is up and running, [Zitadel](https://zitadel.com/) is used for auth in the project.
    - Make sure to create a Zitadel api app inside your project and download the zitadel-key.json to a safe location
    - Make sure the domain, port and key path are set properly in **config.yaml**

        > **__Note__** that Zitadel needs to be reachable using the same domain from the backend and the frontend, docker produced notable connection issues regarding the hostname, which requires additional settings like adding a reverse proxy to be able to reach it using the same hostname as from the browser.

### Running the web server
- Run `go run .` in the backend root directory to run the server


### DB Seeding
- Run `go run . seed` in the backend directory which will prompt for entities to seed.
    > **Warning**  quiting the prompt with `ctrl+q` or `esc` will run the seeding process, if you want to quit, just deselect all the entities


### GUI
Since nutrix is an api based project, you will need a GUI to let end users interact with the api, you are free to create your own GUI and integrate it with the api. if you have an api you can open a discussion with a reference url.

Following are suggested GUI(s) :

- [https://github.com/elmawardy/nutrix-frontend](https://github.com/elmawardy/nutrix-frontend)


