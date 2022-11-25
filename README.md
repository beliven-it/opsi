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
postmark:
  api_url: "https://api.postmarkapp.com"
  token: <POSTMARK_TOKEN>
  slack_webhook: "<POSTMARK_SLACK_WEBHOOK>"
onepassword:
  address: "<ONEPASSWORD_ADDRESS>"
```

Some notes about these settings:

- `GITLAB_TOKEN` is an access token. You can generate in your gitlab settings [here](https://gitlab.com/-/profile/personal_access_tokens). Make sure to select the `api` scope in order to work.
- `POSTMARK_SLACK_WEBHOOK` the webhook to use for notifications.
- `ONEPASSWORD_ADDRESS` the 1password address of your tenant. Like: `my-tenant.1password.com`

<br><br><br><br><br><br>
<br><br><br><br><br><br>

## Commands

Below the list of available commands.

<br><br><br>

### Gitlab

There are these commands:

<br>

`opsi gitlab create projects <PROJECT_NAME> -s <SUB_GROUP_ID>`

> This command generate a project with the name specified under the subgroup provided.
>
> Over all it create the three branches:

- `Develop`
- `Staging`
- `Master` as default

**Flags**

- `s[subgroup]` !!**Required**!! for provide a subgroup for the project.
- `p[path]` for provide a custom slugify version of the name. Ex. `my-slug`. If not provided the system slugify the argument automatically.

<br>

---

<br>

`opsi gitlab create subgroup <SUBGROUP_NAME> -s <PARENT_ID> -p <PATH_NAME>`

> This command generate a subgroup with the name specified under the parent ID provided.

**Flags**

- `s[parent]` for provide a parent of the subgroup. Otherwise it will take this value from `.ospi.yml` file.
- `p[path]` for provide a custom slugify version of the name. Ex. `my-slug`. If not provided the system slugify the argument automatically.

<br>

---

<br>

`opsi gitlab create env <PROJECT_ID> <ENV_PATH> -e <ENV_NAME>`

> This command create environments variables for the specific project.

**Args**

- `<PROJECT_ID>` The ID of the project.
- `<ENV_PATH>` The path of the env file to upload.

**Flags**

- `e[env]` the env file. For example `staging`, `production`. If not provided it use the "*" value

<br>

---

<br>

`opsi gitlab delete env <PROJECT_ID> -e <ENV_NAME>`

> This command delete environments variables for the specific project.

**Args**

- `<PROJECT_ID>` The ID of the project.

**Flags**

- `e[env]` the env file. For example `staging`, `production`. If not provided it use the "*" value

<br>

---

<br>

`opsi gitlab bulk settings`

> This command is a bulk massive fix for old repos's branches with the actual standards.

<br>

---

<br>

`opsi gitlab deprovisioning <USERNAME>`

> This command is a deprovisioning an user from our gitlab workspace. It accept the USERNAME of the user.

<br><br><br><br><br><br>

### 1Password

There are these commands:

<br>

`opsi 1password create <PROJECT_NAME>`

> This command generate `private` and `public` VAULTS by project name

<br>

---

<br>

`opsi 1password deprovisioning <USER_EMAIL>`

> This command remove a specific user by email. If `USER_EMAIL` is not 
> provided, the command remove all users marked as inactive.

<br><br><br><br><br><br>

### Postmark

There are these commands:

<br>

`opsi postmark list servers`

> This command show a list of servers created in postmark.

<br>

The output follow this schema:

`[<INDEX>] <SERVER_NAME> (<SERVER_ID>)`

**Output**

```bash
[0] Server-test (123456)
[1] Server-test-02 (123457)
...
```

<br>

---

<br>

`opsi postmark create server <SERVER_NAME> -c <COLOR>`

> This command create a new server in postmark.

<br>

**Flags**

- `c[color]` for provide a color to assign to server. The default color is `blue` if not provided.

<br>

---

<br>

`opsi postmark edit server <SERVER_ID>`

> This command edit a new server in postmark.
>
> The edit action just update the Slack webhook provided in the configuration.

<br>

---

<br>

`opsi postmark bulk servers`

> This command edit all servers in postmark.
>
> The edit action just update the Slack webhook provided in the configuration.


<br><br><br><br><br><br>

### Hosts

There are these commands:

<br>

`opsi hosts check-reboot`

> This command show a list of servers needed to reboot.
>
> This command make use of `hssh`. You can find the tool and the istructions to install [here](https://github.com/beliven-it/hssh).

<br>




