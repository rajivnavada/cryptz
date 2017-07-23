cryptz
======

`cryptz` aims to help with storing and sharing application credentials securely. This repo contains code for a client to communicate with [cryptzd][cryptzd].

In order to use the client, you will need to provide the fingerprint of an activated public key. The command takes the form:

`cryptz -fpr $KEY_FINGERPRINT`

Available commands
------------------

* `projects list` - lists projects accessible by the user
* `project create $PROJECT_NAME $PROJECT_ENVIRONMENT` - creates a project with a name and a distinguishing environment
* `project $PROJECT_ID list credentials` - lists credentials belonging to project with given ID
* `project $PROJECT_ID add member $MEMBER_EMAIL` - adds a member (activated user) to the project with given ID
* `project $PROJECT_ID remove member $MEMBER_ID` - deletes member identified by given ID
* `project $PROJECT_ID add credential $CREDENTIAL_KEY $CREDENTIAL_VALUE` - set credential (key/value pair) attached to project with given ID
* `project $PROJECT_ID get credential $CREDENTIAL_KEY` - get value of credential. The returned credential is encrypted to the key with fingerprint specified when the command is run. However, cryptz will try to decrypt the credential on the client side before display
* `project $PROJECT_ID delete credential $CREDENTIAL_KEY` - deletes credential identified by given key

[cryptzd]: https://github.com/rajivnavada/cryptzd

