#!/usr/bin/env bash

# This function checks whether we have a given program on the system.
_have()
{
	# Completions for system administrator commands are installed as well in
	# case completion is attempted via `sudo command ...'.
	PATH=$PATH:/usr/sbin:/sbin:/usr/local/sbin type $1 &> /dev/null
}

_have cerberus &&
_cerberus_complete()
{
	# Internal field separator
	local IFS=$'\t\n'

	# current and previously typed words
	local cur prev
	cur=${COMP_WORDS[COMP_CWORD]}
	prev=${COMP_WORDS[COMP_CWORD-1]}

	# completion array
	COMPREPLY=( )

	# cerberus <command>
	if [[ ${COMP_CWORD} -eq 1 ]]; then
		COMPREPLY=( $(compgen -W "$(printf "file \nsecret \nsdb ")" -- ${cur}) )

	# cerberus <command> <command>
	elif [[ $COMP_CWORD -eq 2 ]]; then
		case "$prev" in
		"file")
			COMPREPLY=( $(compgen -W "$(printf "read \ndownload \nedit \ndelete \nupload \nlist ")" -- ${cur}) )
			;;
		"secret")
			COMPREPLY=( $(compgen -W "$(printf "read \nwrite \nedit \ndelete \nlist ")" -- ${cur}) )
			;;
		"sdb")
			COMPREPLY=( $(compgen -W "$(printf "create \ndelete \nupdate ")" -- ${cur}) )
			;;
		*)
			# no autocomplete
			;;
		esac

	# cerberus <command> <command> <path to sdb/file/secret>
	elif [[ ${COMP_CWORD} -eq 3 ]]; then
		# grab secondary command
		local command=${COMP_WORDS[COMP_CWORD-2]}

		# Internal field separator
		IFS=$'\t\n\ '

		# switch on secondary command
		case "$command" in
		"file"|"secret")
			# count number of slashes
			NUMSLASH=$(echo ${cur} | grep -o / | wc -l | xargs echo)

			# empty completion
			completion=""

			# switch on number of slashes
			case "${NUMSLASH}" in
			"0")
				# list categories
				completion=$(cerberus category list -q)
				;;
			"1")
				# list SDBs
				completion=$(cerberus sdb list -a -q -c "${cur}")
				;;
			*)
				# change IFS to newline only
				IFS=$'\n'

				# isolate path to sdb from what has been typed
				sdbpath="$(echo ${cur} | sed -e 's/\(.*\)\/.*/\1/')""/"

				# remove escaped spaces
				sdbpath="$(echo ${sdbpath} | sed -e 's/\\//g')"

				# list secrets or files
				list="$(cerberus ${command} list "${sdbpath}" -a -q)"

				# if list not empty
				if [[ ${#list} -gt 0 ]]; then
					# use printf to separate fields using IFS
					completion="$(printf ${list})"
				fi
				;;
			esac
			# construct COMPREPLY with completion
			COMPREPLY=( $(compgen -W "${completion}" -- ${cur}) )
			;;
		"sdb")
			# count number of slashes
			NUMSLASH=$(echo ${cur} | grep -o / | wc -l | xargs echo)

			# empty completion
			completion=""

			# no completion for create command
			if [[ ${prev} != "create" ]]; then
				# switch on number of slashes
				case "${NUMSLASH}" in
				"0")
					# list categories
					completion=$(cerberus category list -q)
					;;
				"1")
					# list SDBs
					completion=$(cerberus sdb list -a -q -c "${cur}")
					;;
				*)
					# no autocomplete
					;;
				esac
			fi
			# construct COMPREPLY with completion
			COMPREPLY=( $(compgen -W "${completion}" -- ${cur}) )
			;;
		*)
			# no autocomplete
			;;
		esac

		local escaped_single_quote="'\''"
		local i=0

		# correctly escape all autocompletion entries
		for entry in ${COMPREPLY[*]}
		do
			if [[ "${cur:0:1}" == "'" ]]
			then
				# started with single quote, escaping only other single quotes
				# [']bla'bla"bla\bla bla --> [']bla'\''bla"bla\bla bla
				COMPREPLY[$i]="${entry//\'/${escaped_single_quote}}"
			elif [[ "${cur:0:1}" == "\"" ]]
			then
				# started with double quote, escaping all double quotes and all backslashes
				# ["]bla'bla"bla\bla bla --> ["]bla'bla\"bla\\bla bla
				entry="${entry//\\/\\\\}"
				COMPREPLY[$i]="${entry//\"/\\\"}"
			else
				# no quotes in front, escaping _everything_
				# [ ]bla'bla"bla\bla bla --> [ ]bla\'bla\"bla\\bla\ bla
				entry="${entry//\\/\\\\}"
				entry="${entry//\'/\'}"
				entry="${entry//\"/\\\"}"
				COMPREPLY[$i]="${entry// /\\ }"
			fi
			(( i++ ))
		done

		# append spaces to all non-empty autocompletion entries that end with a slash
		for (( i=0; i<${#COMPREPLY[@]}; i++ ));
		do
			if [[ "${COMPREPLY[$i]}" != */ ]] && [[ "${COMPREPLY[$1]}" != "" ]]; then
				COMPREPLY[$i]+=" "
			fi
		done
	fi
  return 0
} &&
complete -o nospace -o default -F _cerberus_complete cerberus
