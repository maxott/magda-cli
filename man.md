# NAME

magda-cli

# SYNOPSIS

  - **magda-cli \[\<flags\>\] \<command\> \[\<args\> ...\]**  

# DESCRIPTION

Managing records & schemas in Magda.

# OPTIONS

  - **-h, --help**  
    Show context-sensitive help (also try --help-long and --help-man).

  - **-H, --host=HOST**  
    DNS name/IP of Magda host \[MAGDA\_HOST\]

  - **--tenantID=TENANTID**  
    Tenant ID \[TENANT\_ID\]

  - **--authID=AUTHID**  
    Authorization Key ID \[AUTH\_ID\]

  - **--authKey=AUTHKEY**  
    Authorization Key \[AUTH\_KEY\]

  - **--useTLS**  
    Use https

  - **--skipGateway**  
    Skip gateway server and call registry server directly
    \[SKIP\_GATEWAY\]\]

# COMMANDS

## **help \[\<command\>...\]**

Show help.

## **record list \[\<flags\>\]**

List some records

  - **-a, --aspects=ASPECTS**  
    The aspects for which to retrieve data

  - **-q, --query=QUERY**  
    Record Name

  - **-o, --offset=OFFSET**  
    Index of first record retrieved

  - **-l, --limit=LIMIT**  
    The maximumm number of records to retrieve

  - **-t, --pageToken=PAGETOKEN**  
    Token that identifies the start of a page of results

## **record read --id=ID \[\<flags\>\]**

Read the content of a record

  - **-i, --id=ID**  
    Record ID

  - **--add-aspects=ADD-ASPECTS**  
    Add aspects to record listing (comma separated)

  - **-a, --aspect=ASPECT**  
    Show only this aspect of the record as result

## **record create --name=NAME \[\<flags\>\]**

Creates a new record

  - **-i, --id=ID**  
    Record ID (defaults to UUID)

  - **-n, --name=NAME**  
    Record Name

  - **-a, --aspectName=ASPECTNAME**  
    Name of aspect to add (requires --aspectFile)

  - **-f, --aspectFile=ASPECTFILE**  
    File containing aspect data

## **record update --id=ID \[\<flags\>\]**

Update an existing record

  - **-i, --id=ID**  
    Record ID (defaults to UUID)

  - **-n, --name=NAME**  
    Record Name

  - **-a, --aspectName=ASPECTNAME**  
    Name of aspect to add (requires --aspectFile)

  - **-f, --aspectFile=ASPECTFILE**  
    File containing aspect data

## **record delete --id=ID \[\<flags\>\]**

Delete a record or one of it's aspects

  - **-i, --id=ID**  
    Record ID

  - **-a, --aspect=ASPECT**  
    Only delete this aspect

## **record history --id=ID \[\<flags\>\]**

Get a list of all events for a record

  - **-i, --id=ID**  
    Record ID

  - **-e, --event-id=EVENT-ID**  
    Only show event wiht event-id

  - **-o, --offset=OFFSET**  
    Index of first record retrieved

  - **-l, --limit=LIMIT**  
    The maximumm number of records to retrieve

  - **-t, --pageToken=PAGETOKEN**  
    Token that identifies the start of a page of results

## **schema list**

List all aspect schemas

## **schema create --name=NAME --id=ID --schemaFile=SCHEMAFILE**

Creates a new schema

  - **-n, --name=NAME**  
    Descriptive name

  - **-i, --id=ID**  
    Schema ID

  - **-f, --schemaFile=SCHEMAFILE**  
    File containing schema/aspect decalration

## **schema read --id=ID**

Read the content of a record

  - **-i, --id=ID**  
    Record ID

## **schema update --id=ID --schemaFile=SCHEMAFILE \[\<flags\>\]**

Update existing schema

  - **-n, --name=NAME**  
    Descriptive name

  - **-i, --id=ID**  
    Schema ID

  - **-f, --schemaFile=SCHEMAFILE**  
    File containing schema/aspect decalration
