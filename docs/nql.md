NQL syntax overview
Specifying table
Every NQL query starts with a short statement specifying the table to select data from. The syntax to specify the table is:

<namespace>.<table>

For example, listing all records in the events table from the execution namespace translates into the following statement:


Copy
execution.events
Syntax shortcuts
Instead of typing the namespace and the table, you can also use the predefined shortcuts. Type the table name only, without a namespace first to retrieve data from the following tables:

Namespace
Table
Shortcut
application

applications

applications

binary

binaries

binaries

campaign

campaigns

campaigns

device

devices

devices

user

user

users

For example, type devices instead of device.devices to list all the records within the devices table in the device namespace.


Copy
devices
You do not need to specify the table fields included in the results to query data from the table. The system includes default fields that are most relevant to identify the records. For more information about fields contained in specific table, refer to the NQL data model page. Use the NQL list keyword to access other fields in the specific table.

Specifying time frame
You have the option to filter your results over a specific period of time by putting a time frame selection right after the table name in your NQL statement. Depending on what you need, you can choose from various data selection formats and time precisions. For example you can specify the number of days back:


Copy
execution.crashes during past 7d
Or specific date:


Copy
execution.crashes on Feb 8, 2024
You can also use a time selection when querying the following inventory objects: devices, users, binaries. If you specify the time frame for the inventory objects, the system refers to the events behind the object's activity.

For example, the following queries refer to the same set of data.


Copy
devices during past 7d

Copy
devices 
| with device_performance.events during past 7d
For more information regarding the time selection formats refer to the NQL time selection

Customizing Query Results
After specifying the table and timeframe, you can further refine your query by providing additional instructions to the system using keywords, operators and functions. These refinements allow you to organize, filter or aggregate your results to gather more comprehensive insights.

For example:

Filter the results using the where clause


Copy
binaries during past 24h
| where binary.name == "dllhost.exe"
Select specific data to display using the list clause


Copy
binaries during past 24h
| where binary.name == "dllhost.exe"
| list name, version, platform, architecture, size
Order results using the sort ... desc clause


Copy
binaries during past 24h
| where binary.name == "dllhost.exe"
| list name, version, platform, architecture, size
| sort size desc
Set a maximum number of results using the limit clause


Copy
binaries during past 24h
| where binary.name == "dllhost.exe"
| list name, version, platform, architecture, size
| sort size desc
| limit 10
For more information about specific instructions, refer to the NQL keywords section.

Pattern matching
Use wildcard characters such as * and ? for text filters.

* replaces any number of characters

? replaces any single character

For example, listing all binaries with a name starting with dll and finishing with .exe translates into the following query:


Copy
binaries during past 24h
| where binary.name == "dll*.exe"
| list size,name,version 
| sort size desc 
| limit 100
NQL time selection

NQL data types

Commenting
Use comments in your NQL queries to include explanatory notes that are ignored during execution. Comments help clarify the intent of the query, making it easier to read, maintain, and understand.

Use /* to begin a comment and */ to end it. All text between these symbols will be ignored during execution.


Copy
devices during past 7d 
/* This comment spans
multiple lines and is ignored */
| list device.name, device.entity
Keyboard shortcuts for commenting
Use the following keyboard shortcuts to quickly add or remove comments in your NQL queries.

Line-based comment
Toggle comment on the current line or selected multiple lines.

The following shortcuts add or remove comment markers (/* ... */) around the entire line where your cursor is placed. If multiple lines are highlighted, it wraps full lines in a single comment block.

Windows: Press Ctrl + /

macOS: Press Cmd + /

Inline or block comment
Toggle comment on selected code.

The following shortcuts add or remove comment markers (/* ... */) around the highlighted portion of code, even if the selection starts or ends mid-line.

Windows: Press Shift + Alt + A

macOS: Press Shift + Option + A

Pressing either shortcut again will remove the comment block it added.

Valid comment placement
Comments can be added in most parts of an NQL query. However, there are specific cases where comments are not allowed.

The following table outlines invalid comment placements. If a comment is added in one of these locations, the NQL editor will display an error.

Description of invalid comment placement
Example of invalid comment placement
Between | and statement keyword


Copy
devices | /* comment */ where name == "test"
Inside expressions


Copy
devices | where name /* comment */ == "test"
Between operator and operand


Copy
devices | compute x = count /* comment */ ()
Inside function calls


Copy
devices | summarize c = count() by /* comment */ 1d, start_time
Last updated 1 month ago

NQL time selection
In NQL you can specify the time frame in various formats.

Time selection formats
NQL during past
The during past clause allows you to filter your results by specifying a particular time period leading up to the present. The time can be expressed in minutes, hours, or days.

Examples:

Retrieving the number of navigations in the past 45 minutes.


Copy
web.page_views during past 45min
| summarize total_navigations = number_of_page_views.sum()
Retrieving the number of navigations in the past 1 hour.


Copy
web.Page_views during past 1h
| summarize total_navigations = number_of_page_views.sum() 
Retrieving the number of navigations in the past 12 hours.


Copy
web.Page_views during past 12h
| summarize total_navigations = number_of_page_views.sum() 
Retrieving the number of navigations in the past 3 days.


Copy
web.Page_views during past 3d
| summarize total_navigations = number_of_page_views.sum() 
NQL from to
The from to clause allows you to apply custom timeframe filters by specifying the start and end times for the desired period.

Specifying a fixed timeframe
Apply a timeframe filter by specifying fixed datetime values for the start and end of the period.

Examples:

The number of navigations from June 1, 2023 to June 15, 2023


Copy
web.page_views from Jun 1, 2023 to Jun 15, 2023
| summarize total_navigations = number_of_page_views.sum() 
The number of navigations from June 15, 2023 at 12:30 to June 15, 2023 at 16:15


Copy
web.page_views from Jun 15, 2023, 12:30 to Jun 15, 2023, 16:15
| summarize total_navigations = number_of_page_views.sum()
The number of navigations from 2023-02-01 00:00:00 to 2023-02-28 23:45:00


Copy
web.page_views from 2023-02-01 00:00:00 to 2023-02-28 23:45:00
| summarize total_navigations = number_of_page_views.sum()
The number of navigations from 2023-02-01 to 2023-02-28


Copy
web.page_views from 2023-02-01 to 2023-02-28
| summarize total_navigations = number_of_page_views.sum()
For more information about the allowed date formats, refer to the NQL data types section. Note that the autocomplete functionality in the NQL editor provides suggestions with available data formats.

Specifying a relative timeframe
Apply a timeframe filter by defining a time window relative to the current time, for example: 15m ago, 2h ago, 1d ago.

The time can be expressed in minutes, hours or days.

Examples:

The number of navigations from the previous day.


Copy
web.page_views from 1d ago to 1d ago
| summarize total_navigations = number_of_page_views.sum() 
The number of navigations grouped into 7-day intervals over a consecutive three-week period.


Copy
devices
| include web.page_views from 21d ago to 13d ago
| compute week1 = number_of_page_views.sum()
| include web.page_views from 14d ago to 8d ago
| compute week2 = number_of_page_views.sum()
| include web.page_views during past 7d
| compute current_week = number_of_page_views.sum()
NQL on
The on clause allows you to select a specific day when querying data.

Examples:

The number of navigations on July 15, 2023


Copy
web.page_views on Jul 15, 2023
| summarize total_navigations = number_of_page_views.sum() 

Copy
web.page_views on 2023-06-15
| summarize total_navigations = number_of_page_views.sum() 
NQL datetime functions
You can further customize your time selection with an NQL where conditions, to specify time windows relative to points in time other than the current time, for example, business hours, business week, or specific days of the month.

Refer to NQL datetime functionsfor more information.

Time granularity and retention
You have the flexibility to choose the precision level for time selection. Use minutes or hours in NQL time specification to retrieve more granular data. Use days to retrieve less granular data typically covering a longer time span.

When specifying timeframes at the day level, without a specific time (for example, during past 2d, June 1, 2023 or 1d ago), the system defaults to include:

Start of the period: 00:00:00 (midnight) on the start date.

End of the period: 23:59:59 on the end date.

This ensures that the entire day(s) within the specified range are included.

Data storage and granularity also depend on specific tables. Refer to the Data resolution and retention documentation page for more details.

Retrieving high-resolution data for Desktop Virtualization
By default, VDI event data are available with 5-minute or 1-day resolution, depending on time selection in your NQL query. To increase VDI data resolution, add by 30s at the end of the time selection in the query. High-resolution data is available for the past 2 days.

The examples below show how to retrieve high-resolution data from vdi_events.

Increasing data resolution from 1 day to 30 seconds.

Copy
session.vdi_events during past 1d by 30s

Copy
session.vdi_events on 2024-08-04 by 30s
Increasing data resolution from 5 minutes to 30 seconds.

Copy
session.vdi_events during past 24h by 30s

Copy
session.vdi_events from Nov 08, 2024, 15:15 to Nov 08, 2024, 17:30 by 30s
Timezones
When the Nexthink cloud instance is located in a different timezone from that of the user, the time selection units determine which timezone is considered for defining the beginning and end of the specified time period.

Full-day timeframes (for example, during past 2d, from 2024-02-07 to 2024-02-08, on Feb 8, 2024) use the cloud instance timezone.

Timeframes expressed in hours and minutes (for example, during past 15min, from 2024-02-07 14:45:00 to 2024-02-08 14:45:00) use the user timezone.

This distinction applies solely to the time period covered in the query. The results will always be displayed in the timezone of the user.

Example:
Let's consider how this would work in a real-world scenario.

Suppose two Nexthink users query the data using the Nexthink platform set to Eastern Time (ET).

The first user operates in the same timezone as the Nexthink platform. The current time for them is November 11, 05:26:15.

The second user operates in the Central European Time (CET) zone. The current time for them is November 11, 11:26:15.

In such a case, time-related queries made by the second Nexthink user will be translated into the corresponding timeframes, considering the timezone differences between CET and ET. This ensures accurate data retrieval and analytics, regardless of geographical location or timezone.

Timeframe selection
Nexthink user in the Eastern Time (ET) zone - the same zone as the Nexthink platform:
Nexthink user in the Central European Time (CET) zone:
past 15min

Nov 11, 05:15:00 AM
– 05:30:00 AM ET​

Nov 11, 11:15:00 AM
– 11:30:00 AM CET

past 2h

Nov 11, 04:00:00 AM
– 06:00:00 AM ET​

Nov 11, 10:00:00 AM
– 12:00:00 PM CET

past 24h

Nov 10, 06:00:00 AM
– Nov 11 06:00:00 AM ET​

