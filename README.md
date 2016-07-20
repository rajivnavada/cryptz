cryptz-client
=============

A client to help communicate with cryptz server.

Anatomy of a cryptz credential
------------------------------

* At the top level we'll have the concept of a "project"
* A project also contains a further specialization for "environment". (QA / Staging / Production etc...)
* Projects will have various team "members".
* Projects will have various "credentials".
* A credential is simply a key/value pair. The key is the identifier. Ex: storm-postgres-db. The value is the encrypted text. There will be one credential record per team member per project.
* At this point, the project table will contain a column for "adminId" which will be the userId of the project admin. For now, only one admin per project.
* Admins have the right to add / remove team members.
* Admins can also add / update / remove credentials.


HTTP routes
-----------

* GET /projects (returns a list of projects. name => ID pairs)
* GET /projects/:projectId/credentials (returns a list of available credentials for the project. list of keys => ID pairs)
* GET /projects/:projectId/credentials/:credential (returns encrypted text for requesting user)

TODO
----

Admin routes can be tackled a little later.
