[![Go Reference](https://pkg.go.dev/badge/github.com/nutrixpos/pos.svg)](https://pkg.go.dev/github.com/nutrixpos/pos)

![screenshot of the application](https://elmawardy.sirv.com/Images/nutrix_wallpaper2-min.png)

Nutrix is a point of sale management system. It allows you to manage inventory, sales and products for your restaurant or shop.


> :warning: Warning : 
> NutrixPOS is currently in active development, and we are continuously making changes, updates, and working on new features and improvements. Please be aware that until a stable release is reached, backward compatibility is not guaranteed. We make every effort to maintain compatibility.


You need to integrate with a GUI service to provide an interface for end users, you can find the standard GUI in the [frontend](./frontend) dir, it's a vuejs app, you can build it and put the content of the dist into /mnt/frontend

> Currently nutrix supports only restaurant style cycle, where the order is sent first to the kitchen awaiting to be started then begin processing the inventory.

**OpenAPI 3.x Docs :**
[`modules/core/specs.api.yaml`](modules/core/specs.api.yaml)

You can use the openapi docs to run [`Swagger`](https://swagger.io/) docs and test the api and you can also use the docs with a mock server like [`Prism`](https://github.com/stoplightio/prism) to develop the frontend without the need of installing nutrix.

Check the [`Getting started`](https://nutrixpos.com/getting_started.html) guide for installation guide and docs.
