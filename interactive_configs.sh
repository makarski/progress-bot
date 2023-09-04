#!/bin/bash

GREEN="\e[32m"
YELLOW="\e[93m"
NOCOLOR="\e[0m"
DIM_COLOR="\e[2m"

source .env 2>/dev/null
config_key=progress-bot-config.toml

if [ -f ${config_key} ]
then
    echo -e $(printf "${GREEN}Configuration file ${YELLOW}${config_key}${GREEN} already exists.
    To rerun the config remove the file${NOCOLOR}")
    exit 0;
fi

# ------

echo -e $(printf "${GREEN}Configuring Jira Access")

read -p  "$(printf "> ${YELLOW}Enter jira email":${NOCOLOR}) " jira_email
read -p  "$(printf "> ${YELLOW}Enter jira account_id":${NOCOLOR}) " jira_account_id
read -p  "$(printf "> ${YELLOW}Enter jira base_url":${NOCOLOR}) " jira_base_url
read -p  "$(printf "> ${YELLOW}Enter jira token":${NOCOLOR}) " jira_token
read -p  "$(printf "> ${YELLOW}Enter jira project name":${NOCOLOR}) " jira_project_name

echo $'\n'
# ------

echo -e $(printf "${GREEN}Configuring Slack Access")

read -p  "$(printf "> ${YELLOW}Enter slack bot token":${NOCOLOR}) " slack_token
read -p  "$(printf "> ${YELLOW}Enter slack channel":${NOCOLOR}) " slack_channel

echo $'\n'
# ------

echo -e $(printf "${GREEN}Configuring OpenAI Access")

read -p  "$(printf "> ${YELLOW}Enter openai token":${NOCOLOR}) " openai_token
read -p  "$(printf "> ${YELLOW}Enter openai gpt_model(ex: gpt-3.5-turbo, gpt-4)":${NOCOLOR}) " openai_gpt_model
echo $'\n'

PB_JIRA_EMAIL=\"${jira_email}\" \
PB_JIRA_ACCOUNT_ID=\"${jira_account_id}\" \
PB_JIRA_BASE_URL=\"${jira_base_url}\" \
PB_JIRA_TOKEN=\"${jira_token}\" \
PB_JIRA_PROJECT=\"${jira_project_name}\" \
PB_SLACK_TOKEN=\"${slack_token}\" \
PB_SLACK_CHANNEL=\"${slack_channel}\" \
PB_OPENAI_TOKEN=\"${openai_token}\" \
PB_OPENAI_MODEL=\"${openai_gpt_model}\" \
envsubst < progress-bot-config.toml.template > ${config_key}

echo -e "$(printf "> ${GREEN}Successfully generated config file ${config_key}":${NOCOLOR})"
