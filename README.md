cryptz-client
=============

A client to help communicate with cryptz server. You will need to provide the fingerprint of an activated public key to successfully start the client. The command takes the form:

`cryptz -fpr $KEY_FINGERPRINT`

Available commands
------------------

* `project create $PROJECT_NAME $PROJECT_ENVIRONMENT` - creates a project with a name and a distinguishing environment
* `project list-credentials $PROJECT_ID` - lists credentials belonging to project with given ID
* `member add $PROJECT_ID $MEMBER_EMAIL` - adds a member (activated user) to the project with given ID
* `member delete $MEMBER_ID` - deletes member identified by given ID
* `credential set $PROJECT_ID $CREDENTIAL_KEY $CREDENTIAL_VALUE` - set credential (key/value pair) attached to project with given ID
* `credential get $PROJECT_ID $CREDENTIAL_KEY` - get value of credential. The returned credential is encrypted to the key with fingerprint specified when the command is run.

