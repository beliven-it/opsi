# OPSI

> All-in-one CLI for Beliven Ops daily usage!

The aim of this tool is to group all the daily and most used commands for OPS activities into a single tool.

<br>

This tool interact with some entities:

- `GitLab` for create, setup the repos.
- `1Password` for create groups and templates.
- `Postmark` for interact with the servers linked.
- `Hosts`for interact with ours hosts

<br><br><br><br><br><br>

## Configuration

Opsi automatically create the configuration file if not exist
in the system in `~/.config/opsi/config.yml`


The configuration file follow this schema:

```yml
gitlab:
  api_url: "https://gitlab.com/api/v4"
  token: "<GITLAB_TOKEN>"
  group_id: <GROUP_ID>
onepassword:
  address: "<ONEPASSWORD_ADDRESS>"
```

Some notes about these settings:

- `GITLAB_TOKEN` is an access token. You can generate in your gitlab settings [here](https://gitlab.com/-/profile/personal_access_tokens). Make sure to select the `api` scope in order to work.
- `ONEPASSWORD_ADDRESS` the 1password address of your tenant. Like: `my-tenant.1password.com`

<br><br><br><br><br><br>
<br><br><br><br><br><br>

## License

Licensed under [MIT](./LICENSE)