Nov 10, 12:00:00 PM
– Nov 11, 12:00:00 PM CET

from 2021-11-11 00:00:00
to 2021-11-11 12:00:00

from 2021-11-11 12:00:00 AM
to 2021-11-11 12:00:00 PM

from 2021-11-11 12:00:00 AM
to 2021-11-11 12:00:00 PM

past 1d

Nov 11, 12:00:00 AM
– Nov 12, 12:00:00 AM ET

Nov 11, 06:00:00
– Nov 12, 06:00:00 CET

on Nov 10, 2021

Nov 10, 12:00:00 AM
– Nov 11, 12:00:00 AM ET​

Nov 10, 06:00:00
– Nov 11, 06:00:00 CET

Last updated 12 months ago

NQL data types
The data type is an attribute of the value stored in a field. It dictates what type of data a field can store.

When applying conditions to the NQL query using a where clause, only values of the same data types can be compared which is reflected in the format of the value.

For example, in the following query:

The first where clause compares values of the string data type. Consequently, the comparison value is enclosed in quotes to denote its string nature.

The second where clause compares versions. Here, the comparison value is prefixed with 'v' and includes multiple points to represent a version number.

The last where clause compares integers. In this case, the comparison value is expressed solely as a standalone number without any additional characters.


Copy
devices during past 1d
| include execution.crashes during past 1d
| where application.name == "Microsoft 365: Teams"
| where binary.version == v1.7.0.1864
| compute number_of_crashes_ = number_of_crashes.sum()
| where number_of_crashes_ >= 3
The following data types are present in the NQL data model:

Data type
Valid operators
Definition
Value example
string

== or =

!=

in

!in

a string of text characters

"abc" or 'abc'

int

=

!=

<

>

<=

>=

in

!in

a whole number

10

float

=

!=

<

>

<=

>=

a floating point number

10.1

Boolean

=

!=

a true or false value

true

false

date time

=

!=

<=

>=

a date with a time

2024-07-15 10:15:00

enumeration

=

!=

sets of named things

for example red blue white

status == red

byte

<

>

<=

>=

a number of bytes

(an int with a unit)

100B

200KB

3MB

12GB

2TB

duration

=

!=

<

>

<=

>=

a duration in time

(an int with a unit)

5ms

10s

4min

3h

2d

IP address

=

!=

IPv4 or IPv6 addresses

with optional mask

123.123.0.0

123.123.0.0/24

f164:b28c:84a5:9dd3:ef21:8c9d:d3ef:218c

f164:b28c:84a5:9dd3::/32

version

<

>

<=

>=

==

!=

a set of numbers separated by a .

v12.212

v1.2.5.9

v13.5.10

v2022.6

v1.2.4125

v6.8.9.7.6.5.4.3

string array

contains

!contains

an array of strings

for example ['abc', 'def', 'xyz']

tags contains "abc"

tags !contains "*xyz"

Last updated 1 year ago

NQL keywords
An NQL keyword is a reserved term used to construct a statement or a clause. To use a keyword, put a | (pipe) symbol before it and provide instructions in accordance with NQL syntax specific for this keyword.

NQL list

NQL limit

NQL where

NQL sort

NQL with

NQL include

NQL compute

NQL summarize

NQL summarize by

Last updated 1 month ago

NQL list
A list clause allows you to specify which fields you want to select.

Syntax

Copy
...
| list <field_1>, <field_2>, <metric_1>, <metric_2> ...
Example
Select the name and the type from the users table.


Copy
users during past 7d
| list username, type
Name
Type
Hemi Charmian

LOCAL_USER

Wangchuk Mirjam

LOCAL_ADMIN

NQL limit
A limit clause restricts the maximum number of rows returned.

Syntax

Copy
...
| limit <number of rows>
Example
Select a maximum of 15 users.


Copy
users during past 7d
| limit 15

NQL where
A where clause allows you to add conditions to your query to filter the results using NQL comparison operators and NQL logical operators.

Comparing field value to a fixed reference
Compare field value to a fixed reference to filter results that match a specific, unchanging criterion. For example:

Filter devices with a specific operating system.

Filter devices with free memory below a specified threshold.

Filter specific binary versions.

Syntax

Copy
...
| where <field name> <comparison operator> <static value>
Examples
Select the devices running the Windows operating system.


Copy
devices during past 7d
| where operating_system.platform == Windows
Name
Platform
nxt-gcarlisa

Windows

nxt-wmirjam

Windows

Select the devices not running the Windows operating system.


Copy
devices during past 7d
| where operating_system.platform != Windows
| list name, operating_system.platform
Name
Platform
nxt-jdoe

macOS

nxt-vlatona

macOS

Select the users whose name contains “jo”.


Copy
users during past 7d
| where username == "*jo*"
Name
John Fisher

John Doe

Comparing two field values against each other
Compare two field values against each other when you wish to filter results based on a dynamic relationship between fields. Only fields from the same table can be compared against each other.

You can compare the following fields:

native fields

context fields

metrics (aliases) computed in the query

manual custom fields

Syntax

Copy
...
| where <field-a name> <comparison operator> <field-b name>
Examples
Comparing native fields
Identify users which don't use the same peripheral for both the speaker and the microphone.


Copy
users
| with collaboration.sessions
| where participant_device.microphone != participant_device.speaker
Comparing a native field with a context field
Filter out events where the device has changed location


Copy
connection.events during past 7d
| where destination.country == context.location.country
Comparing native field to computed metric
Identify devices which have not had any Collector activity after an execution crash.


Copy
devices during past 7d
| include execution.crashes during past 7d
| compute last_crash_time = time.last()
| where last_crash_time > last_seen
Comparing native field to a manual custom field
Compare the package version to a required compliant version that is stored in a manual custom field.


Copy
packages 
| where package.version == package.#required_version
Using multiple conditions
Use multiple filters separated by NQL bitwise operators(and or or) to apply more complex conditions. The conditions in the filter are grouped together to preserve the order of precedence. When you put where clauses on separate lines, the result is the same as if you created one where clause with multiple and conditions.

The following queries provide the exact same results.


Copy
devices during past 7d
| where device.entity == "Lausanne" and device.hardware.type == laptop

Copy
devices during past 7d
| where device.entity == "Lausanne" 
| where device.hardware.type == laptop
Last updated 7 months ago

NQL sort
A sort ... asc or sort ... desc clause orders the results by a field in ascending or descending order, respectively.

Syntax
Sort data starting from the lowest value:


Copy
....
| sort <field name> asc
Sort data starting from the highest value:


Copy
....
| sort <field name> desc
Examples
Sort users by their username in alphabetical order:


Copy
users during past 7d
| list username, type
| sort username asc
Name
Type
Alice Smith

LOCAL_USER

Amanda Carella

LOCAL_ADMIN

Sort users by their username in reverse alphabetical order:


Copy
users during past 7d
| list username, type
| sort username desc
Name
Type
Zion Bush

LOCAL_USER

Zachary Doe

LOCAL_ADMIN

Last updated 1 year ago

NQL with
A with clause allows you to join an inventory object table with an event table. It returns data per object only when there is at least one event recorded for a specific object. Use it to query inventory objects with conditions on events.

Syntax

Copy
<object table> ...
| with <event table> ...
Example
Select all the devices with at least one error during the last seven days.


Copy
devices
| with web.errors during past 7d
| list device.name, operating_system.name
Name
OS name
device-54304276

Windows 10 Pro 21H1 (64 bits)

device-c0b53b3f

Windows 10 Enterprise 21H1 (64 bits)

device-71cedc8f

Windows 10 Enterprise 21H1 (64 bits)

device-dc98cd15

Windows 10 Enterprise 21H1 (64 bits)

device-b5d55bd0

Windows 10 Pro 21H1 (64 bits)

device-706d3c09

Windows 10 Pro 21H1 (64 bits)

device-a56b63f1

Windows 10 Enterprise 21H1 (64 bits)

device-259c7017

Windows 10 Pro 20H2 (64 bits)

device-d0ce2109

Windows 10 Enterprise 21H1 (64 bits)

Computing new metric
The with clause can be used along with a compute clause that appends the object table with a new column with metric per object. Refer to the NQL compute keyword documentation page for more information.

Using multiple ‘with’ clauses
An NQL query can contain multiple with clauses.


Copy
binary.binaries
| with execution.crashes during past 1d
| compute total_number_of_crashes = count()
| with execution.events during past 1d
| compute sum_of_freezes = number_of_freezes.sum()
| list total_number_of_crashes, sum_of_freezes, name
Number of crashes
Sum of freezes
Binary name
MD5 hash
7

0

odio.exe

f32bd724cb4b8593c9789ec584eb38dc

12

0

volutpat.exe

5ec62b81e594367fa20a3fbdf4e4e7f3

24

0

eget.exe

dc182b7939eba5ca8b1d64396b88fcd2

3

0

euismod.exe

2d0c540521f7e5683487c42c6ff52479

9

0

euismod.exe

2d0c540521f7e5683487c42c6ff52479

17

0

aliquet.exe

f4c4ad04db18ff1d225cbc43e864748a

Filtering data
Only the computed values are available outside of the with clause. When you start with devices, only the fields of that table are available for other statements. Adding a with and a compute makes new fields available.


Copy
devices
| with web.errors during past 7d
| compute total_errors_device = number_of_errors.sum()
| where total_errors_device > 10
| list device.name, total_errors_device
| sort total_errors_device desc
Name
total_errors_device
device-741da9be

125

device-c91fa737

120

device-08469fee

62

device-f2301dea

51

device-9e07abe9

45

device-03680882

42

device-25c67269

42

device-f8586bb6

41

device-b5d55bd0

39

device-60ea7a88

39

Last updated 1 year ago

NQL include
An include clause allows you to join an inventory object table with an event table. It returns data per object even when there is no event recorded for a specific object. Use it to make sure to take into account all objects when computing metrics.

Syntax

Copy
<object table> ...
| include <event table> ...
| compute <new metric name> = <metric>.<aggregation function>
...
Example
List the binaries that triggered an execution crash and the associated number of crashes, during the last 24 hours.


Copy
binaries
| include execution.crashes during past 24h
| compute total_number_of_crashes = count()
| list total_number_of_crashes, name
| sort total_number_of_crashes desc
Number of crashes
Binary name
83

lorem.exe

20

bibendum.exe

10

imperdiet.exe

9

tempor.exe

7

egestas.exe

6

semper.exe

6

justo.exe

Using multiple ‘include’ clauses
An NQL query can contain multiple include clauses , allowing you to join the same event table with different conditions or to join several different event tables.


