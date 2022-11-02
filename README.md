# OPSI

> A all-in-one CLI for daily usage!!

The aim of this tool is to group all the daily and most used commands for OPS activities into a single tool. 

This tool interact with some entities:

- `Gitlab` for create, setup the repos.
- `1Password` for create groups and templates.
- `Postmark` for interact with the servers linked.

<br><br><br><br><br><br>

## Configuration

You must start by creating the configuration file for opsi:

```bash
# Go to home
cd

# Create the configuration file
touch .opsi.yml

```

Then copy inside the file the following configuration schema:

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

## Commands

Below the list of available commands.

<br><br><br><br><br><br>

### Gitlab

There are this commands:

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

<br><br>

`opsi gitlab create subgroup <SUBGROUP_NAME> -s <PARENT_ID> -p <PATH_NAME>`

> This command generate a subgroup with the name specified under the parent ID provided.

**Flags**

- `s[parent]` for provide a parent of the subgroup. Otherwise it will take this value from `.ospi.yml` file.
- `p[path]` for provide a custom slugify version of the name. Ex. `my-slug`. If not provided the system slugify the argument automatically.

<br><br>

`opsi gitlab bulk settings`

> This command is a bulk massive fix for old repos's branches with the actual standards.

<br><br>

`opsi gitlab deprovisioning <USERNAME>`

> This command is a deprovisioning an user from our gitlab workspace. It accept the USERNAME of the user.

<br><br><br><br><br><br>

### 1Password

There are this commands:

<br>

`opsi 1password create <PROJECT_NAME>`

> This command generate `private` and `public` VAULTS by project name

<br><br><br><br><br><br>

### Postmark

There are this commands:

<br>

`opsi postmark list servers`

> This command show a list of servers created in postmark.

<br>

The output follow this schema:

`[<TIMESTAMP>] <SERVER_NAME> (<SERVER_ID>)`

**Output**

```bash
[2022-02-12] Server-test (123456)
[2022-02-12] Server-test-02 (123457)
...
```




