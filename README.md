sshackd trades a valid Slack user token for a SSH certificate.

_If you have access to a LDAP service use that, don't use this._

To get started you'll need to setup a Slack app for your specific Slack team.

The various `ENV` variables you'll need are:

- `CLIENT_ID`: The Slack app client id
- `CLIENT_SECRET`: The Slack app client secret
- `EMAIL_DOMAIN`: The domain used for users in your org. This is used to filter the list of users down.
- `PORT`: The port the app should run on.
- `TEAM_ID`: The Slack team id. This is used to filter the list of users down.
- `CAKEY`: The signing key