Copy
binaries
| include execution.crashes during past 1d
| compute total_number_of_crashes = count()
| include execution.events during past 1d
| compute sum_of_freezes = number_of_freezes.sum()
| list total_number_of_crashes, sum_of_freezes, name
| sort total_number_of_crashes desc
Number of crashes
Sum of freezes
Binary name
MD5 hash
60

0

odio.exe

f32bd724cb4b8593c9789ec584eb38dc

26

0

volutpat.exe

5ec62b81e594367fa20a3fbdf4e4e7f3

12

0

eget.exe

dc182b7939eba5ca8b1d64396b88fcd2

7

0

euismod.exe

2d0c540521f7e5683487c42c6ff52479

7

0

euismod.exe

2d0c540521f7e5683487c42c6ff52479

6

0

aliquet.exe

f4c4ad04db18ff1d225cbc43e864748a

6

0

vitae.exe

bd85d77734d35c5ee00edeffc44e1dcd

Understanding the purpose of ‘with’ and ‘include’ clauses
The include and with keywords are very similar but have very different purposes.

Keyword
Meaning
Scope
Purpose
Compute
with

Retain only those objects which have an event recorded

Modifies the scope

Filter and/or compute values for objects with events

A value is always computed and added

include

Retain all objects, including those that do not have an event recorded

Without a compute statement, no effect on scope

Only useful when a value is computed for all objects

Objects without events have no computed value

Last updated 1 year ago

NQL compute
The compute command aggregates and extracts metrics from the events table and appends it to the results table as a new column with metric per object. It can be used only after a with or include clause.

Syntax

Copy
...
| include... 
| compute <new_metric> = <metric>.<aggregation function>

Copy
...
| with... 
| compute <new_metric> = <metric>.<aggregation function>
Example

Copy
devices during past 7d
| include execution.crashes during past 7d
| compute nb_crashes = number_of_crashes.sum()
Using with the ‘count()’ function
When used without a field specified, the count() aggregation function applies to the event table. For example, in the following query the compute clause appends new column with the number of boots per device.


Copy
devices during past 7d
| include device_performance.boots during past 7d
| compute nb_boots = count()
You can also count the unique inventory objects as a new column, using the <object>.count() syntax. It appends a new column with either 1 or 0 as the value, based on whether the object has relevant events or not. In the following example, the compute clause returns 1 for the devices that have been booted during past 7 days, and 0 for devices with no boots recorded in that time period. In the last statement, summarize clause is used for computing the ratio of devices with boots.


Copy
devices during past 7d
| include device_performance.boots during past 7d
| compute nb_devices_with_boots = device.count()
| summarize ratio_devices_with_boots = nb_devices_with_boots.sum()/count()
Last updated 1 year ago

NQL summarize
The summarize statement condenses the information into a single result.

Syntax

Copy
...
| summarize <new metric name> = <metric>.<aggregation function>
Examples
Compute the average value of the number of page views per device.


Copy
devices
| with web.page_views from Jun 1 to Jun 7
| compute num_navigations = number_of_page_views.sum()
| summarize average_num_navigation_per_device = num_navigations.avg()
average_num_navigation_per_device
115.9

Count the number of devices active in the last 7 days. In case of the count() aggregation function, you can omit the filed name before the aggregation to count the number of records of the root table.


Copy
devices during past 7d
| summarize number_of_devices = count()
number_of_devices
285

Compute the total size of all the binaries.


Copy
binaries during past 7d
| summarize total_size = size.sum()
total_size
611.3 GB

Last updated 1 year ago

NQL summarize by
The summarize by statement condenses the information into aggregated results grouped by properties or time interval.

Grouping by property
Enter the field name after by to create a breakdown by a property. Enter additional field names separated by a comma to create more breakdown dimensions.

The summarize by clause does not support grouping by properties with numeric data types such as days_since_last_seen (integer), last_seen (date time) or hardware.memory (byte).

Syntax

Copy
...
| summarize <new metric name> = <metric>.<aggregation function> by <field_1>, <field_2> ...
Example
Display the average Confluence backend page load time per device in the last 7 days.


Copy
web.page_views during past 7d
| where application.name == "Confluence"
| summarize backendTime = page_load_time.backend.avg() by device.name
| list device.name, backendTime
| sort backendTime desc
Device name
backendTime
device-10d267d2

508.2 ms

device-d1d5abc9

498.9 ms

device-5117c4c3

432.1 ms

device-16834449

431.9 ms

device-b634ce84

429.4 ms

device-731db075

349.8 ms

device-7fb313ef

293.9 ms

device-a834a720

277.6 ms

…

…

Grouping by period
The summarize by statement when used in combination with a time period, groups the metric values into time buckets.

Syntax

Copy
...
| summarize <new metric name> = <metric>.<aggregation function> by <time period>
Valid period values are:

15 min 30 min 45 min …
The value must be a multiple of 15.

1 h 2 h 3 h ...
The value must be a whole number.

1 d 2 d 3 d ...
The value must be a whole number.

Example
Display daily number of crashes in the last 7 days in chronological order.


Copy
execution.crashes during past 7d
| summarize total_number_of_crashes = count() by 1d
| sort start_time asc
start_time
end_time
bucket_duration
number_of_crashes
2021-03-05
00:00:00

2021-03-06
00:00:00

1 d

758

2021-03-06
00:00:00

2021-03-07
00:00:00

1 d

700

2021-03-07
00:00:00

2021-03-08
00:00:00

1 d

954

2021-03-08
00:00:00

2021-03-09
00:00:00

1 d

493

2021-03-09
00:00:00

2021-03-10
00:00:00

1 d

344

2021-03-10
00:00:00

2021-03-11
00:00:00

1 d

765

2021-03-11
00:00:00

2021-03-12
00:00:00

1 d

857

Grouping by property and period
Combine properties and time period to generate time buckets with additional breakdowns. You can use multiple fields, but only one time period selector. The sequence of items is arbitrary; the time period selector can be positioned anywhere within the list of fields.

Syntax

Copy
...
| summarize <new metric name> = <metric>.<aggregation function> by <field_1>, <field_2>, ... <time period>, ...
Example
Display daily number of crashes in the last 30 days broken down by operating system platform and sorted starting from the highest number of crashes.


Copy
execution.crashes during past 30d
| summarize total_number_of_crashes = count() by 1d, device.operating_system.platform 
| sort total_number_of_crashes desc
Device platform
start_time
end_time
bucket_duration
number_of_crashes
Windows

2021-12-07
00:00:00

2021-12-08
00:00:00

1 d

690

Windows

2021-12-08
00:00:00

2021-12-09
00:00:00

1 d

533

macOS

2021-12-20
00:00:00

2021-12-21
00:00:00

1 d

511

Windows

2021-12-17
00:00:00

2021-12-18
00:00:00

1 d

493

Windows

2021-12-08
00:00:00

2021-12-09
00:00:00

1d

356

macOS

2021-12-20
00:00:00

2021-12-21
00:00:00

1d

325

…

…

…

…

…

Last updated 1 year ago

NQL operators
NQL syntax operators

NQL arithmetic operators

NQL comparison operators

NQL logical operators

NQL bitwise operators

Last updated 2 months ago

NQL syntax operators
Operator
Definition
Examples
=

Alias operator, only supported in a compute and summarize statements.

num_of_devices = count()

|

Pipe operator, separator between statements.

devices | where name == "abc"

Last updated 2 months ago

NQL arithmetic operators
Use arithmetic operators to calculate field values inNQL compute, NQL summarize and NQL summarize by clauses.

Operator
Definition
Example
+

Addition

column_name.sum() + 10

-

Subtraction

count() - 10

/

Division

count() / 10

*

Multiplication

count() * 0.1

( )

Grouping of operators

(sum() / count()) / 100

Last updated 7 months ago

NQL comparison operators
Use comparison operators with NQL where clause to filter your NQL query results.

Operator
Definition
Supported data types
Examples
== or =

Equals

string

int

float

Boolean

date time

enumeration

duration

IP address

version

| where user.name = "jdoe@kanopy"

| where user.name == "jdoe@kanopy"

!=

Not equals

string

int

float

Boolean

date time

enumeration

duration

IP address

version

| where hardware_manufacturer != "VMWare"
| where hardware_manufacturer != null

>

Greater than

int

float

duration

byte

IP address

version

| where hardware.memory > 8GB

<

Less than

int

float

duration

byte

IP address

version

| where hardware.memory < 16GB

>=

Greater or equal

int

float

date time

duration

byte

IP address

version

| where hardware.memory >= 8GB

<=

Less or equal

int

float

date time

duration

byte

IP address

version

| where hardware.memory <= 16GB

Refer to NQL data types for more information.

In comparison operations, = and == are interchangeable. However, when used for aliasing in NQL compute or NQL summarize statements, only the single = is supported.

All expressions used in combination with these operators are case-insensitive. For example, the following queries return the same results:


Copy
devices during past 24h
| where name == "CORPSYS2022"

Copy
devices during past 24h
| where name == "CoRpSyS2022"
Using wildcards
Use wildcards to match partial values and increase filter flexibility. Expressions used in combination with comparison operators support the following wildcard characters.

Operator
Definition
Examples
*

Replaces any number of characters

| where application.name = "Microsoft*"

Returns application names starting with "Microsoft"

| where application.name = "*Microsoft*"

Returns application names containing "Microsoft"

?

Replaces any single character

| where device.operating_system.name == "Windows 1?"

Returns operating system names with versions above 10, such as Windows 10 and Windows 11.

...

Last updated 2 months ago

NQL logical operators
Use logical operators with NQL where clause to filter your NQL query results.

NQL 'in'
Use the in logical operator to check if a metric value is in the list.

Examples:


Copy
...
| where package.name in [ "MS Teams", "Zoom" ]

Copy
...
| where code in [94011, 94031]
NQL '!in'
Use the !in logical operator to check if a metric value is not in the list.

Examples:


Copy
...
| where hardware.type !in [virtual, null]

Copy
...
| where code !in [ 403, 404]
NQL 'contains'
Use the contains logical operator to check if a string is contained in an array of strings.

Example:


Copy
...
| where monitor.tags contains 'VDI'
NQL '!contains'
Use the !contains logical operator to check if a string is not contained in an array of strings.

Example:


Copy
...
| where monitor.tags !contains 'VDI'
Last updated 2 months ago

NQL bitwise operators
Use bitwise and and or operators in the NQL where clause to apply multiple filters or create complex conditions.

NQL 'and'
Use the and operator to combine multiple conditions and retrieve only records that meet all conditions simultaneously.

