# Mattermost User ID Converter

## Overview

This utility parses a supplied CSV file for a specified User ID field and looks up the matching username or full name in Mattermost. It's particularly useful for CSV exports from Boards or any situation where you have a CSV file containing Mattermost User IDs. Developed in Go, this tool allows for easy distribution and extension without the need for a Go development environment.

## Features

- **User Lookup**: Converts Mattermost User IDs to either the username or the full name (if available).
- **Flexible Output**: Users can opt to output either the username or the full name.
- **Easy to Use**: Simple command line interface for easy operation.
- **Customisable**: Supports various configurations through command line arguments or environment variables.

## Prerequisites

To use this utility, you'll need:

- A personal access token from your Mattermost instance, with the requisite permissions to read user details. More information on obtaining a token can be found [here](https://developers.mattermost.com/integrate/reference/personal-access-token/).
- Your Mattermost URL, the HTTP scheme (http/https), and the port in use. This information is essential for making API calls to Mattermost.

## Installation

As the utility is written in Go and compiled into an executable, installation is as simple as downloading the relevant executable for your platform from the GitHub releases page and placing it in a suitable location on your system's path.

## Usage

Here's how you can use the utility:

```bash
./mattermost_user_id_converter -url <your_mattermost_url> -port <mattermost_port> -scheme <http_or_https> -token <your_personal_access_token> -infile <path_to_input_csv> -outfile <path_to_output_csv> -column <user_id_column_name> [-fullname] [-debug]
```

Some of the parameters can also be supplied as environment variables, as seen in the table below:

| Command Line Arg | Environment Variable | Description |
| --- | --- | --- |
| `-url` | `MM_URL` | The URL of the Mattermost instance (without a schema) |
| `-port` | `MM_PORT` | The Mattermost port to be used [default: 8065] |
| `-scheme` | `MM_SCHEME` | The HHTP scheme to be used (http/https) [default: http] |
| `-token` | `MM_TOKEN` | The user token for Mattermost.  Note that this user must have the appropriate rights to read users. |
| `-infile` |  | **Required**. CSV file for parsing. |
| `-outfile` |  | **Required**.  Filename of output file.   |
| `-column` |  | **Required**. The name of the column containing the User ID. |
| `-fullname` |  | **Optional**. If present, the output will contain the full name of the user, rather than their Mattermost username. |
| `-debug` | MM_DEBUG | If present, will run in debug mode, delivering additional output to stdout |

## Contributing

Your contributions are welcome! If you have suggestions for improving this utility, please open an issue or pull request in this repository. Whether it's bug reports, feature requests, or contributions to code, all are greatly appreciated.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

## Contact

For questions, feedback, or contributions regarding this project, please use the following methods:

- **Issues and Pull Requests**: For specific questions, issues, or suggestions for improvements, feel free to open an issue or a pull request in this repository.
- **Mattermost Community**: Join us in the Mattermost Community server, where we discuss all things related to extending Mattermost. You can find us in the channel [Integrations and Apps](https://community.mattermost.com/core/channels/integrations).
- **Social Media**: Follow and message me on Twitter, where I'm [@jlandells](https://twitter.com/jlandells).