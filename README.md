<br>
<p align="center"><img width="400" src="./assets/logo.svg" /></p>
<br>
<p align="center">
<img src="https://img.shields.io/github/go-mod/go-version/beliven-it/opsi?color=e75a39&style=for-the-badge" />
<img src="https://img.shields.io/github/v/release/beliven-it/opsi?color=e75a39&style=for-the-badge" />
<img src="https://img.shields.io/github/license/beliven-it/opsi?color=e75a39&style=for-the-badge" />
</p>
<p align="center">
<img src="https://img.shields.io/github/issues-pr/beliven-it/opsi?color=e75a39&style=for-the-badge" />
<img src="https://img.shields.io/github/issues/beliven-it/opsi?color=e75a39&style=for-the-badge" />
<img src="https://img.shields.io/github/contributors/beliven-it/opsi?color=e75a39&style=for-the-badge" />
</p>

<br><br>
<br><br>

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
  api_url: "https://company.gitlab.com/api/v4"
  token: "<GITLAB_TOKEN>"
  mirror:
    api_url: "https://gitlab.com/api/v4"
    group_id: "<GITLAB_MIRROR_GROUP_ID>"
    group_path: "gitlab.com/<GROUP_NAME>"  
    username: "<GITLAB_MIRROR_USERNAME>"
    token: "<GITLAB_MIRROR_TOKEN>"
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