Example:

Retrieve binaries where the name is "chrome.exe" and they run on Windows.


Copy
binaries during past 30d
| where name == "chrome.exe" and platform == windows 
NQL 'or'
Use the or operator to combine multiple conditions and retrieve records that meet at least one of them.

Example:

Retrieve binaries where the name contains "chrome" or "firefox".


Copy
binaries  during past 7d
| where name == "*chrome*" or name == "*firefox*"
Last updated 2 months ago

NQL functions
Functions are predefined operations that aggregate, format or extract data, enabling further analysis. They include operations like summing, averaging, and counting, often within grouped data.

Depending on the specific function, you can use it with:

compute

summarize

list

sort

where

Syntax

Copy
...
... <metric>.<function>.(<optional: function parameters>)
Examples
Aggregate functions
sum()

Copy
devices during past 7d 
| include device_performance.system_crashes during past 7d 
| compute number_of_crashes = number_of_system_crashes.sum()
countif()

Copy
collaboration.sessions during past 24h
| summarize ratio_of_poor_calls = countif(session.audio.quality = poor or session.video.quality = poor) / count() by connection_type
Format function
as()

Copy
devices
| summarize total_cost = count() * 2000
| list total_cost.as(format = currency, code  = USD)
Timestamp functions
time_elapsed()

Copy
devices
| where operating_system.last_update.time_elapsed() > 15d
hour()

Copy
device_performance.events during past 24h
| where start_time.hour() >= 9 and end_time.hour() <= 17
Aggregated metrics
It's important to differentiate between functions and aggregated metrics. The data model contains various aggregated metrics simplifying access to information. They are defined as fields of the data model.

Field
Description
Example
<metric>.avg

Average value of the metric aggregated in the bucket.

where unload_event.avg > 1.0

<metric>.sum

Sum of all values of the metric aggregated in the bucket.

where unload_event.sum == 10

<metric>.count

Number of aggregated values in the bucket.

where unload_event.count <= 4

<metric>.min

Minimum value of the metric in the bucket.

where unload_event.min < 1.0

<metric>.max

Maximum value of the metric in the bucket.

where unload_event.max > 1.0

Smart aggregates
A smart aggregate is an aggregate on an aggregated metrics that abstracts the underlying computation. They are not fields of the data model. During the execution of a query, the parser computes them on the fly.

Aggregate
Description
<metric>.avg()

Average value of the metric.
It is equivalent to <metric>.sum.sum() and <metric>.count.sum()

<metric>.sum()

Sum of all values of the metric.
It is equivalent to <metric>.sum.sum()

<metric>.max()

Maximum value of the metric.
It is equivalent to <metric>.max.max()

<metric>.min()

Minimum value of the metric.
It is equivalent to <metric>.min.min()

<metric>.p95()

95th percentile of the metric.

<metric>.p05()

5th percentile of the metric.

<metric>.count()

Number of aggregated values.
It is equivalent to <metric>.count.sum()

Example:

Retrieve a list of devices with less than 3GB of average free memory. The following query includes the free_memory.avg() smart aggregate in a compute clause. It computes the average free memory based on the same underlying data points as free_memory.avg aggregated metrics. It is equivalent to free_memory.avg.avg().


Copy
devices during past 7d
| with device_performance.events during past 7d
| compute avg_free_memory = free_memory.avg()
| where avg_free_memory < 3GB
Chaining of functions
You can call more than one function on the same field. Currently, the system supports chaining of the time_elapsed() function.

Example:

The following query returns the list of devices with the time elapsed since their last fast startup.


Copy
devices
| include device_performance.boots
| where type == fast_startup
| compute time_since_last_fast_startup = time.last().time_elapsed()
In the following section you can find a list of all available functions with usage rules and examples.

Last updated 2 months ago

NQL as()
The as() function allows you to display the output from NQL queries in a formatted manner by assigning specific formatting information to the metrics. The Nexthink interface displays the data with the unit specified.​ The as() function is supported in the list and sort clauses.

Example:


Copy
devices
| summarize total_cost = count() * 2000
| list total_cost.as(format = currency, code  = USD)

Use the following formatting options:

Formatter
Unit
Example query
energy

Wh

kWh

MWh

GWh


Copy
device_performance.events during past 1d
| where device.operating_system.platform = macos
| where hardware.type == laptop
| summarize estimated_energy = ((8 * 0.070) * device.count())
| list estimated_energy.as( format = energy )
weight

g

kg

t

kt


Copy
devices
| where hardware.type = laptop or hardware.type = desktop
| summarize no_of_devices = (count()^1) * 422.5
| list no_of_devices.as( format = weight )
currency

Specify the currency using the additional code parameter.

CAD: CA$

CHF: CHF

GBP: £
EUR: €

USD: $


Copy
devices
| summarize total_cost = count() * 2000
| list total_cost.as(format = currency, code  = USD)

Copy
devices
| summarize total_cost = count() * 2000
| list total_cost.as(format = currency, code  = EUR)
percent

%

To compute the percentage, the system multiplies the value by 100 and appends the % symbol. For example, a value of 0.47 is converted to 47%.


Copy
devices
| with execution.crashes
| summarize impacted_devices = countif(operating_system.name == "Windows 10 Pro 22H2 (66 bits)")/count()
| list impacted_devices.as(format = percent)
bitrate

bps

Kbps

Mbps

Gbps


Copy
connectivity.events during past 30d
| where primary_physical_adapter.type == wifi
| where wifi.signal_strength.avg <= -67 or (context.device_platform == "macOD" and wifi.noise_level.avg >= -80)
| summarize average_receive_rate = wifi.receive_rate.avg()
| list average_receive_rate.as(format = bitrate)
Last updated 1 year ago

NQL avg()
The avg() command returns the average value of a metric.

Examples:


Copy
devices during past 30min
| include web.page_views during past 30min
| compute avg_page_load = page_load_time.overall.avg()

Copy
web.page_views during past 30min
| summarize c = page_load_time.overall.avg()
Last updated 2 months ago

NQL count()
The count() function returns the number of unique objects or punctual events.

Using with the ‘compute’ clause
For objects:
It returns the number of unique objects.


Copy
devices during past 7d
| include execution.events during past 7d
| compute number_of_devices = device.count()
For punctual events:
It computes the number of events per object.


Copy
devices during past 7d
| include execution.crashes during past 7d
| compute number_of_crashes_ = count()
For sampled events:
It is not recommended to use the count() function on sampled events as it will return the number of data samples, not the actual number of events.

Using with the ‘summarize’ clause
When used with the summarize clause, the count() function always returns the number of records in the root table.

For objects:
It returns the number of objects.


Copy
devices during past 7d
| summarize c1 = count()
For punctual events:
It returns the number of events.


Copy
execution.crashes during past 7d
| summarize c1 = number_of_crashes.count()
Note that the following query returns the number of records of root table (in this case, devices), not the number of unique events. To count events, use the sum() function in the summarize clause instead.


Copy
devices during past 7d
| include execution.crashes during past 7d
| compute number_of_crashes_ = number_of_crashes.count()
| summarize c1 = number_of_crashes_.count()
For sampled events:
It is not recommended to use the count() function on sampled events as it will return the number of data samples, not the actual number of events.

Last updated 2 months ago

NQL countif()
The countif() function counts the number of rows that match specified criteria. It accepts extra arguments, the conditions that determine which rows to include. Provide one or more conditions using bitwise operators.

The following query returns a ratio of poor calls. It includes two conditions that filter rows with either poor audio or poor video quality.


Copy
collaboration.sessions during past 24h
| summarize ratio_of_poor_calls = countif(session.audio.quality = poor or session.video.quality = poor) / count() by connection_type
Connection type
Ratio of poor calls
Ethernet

0

Wi-Fi

0.01

Cellular

0.13

Last updated 2 months ago

NQL last()
The last() function returns the last value recorded.


Copy
devices
| with execution.events past 7d
| where binary.name == "zoom*"
| compute last_execution = timestamp.last()
| list last_execution, operating_system.name, device.name, operating_system.platform
Last updated 2 months ago

NQL max()
The max() function returns the maximum value recorded.


Copy
devices during past 30min
| summarize device_max_memory = hardware.memory.max() 
Last updated 2 months ago

NQL min()
The min() function returns the minimum value recorded.


Copy
devices during past 7d
| summarize c = collector.tag_id.min()
Last updated 2 months ago

NQL sum()
The sum() function returns the total sum of a numeric column.


Copy
users
| include web.events during past 24h
| where application.name == "Salesforce"
| compute usage_duration = duration.sum()
| where usage_duration != 0
| sort usage_duration desc
Username
AD → Full name
Usage duration
jmansel@kanopy

Jessy Mansel

34min 20s

lahmir@kanopy

Latrell Ahmir

26min 16s

abilete@kanopy

Adhibhaar Bilete

25min 46s

kkamilah@kanopy

Kolton Mariam Kamilah

24min 53s

jadnan@kanopy

Jimmie Adnan

24min 47s

bkyrie@kanopy

Braedyn Branden Kyrie

24min 7s

melyssa@kanopy

Micah Mackenna Elyssa

23min 6s

csami@kanopy

Cody Treyton Sami

23min 3s

Last updated 2 months ago

NQL sumif()
The sumif() command summarizes values from the rows that match specified criteria. It accepts extra arguments, the conditions that determine which rows to include.


Copy
devices during past 24h
| include web.page_views during past 24h
| where application.name == "Confluence"
| compute total_cnt = page_view.number_of_page_views.sumif(page_view.experience_level == frustrating)
| sort total_cnt desc
| limit 10
Name
Entity
Hardware → Device model
Hardware → Device type
Operating system -> Name
Total cnt
NXDOCS-1704355664

Switzerland

MacBookPro18,3

laptop

macOS Sonoma 14.2.1 (ARM 64 bits)

7

NXDOCS-1704355669

Spain

ThinkPad P1 Gen 4i

laptop

Windows 11 Enterprise 22H2 (64 bits)

6

NXDOCS-1704355676

India

HP ZBook Power G7 Mobile Workstation

laptop

Windows 11 Enterprise 22H2 (64 bits)

6

NXDOCS-1704355687

Spain

ThinkPad T14s Gen 3

laptop

Windows 11 Pro 22H2 (64 bits)

5

NXDOCS-1704355691

India

MacBookPro16,1

laptop

macOS Ventura 13.6.3 (64 bits)

5

NXDOCS-1704355695

United States

MacBookPro18,3

laptop

macOS Sonoma 14.2.1 (ARM 64 bits)

