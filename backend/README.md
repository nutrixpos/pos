# Getting started :

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


#### DB Seeding
- Run `go run . seed` in the backend directory which will prompt for entities to seed.
    > **Warning**  quiting the prompt with `ctrl+q` or `esc` will run the seeding process, if you want to quit, just deselect all the entities



