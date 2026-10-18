---
outline: deep
---

# Architecture

Nutrix is a web application that consists of the following three main subsystems:

#### Backend
The core system and backend of nutrix, a REST api that's build in [golang](go.dev), for more information check the dedicated git repo [nutrixpos/pos](https://github.com/nutrixpos/pos), pos uses a document based database built with [mongo](https://www.mongodb.com/)

#### Frontend
The user interface of nutrix, a web application that's build in [vue](https://vuejs.org/), the code is included in the [pos/frontend](https://github.com/nutrixpos/pos/tree/main/frontend) subdirectory.