5

NXDOCS-1704355714

United Kingdom

MacBookPro18,2

laptop

macOS Sonoma 14.2.1 (ARM 64 bits)

4

NXDOCS-1704355719

Switzerland

MacBookPro18,1

laptop

macOS Ventura 13.6.3 (ARM 64 bits)

4

NXDOCS-1704355725

India

ThinkPad T480

laptop

Windows 10 Pro 22H2 (64 bits)

4

NXDOCS-1704355730

United States

MacBookPro18,1

laptop

macOS Sonoma 14.0.0 (ARM 64 bits)

3

Last updated 2 months ago

NQL time_elapsed()
The time_elapsed() function calculates the time elapsed since an event. The function returns the values in seconds.

Use this function with a field of datetime data type, for example:

the last_seen field from the devices table

time from the device_performance.boots table

The timeframe specified in your query does not restrict the values returned by the time_elapsed() function. For example, the following query retrieves only devices active in the past day, but the values returned by time_elapsed() may extend beyond that timeframe.


Using with the ‘where’ clause
Use the time_elapsed() function in a where clause.

Example:

Retrieve the list of devices where the last operating system update was more than 15 days ago.


Copy
devices
| where operating_system.last_update.time_elapsed() > 15d
Using with the ‘list’ clause
Use the time_elapsed() function in a list clause.

Example:

List devices and the time elapsed from their last startup.


Copy
devices
| include device_performance.boots
| where type == fast_startup
| compute last_fast_startup_time = time.last()
| list name, last_fast_startup_time.time_elapsed()
Name
Last fast startup time → time elapsed
device-10d267d2

1w 0d 1h 8min 22s 0ms

device-d1d5abc9

17h 38min 22s 0ms

device-5117c4c3

3w 1d 10h 33min 8s 0ms

device-16834449

57min 18s 0ms

…

…

Using with the ‘compute’ clause
Use the time_elapsed() function in a compute clause.

Example:

List devices and the time elapsed from their last startup. Applying chaining of functions (call multiple functions on the same field).


Copy
devices
| include device_performance.boots
| where type == fast_startup
| compute time_since_last_fast_startup = time.last().time_elapsed()
Last updated 2 months ago

NQL datetime functions
Datetime functions return specific time components— such as hour, day of the week, day of the month—from the timestamp field. This enables you to identify patterns or trends within time windows, for example, business hours, business week, or specific days of the month.

Example
Retrieve device performance data for business hours within the last 24 hours.


Copy
device_performance.events during past 24h
| where start_time.hour() >= 9 and end_time.hour() <= 17
Available functions
hour()
Description: This function allows you to extract the hour from a given timestamp.

Returns: Numbers from 0 to 23.

NQL query example: View all events that occurred during business hours—e.g., between 9 am and 5 pm.


Copy
device_performance.events during past 24h
| where start_time.hour() >= 9 and end_time.hour() <= 17
To extract an hour from the date, you need to use a timeframe expressed in minutes or hours in your query—for example, use during past 168h instead of during past 7d.

day()
Description: This function allows you to extract the day of the month from a given date.

Returns: Numbers from 1 to 31.

NQL query example: Retrieve device performance data from the first week of the month.


Copy
device_performance.events during past 30d
| where start_time.day() >= 1 and end_time.day() <= 7
day_of_week()
Description: This function allows you to extract the day of the week from a given date.

Returns: Numbers from 1 to 7, where 1 represents Monday and 7 represents Sunday.

NQL query example: Retrieve device performance data from working days—e.g., Monday to Friday.


Copy
device_performance.events during past 30d
| where start_time.day_of_week() >= 1 and end_time.day_of_week() < 6
Timezone parameter
By default, the system returns time values in your local timezone, and datetime functions return values in your local time. Provide different timezones using the timezone parameter in the function. Datetime functions return values in the specified timezone. Refer to Supported timezones in datetime functions for more information.

While datetime functions retrieve the time component of a timestamp in the specified timezone, they do not alter the timeframe used for time selection.

Refer to Timezones for more information on how specific time selection formats affect the timeframe included in your query.

Example

You are in Helsinki (EET timezone) at 8:00 AM and want to retrieve device performance data during business hours in London (GMT timezone). Use the following query to narrow the data to London's business hours.


Copy
device_performance.events during past 24h
| where start_time.hour(timezone = 'GMT') >= 9 and end_time.hour(timezone = 'GMT') <= 17
The image below shows the same query, but with different time formats included in a list clause. The timeframe and values returned depend on the time format.

The timeframe reflects time in your current timezone.

The start_time returns the full timestamp of the event in your current timezone.

The start_time.hour() returns the hour of the event in your current timezone.

The start_time.hour(timezone = 'GMT') returns the hour of the event in London's timezone.

Supported timezones in datetime functions
The timezone parameter in the datetime function supports the following timezones:


