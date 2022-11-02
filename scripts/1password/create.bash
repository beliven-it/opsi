#!/bin/bash

PROJECT="$1"
ADDRESS="$2"
PRI="${PROJECT} - PRI"
PUB="${PROJECT} - PUB"

check_me() {
  ME=$(op account list | grep ${ADDRESS} | awk '{print $3}')
  if [ -z ${ME} ]; then
    op account add --address "${ADDRESS}" --signin
    ME=$(op account list | grep ${ADDRESS} | awk '{print $3}')
    eval $(op account add --address "${ADDRESS}" --email "${ME}" --signin)

  else
    eval $(op signin --account "${ADDRESS}")
  fi
}

create_group() {
  match=$(op group list | grep "$1")

  if [ -n "$match" ]; then
    echo "Group $1 already exists"
    return 1
  fi

  command=$(op group create "$1" && op group user revoke --user "${ME}" --group "$1")

  if [ "$?" -eq 0 ];then
    echo "Group $1 created"
  else
    echo "Error on creating Group $1"
    echo "$command"
    exit 1
  fi

  
}

create_vault() {
  match=$(op vault list | grep "$1")

  if [ -n "$match" ]; then
    echo "Vault $1 already exists"
    return 1
  fi

  command=$(op vault create "$1" && op vault user revoke --user "${ME}" --vault "$1")
  if [ "$?" -eq 0 ];then
    echo "Vault $1 created"
  else
    echo "Error on creating Vault $1"
    echo "$command"
    exit 1
  fi
}

set_permissions() {
  VAULT="$1"
  GROUP="$2"
  PERMISSIONS="$3"

  if [[ ${PERMISSIONS} == "full" ]]; then
    PERMISSIONS="view_items,create_items,edit_items,archive_items,delete_items,view_and_copy_passwords,view_item_history,import_items,export_items,copy_and_share_items,print_items,manage_vault"
  fi

  command=$(op vault group grant --vault "${VAULT}" --group "${GROUP}" --permissions "${PERMISSIONS}")

  if [ "$?" -eq 0 ];then
    echo "Permissions set on Vault $1 to Group $2 with permissions $3"
  else
    echo "Error on creating permissions on Vault $1 to Group $2"
    echo "$command"
    exit 1
  fi
}

pri() {
  create_group "${PRI}"
  create_vault "${PRI}"

  set_permissions "${PRI}" "Owners" "full"
  set_permissions "${PRI}" "Administrators" "full"

  # Add vault PRI to group PRI with the right permissions
  set_permissions "${PRI}" "${PRI}" "view_items,create_items,edit_items,archive_items,view_and_copy_passwords,view_item_history,import_items,export_items,copy_and_share_items"
}

pub() {
  create_group "${PUB}"
  create_vault "${PUB}"

  set_permissions "${PUB}" "Owners" "full"
  set_permissions "${PUB}" "Administrators" "full"

  # Add vault PUB to group PUB with the right permissions
  set_permissions "${PUB}" "${PUB}" "view_items,view_and_copy_passwords,view_item_history,copy_and_share_items"
}

check_me

pri
pub

# Add vault PUB to group PRI with the right permissions
set_permissions "${PUB}" "${PRI}" "view_items,create_items,edit_items,archive_items,view_and_copy_passwords,view_item_history,import_items,export_items,copy_and_share_items"