Copy
Africa/Abidjan
Africa/Accra
Africa/Addis_Ababa
Africa/Algiers
Africa/Asmara
Africa/Asmera
Africa/Bamako
Africa/Bangui
Africa/Banjul
Africa/Bissau
Africa/Blantyre
Africa/Brazzaville
Africa/Bujumbura
Africa/Cairo
Africa/Casablanca
Africa/Ceuta
Africa/Conakry
Africa/Dakar
Africa/Dar_es_Salaam
Africa/Djibouti
Africa/Douala
Africa/El_Aaiun
Africa/Freetown
Africa/Gaborone
Africa/Harare
Africa/Johannesburg
Africa/Juba
Africa/Kampala
Africa/Khartoum
Africa/Kigali
Africa/Kinshasa
Africa/Lagos
Africa/Libreville
Africa/Lome
Africa/Luanda
Africa/Lubumbashi
Africa/Lusaka
Africa/Malabo
Africa/Maputo
Africa/Maseru
Africa/Mbabane
Africa/Mogadishu
Africa/Monrovia
Africa/Nairobi
Africa/Ndjamena
Africa/Niamey
Africa/Nouakchott
Africa/Ouagadougou
Africa/Porto-Novo
Africa/Sao_Tome
Africa/Timbuktu
Africa/Tripoli
Africa/Tunis
Africa/Windhoek
America/Adak
America/Anchorage
America/Anguilla
America/Antigua
America/Araguaina
America/Argentina/Buenos_Aires
America/Argentina/Catamarca
America/Argentina/ComodRivadavia
America/Argentina/Cordoba
America/Argentina/Jujuy
America/Argentina/La_Rioja
America/Argentina/Mendoza
America/Argentina/Rio_Gallegos
America/Argentina/Salta
America/Argentina/San_Juan
America/Argentina/San_Luis
America/Argentina/Tucuman
America/Argentina/Ushuaia
America/Aruba
America/Asuncion
America/Atikokan
America/Atka
America/Bahia
America/Bahia_Banderas
America/Barbados
America/Belem
America/Belize
America/Blanc-Sablon
America/Boa_Vista
America/Bogota
America/Boise
America/Buenos_Aires
America/Cambridge_Bay
America/Campo_Grande
America/Cancun
America/Caracas
America/Catamarca
America/Cayenne
America/Cayman
America/Chicago
America/Chihuahua
America/Ciudad_Juarez
America/Coral_Harbour
America/Cordoba
America/Costa_Rica
America/Creston
America/Cuiaba
America/Curacao
America/Danmarkshavn
America/Dawson
America/Dawson_Creek
America/Denver
America/Detroit
America/Dominica
America/Edmonton
America/Eirunepe
America/El_Salvador
America/Ensenada
America/Fort_Nelson
America/Fort_Wayne
America/Fortaleza
America/Glace_Bay
America/Godthab
America/Goose_Bay
America/Grand_Turk
America/Grenada
America/Guadeloupe
America/Guatemala
America/Guayaquil
America/Guyana
America/Halifax
America/Havana
America/Hermosillo
America/Indiana/Indianapolis
America/Indiana/Knox
America/Indiana/Marengo
America/Indiana/Petersburg
America/Indiana/Tell_City
America/Indiana/Vevay
America/Indiana/Vincennes
America/Indiana/Winamac
America/Indianapolis
America/Inuvik
America/Iqaluit
America/Jamaica
America/Jujuy
America/Juneau
America/Kentucky/Louisville
America/Kentucky/Monticello
America/Knox_IN
America/Kralendijk
America/La_Paz
America/Lima
America/Los_Angeles
America/Louisville
America/Lower_Princes
America/Maceio
America/Managua
America/Manaus
America/Marigot
America/Martinique
America/Matamoros
America/Mazatlan
America/Mendoza
America/Menominee
America/Merida
America/Metlakatla
America/Mexico_City
America/Miquelon
America/Moncton
America/Monterrey
America/Montevideo
America/Montreal
America/Montserrat
America/Nassau
America/New_York
America/Nipigon
America/Nome
America/Noronha
America/North_Dakota/Beulah
America/North_Dakota/Center
America/North_Dakota/New_Salem
America/Nuuk
America/Ojinaga
America/Panama
America/Pangnirtung
America/Paramaribo
America/Phoenix
America/Port_of_Spain
America/Port-au-Prince
America/Porto_Velho
America/Puerto_Rico
America/Punta_Arenas
America/Rainy_River
America/Rankin_Inlet
America/Recife
America/Regina
America/Resolute
America/Rosario
America/Santa_Isabel
America/Santarem
America/Santiago
America/Santo_Domingo
America/Sao_Paulo
America/Scoresbysund
America/Shiprock
America/Sitka
America/St_Barthelemy
America/St_Johns
America/St_Kitts
America/St_Lucia
America/St_Thomas
America/St_Vincent
America/Swift_Current
America/Tegucigalpa
America/Thule
America/Thunder_Bay
America/Tijuana
America/Toronto
America/Tortola
America/Vancouver
America/Virgin
America/Whitehorse
America/Winnipeg
America/Yakutat
America/Yellowknife
Antarctica/Casey
Antarctica/Davis
Antarctica/DumontDUrville
Antarctica/Macquarie
Antarctica/Mawson
Antarctica/McMurdo
Antarctica/Palmer
Antarctica/Rothera
Antarctica/South_Pole
Antarctica/Syowa
Antarctica/Troll
Antarctica/Vostok
Arctic/Longyearbyen
Asia/Aden
Asia/Almaty
Asia/Amman
Asia/Anadyr
Asia/Aqtau
Asia/Aqtobe
Asia/Ashgabat
Asia/Ashkhabad
Asia/Atyrau
Asia/Baghdad
Asia/Bahrain
Asia/Baku
Asia/Bangkok
Asia/Barnaul
Asia/Beirut
Asia/Bishkek
Asia/Brunei
Asia/Calcutta
Asia/Chita
Asia/Choibalsan
Asia/Chongqing
Asia/Chungking
Asia/Colombo
Asia/Dacca
Asia/Damascus
Asia/Dhaka
Asia/Dili
Asia/Dubai
Asia/Dushanbe
Asia/Famagusta
Asia/Gaza
Asia/Harbin
Asia/Hebron
Asia/Ho_Chi_Minh
Asia/Hong_Kong
Asia/Hovd
Asia/Irkutsk
Asia/Istanbul
Asia/Jakarta
Asia/Jayapura
Asia/Jerusalem
Asia/Kabul
Asia/Kamchatka
Asia/Karachi
Asia/Kashgar
Asia/Kathmandu
Asia/Katmandu
Asia/Khandyga
Asia/Kolkata
Asia/Krasnoyarsk
Asia/Kuala_Lumpur
Asia/Kuching
Asia/Kuwait
Asia/Macao
Asia/Macau
Asia/Magadan
Asia/Makassar
Asia/Manila
Asia/Muscat
Asia/Nicosia
Asia/Novokuznetsk
Asia/Novosibirsk
Asia/Omsk
Asia/Oral
Asia/Phnom_Penh
Asia/Pontianak
Asia/Pyongyang
Asia/Qatar
Asia/Qostanay
Asia/Qyzylorda
Asia/Rangoon
Asia/Riyadh
Asia/Saigon
Asia/Sakhalin
Asia/Samarkand
Asia/Seoul
Asia/Shanghai
Asia/Singapore
Asia/Srednekolymsk
Asia/Taipei
Asia/Tashkent
Asia/Tbilisi
Asia/Tehran
Asia/Tel_Aviv
Asia/Thimbu
Asia/Thimphu
Asia/Tokyo
Asia/Tomsk
Asia/Ujung_Pandang
Asia/Ulaanbaatar
Asia/Ulan_Bator
Asia/Urumqi
Asia/Ust-Nera
Asia/Vientiane
Asia/Vladivostok
Asia/Yakutsk
Asia/Yangon
Asia/Yekaterinburg
Asia/Yerevan
Atlantic/Azores
Atlantic/Bermuda
Atlantic/Canary
Atlantic/Cape_Verde
Atlantic/Faeroe
Atlantic/Faroe
Atlantic/Jan_Mayen
Atlantic/Madeira
Atlantic/Reykjavik
Atlantic/South_Georgia
Atlantic/St_Helena
Atlantic/Stanley
Australia/ACT
Australia/Adelaide
Australia/Brisbane
Australia/Broken_Hill
Australia/Canberra
Australia/Currie
Australia/Darwin
Australia/Eucla
Australia/Hobart
Australia/LHI
Australia/Lindeman
Australia/Lord_Howe
Australia/Melbourne
Australia/North
Australia/NSW
Australia/Perth
Australia/Queensland
Australia/South
Australia/Sydney
Australia/Tasmania
Australia/Victoria
Australia/West
Australia/Yancowinna
Brazil/Acre
Brazil/DeNoronha
Brazil/East
Brazil/West
Canada/Atlantic
Canada/Central
Canada/Eastern
Canada/Mountain
Canada/Newfoundland
Canada/Pacific
Canada/Saskatchewan
Canada/Yukon
CET
Chile/Continental
Chile/EasterIsland
CST6CDT
Cuba
EET
Egypt
Eire
EST
EST5EDT
Etc/GMT
Etc/GMT+0
Etc/GMT+1
Etc/GMT+10
Etc/GMT+11
Etc/GMT+12
Etc/GMT+2
Etc/GMT+3
Etc/GMT+4
Etc/GMT+5
Etc/GMT+6
Etc/GMT+7
Etc/GMT+8
Etc/GMT+9
Etc/GMT0
Etc/GMT-0
Etc/GMT-1
Etc/GMT-10
Etc/GMT-11
Etc/GMT-12
Etc/GMT-13
Etc/GMT-14
Etc/GMT-2
Etc/GMT-3
Etc/GMT-4
Etc/GMT-5
Etc/GMT-6
Etc/GMT-7
Etc/GMT-8
Etc/GMT-9
Etc/Greenwich
Etc/UCT
Etc/Universal
Etc/UTC
Etc/Zulu
Europe/Amsterdam
Europe/Andorra
Europe/Astrakhan
Europe/Athens
Europe/Belfast
Europe/Belgrade
Europe/Berlin
Europe/Bratislava
Europe/Brussels
Europe/Bucharest
Europe/Budapest
Europe/Busingen
Europe/Chisinau
Europe/Copenhagen
Europe/Dublin
Europe/Gibraltar
Europe/Guernsey
Europe/Helsinki
Europe/Isle_of_Man
Europe/Istanbul
Europe/Jersey
Europe/Kaliningrad
Europe/Kiev
Europe/Kirov
Europe/Kyiv
Europe/Lisbon
Europe/Ljubljana
Europe/London
Europe/Luxembourg
Europe/Madrid
Europe/Malta
Europe/Mariehamn
Europe/Minsk
Europe/Monaco
Europe/Moscow
Europe/Nicosia
Europe/Oslo
Europe/Paris
Europe/Podgorica
Europe/Prague
Europe/Riga
Europe/Rome
Europe/Samara
Europe/San_Marino
Europe/Sarajevo
Europe/Saratov
Europe/Simferopol
Europe/Skopje
Europe/Sofia
Europe/Stockholm
Europe/Tallinn
Europe/Tirane
Europe/Tiraspol
Europe/Ulyanovsk
Europe/Uzhgorod
Europe/Vaduz
Europe/Vatican
Europe/Vienna
Europe/Vilnius
Europe/Volgograd
Europe/Warsaw
Europe/Zagreb
Europe/Zaporozhye
Europe/Zurich
GB
GB-Eire
GMT
GMT0
Greenwich
Hongkong
HST
Iceland
Indian/Antananarivo
Indian/Chagos
Indian/Christmas
Indian/Cocos
Indian/Comoro
Indian/Kerguelen
Indian/Mahe
Indian/Maldives
Indian/Mauritius
Indian/Mayotte
Indian/Reunion
Iran
Israel
Jamaica
Japan
Kwajalein
Libya
MET
Mexico/BajaNorte
Mexico/BajaSur
Mexico/General
MST
MST7MDT
Navajo
NZ
NZ-CHAT
Pacific/Apia
Pacific/Auckland
Pacific/Bougainville
Pacific/Chatham
Pacific/Chuuk
Pacific/Easter
Pacific/Efate
Pacific/Enderbury
Pacific/Fakaofo
Pacific/Fiji
Pacific/Funafuti
Pacific/Galapagos
Pacific/Gambier
Pacific/Guadalcanal
Pacific/Guam
Pacific/Honolulu
Pacific/Johnston
Pacific/Kanton
Pacific/Kiritimati
Pacific/Kosrae
Pacific/Kwajalein
Pacific/Majuro
Pacific/Marquesas
Pacific/Midway
Pacific/Nauru
Pacific/Niue
Pacific/Norfolk
Pacific/Noumea
Pacific/Pago_Pago
Pacific/Palau
Pacific/Pitcairn
Pacific/Pohnpei
Pacific/Ponape
Pacific/Port_Moresby
Pacific/Rarotonga
Pacific/Saipan
Pacific/Samoa
Pacific/Tahiti
Pacific/Tarawa
Pacific/Tongatapu
Pacific/Truk
Pacific/Wake
Pacific/Wallis
Pacific/Yap
Poland
Portugal
PRC
PST8PDT
ROK
Singapore
Turkey
UCT
Universal
US/Alaska
US/Aleutian
US/Arizona
US/Central
US/Eastern
US/East-Indiana
US/Hawaii
US/Indiana-Starke
US/Michigan
US/Mountain
US/Pacific
US/Samoa
UTC
WET
W-SU
Zulu
Last updated 2 months ago

NQL catalog
The aim of the NQL catalog is to help you successfully query the data in the Nexthink web interface. Click on the links below to access examples of queries most commonly used in specific Nexthink modules. Go through the examples listed on each page and pick the one most similar to your use case. Copy the query and adjust it to your needs, or use it as a starting point for writing your own query.

DEX score NQL examples
Use Nexthink Query Language (NQL) to access DEX score data and other relevant information.

NQL data structure
The dex.scores and dex.application_scores tables contain score data. The system computes each score once a day at 00:00 UTC for a combination of user and device objects, and device dimensions active over the last 7 days. For example, employee A who used device 1 and device 2 over the last 7 days will have two sets of score data for the current day.

DEX score V3 primarily focuses on user-centric experience management rather than a device-centric approach, which differs from previous versions. Even though starting an NQL query with the device table is technically possible, you may witness deviations between the device's metric value and its score. Refer to the FAQ section of the Computation of the DEX score (available to Nexthink Community users) documentation for more information.

dex.scores table
The dex.scores table contains score data for the endpoint and collaboration scores and their subscores. The system structures a set of scores as follows:

Each node of the DEX score has a score value with the syntax [node name]_value. For example, if you want the score for logon speed, you would enter score.endpoint.logon_speed_value.

In addition, each node has a score impact value with the syntax [node name]_score_impact. This value represents the estimated decrease in the Technology component of the DEX score due to issues monitored by this node.

dex.application_scores table
The dex.application_scores table contains score data for the application scores and their subscores. This table is linked to users, devices and application objects. The system structures a set of scores as follows:

node.type represents the type of node of the application score structure:

page_loads

transactions

web_reliability

crashes

freezes

application

node.value indicates the score of a node of the application score structure. It must be used with the field application_score.node.type to specify the target node score.

node.score_impact indicates the estimated decrease in the Technology component of the DEX score due to issues monitored by this node. It must be used with the field application_score.node.type to specify the target score impact.

Refer to the List of hard metrics and their default thresholds documentation (available for Nexthink Community users) for the full list of nodes and their respective NQL names.

score_impact values
Both dex.scores and dex.application_scores include a score_impact value for all nodes, i.e., [node name]_score_impact. This value estimates the number of points removed from the DEX score due to user-level issues for the node. For example, logon_speed_score_impact contains the estimated impact on the DEX score for a user due to slow logons.

To compute the impact on DEX score of a node for a population, use the following formula:

image-2024-03-18-09-22-56-620.png
Examples of NQL queries

Copy
users
| include dex.scores during past 24h
| compute DEX_score_per_user = value.avg(), c1 = count()
| where c1 > 0
| summarize Overall_DEX_score =  DEX_score_per_user.avg()

Copy
users
| include dex.scores during past 24h
| where context.location.country == "Switzerland"
| compute Virtualization_score = endpoint.virtual_session_lag_value.avg(), c1 = count()
| where c1 > 0
| summarize Overall_virtualization_score = Virtualization_score.avg()

Copy
devices
| include dex.scores during past 24h
| compute DEX_score_per_device = value.avg(), c1 = count()
| where c1 > 0
| summarize Overall_DEX_score_per_OS_platform = DEX_score_per_device.avg() by operating_system.platform

Copy
users 
| include dex.application_scores during past 24h 
| where application.name == "miro" and node.type == application 
| compute application_score_per_user = value.avg()
| include dex.application_scores during past 24h 
| where application.name == "miro" and node.type == page_loads 
| compute page_load_score_per_user = value.avg()
| summarize Overall_application_score = application_score_per_user.avg(), Overall_page_load_score = page_load_score_per_user.avg() 

Copy
users
| with dex.scores during past 24h
| compute DEX_score = score.value.avg()
| list name, DEX_score
| sort DEX_score asc
| limit 50

Copy
users
| include dex.scores during past 24h
| compute logon_speed_score_impact_per_user = endpoint.logon_speed_score_impact.avg(), dex_score_per_user = dex.score.value.avg()
| summarize Impact_of_logon_speed_on_technology_score = (logon_speed_score_impact_per_user.avg()*countif(logon_speed_score_impact_per_user != NULL))/countif(dex_score_per_user != NULL)

Copy
users
| include dex.scores during past 24h
| compute logon_speed_score_impact_per_user = endpoint.logon_speed_score_impact.avg(), DEX_score_per_user = value.avg()
| where DEX_score_per_user != NULL
| summarize total_logon_speed_score_impact = logon_speed_score_impact_per_user.avg()*countif(logon_speed_score_impact_per_user != NULL)/count()
Considerations
Timeframe associated with the dex.score table
A user or device object without any data over the last 7 days will have no score.

The score computed today at 00:00 UTC is associated with today's date, not yesterday's date.

This means that querying the dex.scores data with during past 7d is not correct, as this returns seven days of data points, and each data point is already a rolling window of 7 days. The data should be queried for only 1 day, for example

dex.scores during past 24h

dex.scores on 2023-10-30

Example:


Copy
users
| include dex.scores during past 24h
| compute software_reliability_score_per_user = endpoint.software_reliability_value.avg()
The daily computation of the scores starts at 00:00 UTC, but may take several hours to complete. Once the computation is finished, Nexthink tags the outcome with 04:00 UTC.

Timeframe associated with raw data tables
To compare raw metric tables, for example, session.logins, web.page_views with corresponding scores, the timeframe must mimic the one used in the DEX score computation:

00:00 UTC today - 7d to 00:00 UTC today

Ensure that your timeframe matches the timezone of your browser.

Example:

Your browser timezone is CET (i.e., UTC + 1).

You want to compare the raw metrics with the DEX score of 2024-01-17.

You are interested in checking the page loads trends of Outlook.


Copy
web.page_views from 2024-01-10 01:00:00 to 2024-01-17 01:00:00
| where application.name == "outlook" and user.name == "TBD"
| summarize average_page_load_per_hour = page_load_time.overall.avg() by 1h
Hourly samples
The DEX scores take into account hourly samples of data. This data is either an average of a field or a sum of events over the past hour, depending on the metric type. To understand if a metric has breached the configured score thresholds, use hourly aggregation, not 5-minute or 15-minute time buckets.

Example:


Copy
Session.logins from 2024-01-10 01:00:00 to 2024-01-17 01:00:00
| summarize average_logon_time = time_until_desktop_is_visible.avg() by 1h
Computing DEX scores for populations
Computing the DEX scores of a population requires first computing the score of each employee and then averaging them across the entire population.

First, aggregate the DEX score per user or device and then for the population.

In this case, you should not start an NQL query with dex.scores during past 24h, but with the users or the devices table to compute individual scores. Use the summarize statement to compute the DEX score for the population.

Example:


Copy
users
| include dex.scores during past 24h
| compute dex_score_per_user = value.avg()
| summarize dex_score = dex_score_per_user.avg()
Looking at the right combination of user, device, and device dimensions
The system computes a score for a combination of user, device, and device dimensions such as geolocation, or location type. The rationale behind this approach is to enable advanced filtering on each employee context to extract insights.

This means employees who have changed their location during the last 7 days will have several scores for different geolocation or location types. Additionally, a device used by multiple employees will also have several scores.

To compare the raw metric with its corresponding score, you must:

Understand the different dimensions associated with the score data.

Apply the same dimensions when looking at the raw data.

Example:


Copy
device_performance.boots from 2024-01-31 01:00:00 to 2024-02-06 01:00:00
| where device.name == "TBD" and context.location.type == "Remote" and context.location.state == "Vaud"
| summarize average_boot_duration_per_hour = duration.avg() by 1h
Device events without an employee associated with the device
Some metrics used in the DEX score computation do not have any user association in the Nexthink data model, for example:

device_performance.events

device_performance.boots

device_performance.hard_resets

device_performance.system_crashes

connectivity.events

The DEX score pipeline keeps a list of recent users on a device so that these events can be associated with these users’ DEX scores. If no recent users can be found on the device, these events will not be factored into the device’s DEX scores.

Last updated 10 months ago

Detecting issues impacting multiple devices
Refer to the Alerts FAQ to learn how to investigate and query devices associated with an existing alert, using NQL.

Detect issues impacting multiple devices to allow application and network L2+ teams to proactively respond to global issues in their specific areas. Notify relevant application owners about issues impacting their applications. Using the following use cases, evaluate:

The number of impacted devices or users, for example, the number of devices with specific application crashes.

Frequent issues across devices, for example, the number of specific application crashes across all devices.

Both approaches are vital and often complement each other. Use either approach when configuring monitor trigger conditions to avoid triggering alerts and sending notifications when issues are irrelevant to the recipient. For example, the system triggers an alert when the number of specific application crashes across all devices exceeds 20 and affects more than 5 devices. The system then notifies the application owner.

The following sections describe two use cases in detail.

Monitoring the number of devices or users with issues
Detect the number of devices or users with an issue to proactively monitor issues impacting multiple devices.

Create an NQL query that returns a summarized number of devices. Count only the crashes that happened while the application was running in the foreground. Optionally, group your results using the by keyword to group your results. The system triggers an alert per group.


Copy
devices
| with execution.crashes during past 24h
| where binary.name = "outlook.exe"
| compute crashes = countif(process_visibility == foreground)
| summarize nr_of_devices = count() by entity
Notifications
The system sends notifications for all devices at once, or if the query includes the by clause, for each group separately. Only the number of devices is included in the notification as a value. The details of all devices impacted are available in Nexthink web interface.

Alerts overview dashboard
In the Alerts overview dashboard, the alert is displayed in a single line, or if grouping has been added, the alert is displayed for each group in a separate line with context about the grouping.

Monitoring frequent issues across devices
Detect an issue across multiple devices which is reflected in an aggregated metric value.

Create an NQL query that returns a summarized metric value. Optionally, group your results using the by keyword. The system triggers an alert per group.


Copy
execution.crashes during past 24h
| summarize 
  total_number_of_crashes = count(), 
  devices_with_crashes = device.count()
by binary.name
Notifications
The system sends notifications for a single metric, or if the query includes the by clause, for each group separately. The notifications contain information about breeched values for each metric defined in the condition.

Alerts overview dashboard
In the Alerts overview dashboard, the alert is displayed in a single line without the context-related label. If grouping has been added, the alert is displayed for each group in a separate line with context about the grouping.

Refer to the NQL examples below and the NQL data model documentation for more information about NQL.

NQL Examples
Below is a list of NQL query examples to help you create and edit monitors. Review the queries and pick the one most similar to the monitor you create or edit. Copy the query and adjust it to your use case, including the thresholds that have been provided as an example.

Detect specific web errors for an application.
This NQL query returns the aggregated number of errors and devices with errors for a specific application and triggers the alert per specific error code separately:


Copy
web.errors during past 1h
| where application.name  in ["Jenkins"] 
| where error.code !in [405, 404, 403]
| summarize nr_of_devices_impacted = device.count(), nr_of_errors = count() by label

Detect applications with a high web error ratio.
Select other thresholds to make sure there is enough usage volume and that there are enough issues to avoid false positives.


Copy
application.applications
| with web.page_views during past 60min
| where is_soft_navigation = false
| compute total_number_of_page_views = number_of_page_views.sum(), all_users = user.count()
| with web.errors during past 60min
| compute number_page_views_with_error = error.number_of_errors.sum(), users_with_errors = user.count()
| summarize web_errors_ratio = number_page_views_with_error.sum() * 100 / total_number_of_page_views.sum(), number_of_errors = number_page_views_with_error.sum(), users_with_issues = users_with_errors.sum(), ratio_of_users_with_issues = users_with_errors.sum() * 100 / all_users.sum() by application.name

Detect a high number of crashes for binaries.

Copy
execution.crashes during past 24h
| summarize total_number_of_crashes = count(), devices_with_crashes = device.count() by binary.name
| sort total_number_of_crashes desc

Detect a high number of devices with long boot time with Geolocation by country.
The long boot time is defined as time_until_desktop_is_visible>= 60s


Copy
devices
| with session.logins during past 24h
| compute total_devices = device.count(), avg_time_until_desktop_ready = time_until_desktop_is_ready.avg(), avg_time_until_desktop_visible = time_until_desktop_is_visible.avg()
| include session.logins during past 24h
| where time_until_desktop_is_visible>= 60s
| compute number_of_device_with_long_login = device.count()
| summarize percentage_of_devices_with_issue = number_of_device_with_long_login.sum() * 100 / total_devices.sum(), average_time_until_desktop_ready = avg_time_until_desktop_ready.avg(), average_time_until_desktop_visible = avg_time_until_desktop_visible.avg(), number_of_devices_with_issue = number_of_device_with_long_login.sum() by public_ip.country

Virtualization alert for when the average CPU queue length per desktop pool is >= 3

Copy
device_performance.events during past 30min
| where device.virtualization.desktop_pool != null
| summarize Average_cpu_queue_length = cpu_queue_length.avg() / number_of_logical_processors.avg() by device.virtualization.desktop_pool

Detecting issues impacting a single device or user
Refer to the Alerts FAQ to learn how to investigate and query devices associated with an existing alert, using NQL.

Detect an issue that occurs on a single device or for a single user, to help L1 support teams proactively respond and remediate before a user raises a ticket.

To configure a monitor that evaluates a metric per single device or user, create an NQL query that returns a list of devices with a computed metric without further aggregations. You can achieve this in two ways:

Start with the devices or users table, join either table with the events table and compute the metric.


Copy
devices
| with device_performance.system_crashes during past 7d
| compute total_number_of_system_crashes = number_of_system_crashes.sum()
Start with an events table, summarize the metric and add grouping by device.


Copy
device_performance.system_crashes during past 7d
| summarize total_number_of_system_crashes = number_of_system_crashes.sum() by device.collector.uid
Notifications:
The system sends notifications for each impacted device or user separately and includes the device or user name in the payload.

Alerts overview dashboard
In the Alerts overview dashboard, the alerts for all devices/users are combined and displayed in a single line. The impacted devices column informs you about the number of devices with alerts.

Considerations
Do not use this type of query if you expect a large number of objects to trigger an alert at once. Nexthink sets a limit of 500 simultaneous triggers for one monitor. Consider using Data Export or Webhooks for reporting purposes to external systems.

For this type of query do not use the summarize... by device.name syntax as it will not trigger an alert per device as you might expect.

NQL examples
Below is a list of NQL query examples to help you create and edit monitors. Review the queries and pick the one most similar to the monitor you are creating or editing. Copy the query and adjust it to your use case, including the thresholds that have been provided as an example.

Devices with a high number of system crashes per week ( >=3)
This alert is triggered per device with a device name in the payload.


Copy
devices
| with device_performance.system_crashes during past 7d 
| compute total_number_of_system_crashes = number_of_system_crashes.sum()
| sort total_number_of_system_crashes desc
or


Copy
device_performance.system_crashes during past 7d
| summarize total_number_of_system_crashes = number_of_system_crashes.sum() by device.collector.uid
In the second example, use the device.collector.uid for grouping. The system sends the device name in the notification.


Devices with a high system drive usage ratio in the last week ( >=90)
This alert is triggered per device with a device name in the payload.


Copy
devices during past 7d
| include device_performance.events during past 7d
| compute system_drive_usage_ratio_ = event.system_drive_usage.avg()/event.system_drive_capacity.avg()*100
| list system_drive_usage_ratio_

Unauthorized users accessing Salesforce app
This alert is triggered per user for those who have accessed the Salesforce app and are not from the Marketing or Sales departments. The username is in the alert payload.


Copy
users during past 1d
| include web.events during past 1d
| where application.name == "Salesforce"
| compute usage_time_ = event.duration.sum()
| include session.events during past 1d
| where device.organization.entity !in ["MarketingSales"]
| compute interaction_time_ = event.user_interaction_time.sum()

Last updated 1 year ago

Live Dashboards NQL examples
When configuring widgets using NQL, do not hardcode absolute time values or static timestamps.

This list of NQL query examples is designed to help you create live dashboard widgets.

KPI widget

Copy
...
summarize <kpi> = <sum() | count() | avg() | max() | min()>
Examples

Copy
web.errors during past 7d
| summarize total_errors = number_of_errors.sum() 

Copy
web.page_views during past 7d
| summarize 
  backend_dur_ratio = page_load_time.backend.sum() /
  page_load_time.overall.sum()

Copy
remote_action.executions during past 30d
| where status == success
| where purpose == remediation
| summarize amt_saved = (number_of_executions.sum()) * (20)
| list amt_saved.as(format = currency,code = usd)
Line chart

Copy
<event table> <time_duration>
...
summarize <kpi1>, <kpi2>, ... by <time_duration_granularity>
(list <time>, <kpi1>, <kpi2>, ...)
Examples

Copy
web.page_views during past 7d
| summarize 
    backend_duration = page_load_time.backend.avg() , 
    client_duration = page_load_time.client.avg() , 
    network_duration = page_load_time.network.avg() by 1d
 

Copy
web.page_views during past 7d
| summarize 
    backend_duration = page_load_time.backend.avg() , 
    client_duration = page_load_time.client.avg() , 
    network_duration = page_load_time.network.avg() by 1d
| list end_time, backend_duration, client_duration, network_duration

Copy
execution.events during past 15d
| where device.operating_system.name != "*server*"
| where 
  (device.hardware.type == laptop 
  or device.hardware.type == desktop)
| where binary.name in ["nxtsvc.exe", "nxtsvc"]
| summarize 
  Total_energy_consumption = 
  (((execution_duration.sum()) ^ (1)) / (3600)) * (30) 
  by 1d
| list 
  start_time, 
  Total_energy_consumption.as(format = energy)
Bar chart with disabled default breakdowns
In this case, as the default breakdowns are disabled, you should always specify by <segmentation1>,... in the query.


Copy
...
summarize <kpi1>, <kpi2>, ... by <segmentation1>, <segmentation2>, ...
Examples

Copy
device_performance.hard_resets  during past 7d
| summarize
    num_hard_resets = number_of_hard_resets.sum() ,
    num_devices = device.count()
   by
    device.operating_system.platform ,
    device.hardware.manufacturer ,
    device.hardware.model
| sort num_hard_resets desc

Copy
web.transactions 
| summarize nb_transactions = number_of_transactions.sum() 
   by application.name 
| sort nb_transactions desc

Copy
devices
| where device.public_ip.isp != null
| summarize 
  devices = device.name.count() 
  by device.public_ip.isp
| sort devices desc

Copy
workflow.executions during past 30d
| where status == success
| summarize amt_saved = (number_of_executions.sum()) * (100) 
  by trigger_method
| list trigger_method, amt_saved.as(format = currency,code = usd)
| sort amt_saved desc
Bar chart with enabled default breakdowns
In this case, as the default breakdowns are enabled, you can omit by <segmentation1>,... from the query since the system defaults to the options in the breakdown dropdown.


Copy
...
summarize <kpi1>, <kpi2>, ...
Examples
Since default breakdowns are enabled, the system displays web.page_views according to the available default options in the breakdown dropdown.
In this case, by: URL, Adapted type, Number of active tabs, Number of large resources, Experience level.


Copy
web.page_views
| summarize web.page_views = number_of_page_view.sum() 
Single-metric gauge chart
Ratio of devices or users when there is a bad event
Create a single-metric gauge chart displaying the ratio of devices or users when there is a bad event, for example, a crash. It allows to see how devices or users are affected by the issue.


Copy
<devices|users>
| include <event table>
| compute temp_bad_number = <device|user>.count()
| summarize 
   <metric> = temp_bad_number.sum(), 
   <total> = count()
Example

Copy
devices
| include execution.crashes
| compute crash_cnt = device.count()
| summarize 
   devices_with_crashes = crash_cnt.sum(), 
   total_devices = count()

Copy
devices
| include execution.crashes
| compute crash_cnt = device.countif(process_visibility == foreground)
| summarize
   devices_with_crashes = crash_cnt.sum(),
   total_devices = count()
Ratio of events
Display the ratio of events when there is an event such as a crash, freeze, hard reset, system reset.


Copy
<devices|users>
| include <bad event table>
| compute temp_metric_number = count()
| include <total event table>
| compute temp_total_number = count()
| summarize 
    <metric> = temp_metric_number.sum(), 
    <total> = temp_total_number.sum()
Example

Copy
devices 
| include collaboration.sessions 
| where video.quality == poor or audio.quality == poor 
| compute num_poor_quality_sessions = id.count() 
| include collaboration.sessions 
| compute num_total_sessions = id.count() 
| summarize
    poor_quality = num_poor_quality_sessions.sum(), 
    acceptable_quality = num_total_sessions.sum()
Score metric
Display the DEX score metric.


Copy
<score table>
| summarize <metric> = <score_field>.avg(), <total> = <total>
Example


Copy
dex.scores
| summarize score = value.avg() , total = 100
Multi-metric gauge chart
Ratio of devices or users with bad events against objects without them

Copy
<devices|users>
| include <event table>
| compute temp_bad_number = <device|user>.count()
| summarize 
   <good_label> = count() - temp_bad_number.sum(), 
   <bad_label> = temp_bad_number.sum()
Example


Copy
devices
| include execution.crashes
| compute crash_cnt = device.count()
| summarize 
    without_crashes = count() - crash_cnt.sum(), 
    with_crashes = crash_cnt.sum()
Ratio of devices with bad events against devices without them

Copy
devices
| include <bad event table>
| compute temp_bad_number = count()
| include <total event table>
| compute temp_total_number = count()
| summarize 
   <good_label> = temp_total_number.sum() - temp_bad_number.sum(), 
   <bad_label> = temp_bad_number.sum()
Example


Copy
devices
| include device_performance.hard_resets
| compute hard_reset_cnt = number_of_hard_resets.sum()
| include device_performance.events
| compute total_cnt = count()
| summarize 
   no_hard_resets = total_cnt.sum() - hard_reset_cnt.sum(), 
   hard_resets = hard_reset_cnt.sum()
Ratio of users or devices with a good state against the ones with a bad state

Copy
<devices|users>
| include <event table>
| where <condition is bad>
| compute temp_bad_number = <device|user>.count()
| include <event table>
| where <condition is good>
| compute temp_good_number = <device|user>.count()
| summarize 
   <good_label> = temp_good_number.sum(), 
   <bad_label> = temp_bad_number.sum()
Example


Copy
users
| include web.page_views
| where experience_level == frustrating
| compute frustrating_cnt = user.count()
| include web.page_views
| where experience_level == good 
| compute good_cnt = user.count()
| summarize 
   good = good_cnt.sum(), 
   frustrating = frustrating_cnt.sum()
Ratio of events with a good state against events with a bad state

Copy
<devices|users>
| include <event table>
| where <condition is bad>
| compute temp_bad_number = <sum|count>
| include <event table>
| where <condition is good>
| compute temp_good_number = <sum|count>
| summarize 
   <good_label> = temp_good_number.sum(), 
  <bad_label> = temp_bad_number.sum()
Example


Copy
users
| include web.page_views
| where experience_level == frustrating
| compute frustrating_cnt = number_of_page_views.sum() 
| include web.page_views
| where experience_level == good 
| compute good_cnt = number_of_page_views.sum() 
| summarize 
   good = good_cnt.sum(), 
   frustrating = frustrating_cnt.sum()


Detecting issues impacting multiple devices

Detecting issues impacting a single device or user

Live Dashboards NQL examples

Investigations NQL examples

Workflows NQL examples

Campaigns NQL examples

Remote Actions NQL examples

Software Metering NQL examples

Custom trends NQL examples

Webhooks NQL examples